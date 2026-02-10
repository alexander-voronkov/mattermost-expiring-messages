package server

import (
	"strconv"
	"time"

	"github.com/mattermost/mattermost/server/v8/model"
)

const (
	queueBucketSizeMinutes = 1
	maxPostsPerDeletion    = 100
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
	return "expiration_bucket_" + strconv.FormatInt(bucketMinute, 10) + "_"
}

func (p *Plugin) cleanupOldBuckets(now time.Time) {
	cutoff := now.Add(-24 * time.Hour)
	cutoffBucket := strconv.FormatInt(cutoff.Unix()/60, 10)

	prefix := "expiration_bucket_"

	for i := 0; i < 100; i++ {
		keys, appErr := p.API.KVList(prefix, i*100, 100)
		if appErr != nil {
			break
		}

		if len(keys) == 0 {
			break
		}

		for _, key := range keys {
			bucketNum := key[len("expiration_bucket_"):]
			bucketNum = bucketNum[:len(bucketNum)-1]

			if bucketNum < cutoffBucket {
				if appErr := p.API.KVDelete(key); appErr != nil {
					p.API.LogError("Failed to delete old bucket", "key", key, "error", appErr.Error())
				}
			}
		}
	}
}
