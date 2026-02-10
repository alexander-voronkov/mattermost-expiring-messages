package server

import (
	"sync"
	"time"

	"github.com/mattermost/mattermost/server/v8/model"
	"github.com/mattermost/mattermost/server/v8/plugin"
)

type Plugin struct {
	plugin.MattermostPlugin
	configuration     configuration
	configurationLock sync.RWMutex

	stopExpirationJob chan struct{}
}

func (p *Plugin) OnActivate() error {
	if err := p.ensureIndex(); err != nil {
		p.API.LogError("Failed to create database index", "error", err.Error())
		return err
	}

	p.stopExpirationJob = make(chan struct{})
	go p.runExpirationJob()

	return nil
}

func (p *Plugin) OnDeactivate() error {
	if p.stopExpirationJob != nil {
		close(p.stopExpirationJob)
		p.stopExpirationJob = nil
	}
	return nil
}

func (p *Plugin) MessageWillBePosted(_ *plugin.Context, post *model.Post) (*model.Post, string) {
	p.configurationLock.RLock()
	defer p.configurationLock.RUnlock()

	if !p.configuration.Enabled {
		return post, ""
	}

	if post == nil || post.Props == nil {
		return post, ""
	}

	ttlInterface, ok := post.Props["ttl"]
	if !ok {
		return post, ""
	}

	ttl, ok := ttlInterface.(map[string]interface{})
	if !ok {
		return post, ""
	}

	enabled, _ := ttl["enabled"].(bool)
	if !enabled {
		delete(post.Props, "ttl")
		return post, ""
	}

	duration, ok := ttl["duration"].(string)
	if !ok || duration == "" {
		return post, "TTL duration is required when TTL is enabled"
	}

	if !p.isDurationAllowed(duration) {
		return post, "TTL duration '" + duration + "' is not allowed"
	}

	expiresAt := calculateExpiresAt(duration)
	ttl["expires_at"] = float64(expiresAt)
	post.Props["ttl"] = ttl

	return post, ""
}

func (p *Plugin) PostHasBeenPosted(_ *plugin.Context, post *model.Post) {
	if post == nil || post.Props == nil {
		return
	}

	ttlInterface, ok := post.Props["ttl"]
	if !ok {
		return
	}

	ttl, ok := ttlInterface.(map[string]interface{})
	if !ok {
		return
	}

	enabled, _ := ttl["enabled"].(bool)
	if !enabled {
		return
	}

	expiresAtFloat, ok := ttl["expires_at"].(float64)
	if !ok {
		return
	}

	p.queuePostForDeletion(post.Id, int64(expiresAtFloat))
}

func (p *Plugin) MessageWillBeUpdated(_ *plugin.Context, newPost *model.Post, oldPost *model.Post) (*model.Post, string) {
	if newPost == nil || newPost.Props == nil {
		return newPost, ""
	}

	p.configurationLock.RLock()
	defer p.configurationLock.RUnlock()

	if !p.configuration.Enabled {
		return newPost, ""
	}

	ttlInterface, ok := newPost.Props["ttl"]
	if !ok {
		return newPost, ""
	}

	ttl, ok := ttlInterface.(map[string]interface{})
	if !ok {
		return newPost, ""
	}

	enabled, _ := ttl["enabled"].(bool)
	if !enabled {
		delete(newPost.Props, "ttl")
		return newPost, ""
	}

	duration, ok := ttl["duration"].(string)
	if !ok || duration == "" {
		return newPost, "TTL duration is required when TTL is enabled"
	}

	if !p.isDurationAllowed(duration) {
		return newPost, "TTL duration '" + duration + "' is not allowed"
	}

	expiresAt := calculateExpiresAt(duration)
	ttl["expires_at"] = float64(expiresAt)
	newPost.Props["ttl"] = ttl

	oldExpiresAt := float64(0)
	if oldPost != nil && oldPost.Props != nil {
		if oldTTL, ok := oldPost.Props["ttl"].(map[string]interface{}); ok {
			oldExpiresAt, _ = oldTTL["expires_at"].(float64)
		}
	}

	if int64(expiresAt) != int64(oldExpiresAt) {
		p.queuePostForDeletion(newPost.Id, int64(expiresAt))
	}

	return newPost, ""
}

func (p *Plugin) runExpirationJob() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			p.deleteExpiredPosts()
		case <-p.stopExpirationJob:
			return
		}
	}
}

func (p *Plugin) deleteExpiredPosts() {
	now := model.GetMillis()
	nowTime := time.Unix(now/1000, 0)

	currentBucket := getExpirationBucketKey(nowTime)

	keys, appErr := p.API.KVList(currentBucket, 0, maxPostsPerDeletion)
	if appErr != nil {
		return
	}

	for _, key := range keys {
		postIDBytes, appErr := p.API.KVGet(currentBucket + key)
		if appErr != nil {
			continue
		}

		postID := string(postIDBytes)
		if postID == "" {
			continue
		}

		if err := p.API.DeletePost(postID); err != nil {
			p.API.LogError("Failed to delete expired post", "post_id", postID, "error", err.Error())
		} else {
			p.API.LogInfo("Deleted expired post", "post_id", postID)
		}

		if appErr := p.API.KVDelete(currentBucket + key); appErr != nil {
			p.API.LogError("Failed to delete expiration key", "key", currentBucket+key, "error", appErr.Error())
		}
	}

	p.cleanupOldBuckets(nowTime)
}
