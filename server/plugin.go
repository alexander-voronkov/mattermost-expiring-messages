package main

import (
	"strings"
	"sync"
	"time"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
)

type Plugin struct {
	plugin.MattermostPlugin
	configuration     configuration
	configurationLock sync.RWMutex

	stopExpirationJob chan struct{}
}

func (p *Plugin) OnActivate() error {
	if err := p.ensureIndex(); err != nil {
		p.API.LogError("Failed to initialize plugin", "error", err.Error())
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

	// Set custom post type so webapp can render with countdown
	post.Type = "custom_expiring"

	return post, ""
}

func (p *Plugin) MessageHasBeenPosted(_ *plugin.Context, post *model.Post) {
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
	// Also check the previous bucket in case we missed it
	prevBucket := getExpirationBucketKey(nowTime.Add(-1 * time.Minute))

	bucketsToCheck := []string{currentBucket, prevBucket}

	// KVList returns all keys, we filter by prefix
	for page := 0; page < 10; page++ {
		keys, appErr := p.API.KVList(page, maxPostsPerDeletion)
		if appErr != nil {
			p.API.LogError("Failed to list KV keys", "error", appErr.Error())
			break
		}

		if len(keys) == 0 {
			break
		}

		for _, key := range keys {
			matchesBucket := false
			for _, bucket := range bucketsToCheck {
				if strings.HasPrefix(key, bucket) {
					matchesBucket = true
					break
				}
			}
			if !matchesBucket {
				continue
			}

			postIDBytes, appErr := p.API.KVGet(key)
			if appErr != nil {
				continue
			}

			postID := string(postIDBytes)
			if postID == "" {
				continue
			}

			// Use permanent delete via HTTP API to avoid "(message deleted)" placeholder
			if err := p.permanentDeletePost(postID); err != nil {
				p.API.LogError("Failed to permanently delete expired post", "post_id", postID, "error", err.Error())
			} else {
				p.API.LogInfo("Permanently deleted expired post", "post_id", postID)
			}

			if appErr := p.API.KVDelete(key); appErr != nil {
				p.API.LogError("Failed to delete expiration key", "key", key, "error", appErr.Error())
			}
		}
	}

	p.cleanupOldBuckets(nowTime)
}

// permanentDeletePost attempts to permanently delete a post
// For now, we use the standard DeletePost which may show "(message deleted)"
// TODO: Implement permanent delete via REST API with proper auth when bot user is configured
func (p *Plugin) permanentDeletePost(postID string) error {
	// Standard delete - may show "(message deleted)" placeholder
	// To enable true permanent delete:
	// 1. Configure a bot user with system admin permissions
	// 2. Use REST API: DELETE /api/v4/posts/{post_id}?permanent=true
	appErr := p.API.DeletePost(postID)
	if appErr != nil {
		return appErr
	}
	return nil
}
