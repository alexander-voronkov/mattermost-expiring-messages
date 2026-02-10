package main

import (
	"strconv"
	"strings"
	"time"
)

const (
	maxPostsPerDeletion = 100
	expirationPrefix    = "expiration_bucket_"
)

func (p *Plugin) queuePostForDeletion(postID string, expiresAt int64) {
	expiresTime := time.Unix(expiresAt/1000, 0)
	bucketKey := getExpirationBucketKey(expiresTime)

	if err := p.API.KVSet(bucketKey+postID, []byte(postID)); err != nil {
		p.API.LogError("Failed to queue post for deletion", "post_id", postID, "error", err.Error())
	}
}

func getExpirationBucketKey(t time.Time) string {
	bucketMinute := t.Unix() / 60
	return expirationPrefix + strconv.FormatInt(bucketMinute, 10) + "_"
}

func extractBucketNumber(key string) (int64, bool) {
	if !strings.HasPrefix(key, expirationPrefix) {
		return 0, false
	}
	rest := key[len(expirationPrefix):]
	// Find the trailing underscore
	idx := strings.Index(rest, "_")
	if idx == -1 {
		return 0, false
	}
	numStr := rest[:idx]
	num, err := strconv.ParseInt(numStr, 10, 64)
	if err != nil {
		return 0, false
	}
	return num, true
}

func (p *Plugin) cleanupOldBuckets(now time.Time) {
	cutoffBucket := now.Add(-24*time.Hour).Unix() / 60

	for page := 0; page < 100; page++ {
		keys, appErr := p.API.KVList(page, 100)
		if appErr != nil {
			p.API.LogError("Failed to list KV keys", "error", appErr.Error())
			break
		}

		if len(keys) == 0 {
			break
		}

		for _, key := range keys {
			if !strings.HasPrefix(key, expirationPrefix) {
				continue
			}

			bucketNum, ok := extractBucketNumber(key)
			if !ok {
				continue
			}

			if bucketNum < cutoffBucket {
				if appErr := p.API.KVDelete(key); appErr != nil {
					p.API.LogError("Failed to delete old bucket", "key", key, "error", appErr.Error())
				}
			}
		}
	}
}
