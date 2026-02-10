package main

import (
	"testing"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
	"github.com/stretchr/testify/assert"
)

func TestMessageWillBePosted(t *testing.T) {
	tests := []struct {
		name             string
		post             *model.Post
		configEnabled    bool
		allowedDurations string
		wantErrorMsg     string
		shouldHaveTTL    bool
	}{
		{
			name:          "nil post",
			post:          nil,
			configEnabled: true,
			wantErrorMsg:  "",
			shouldHaveTTL: false,
		},
		{
			name:          "post with no props",
			post:          &model.Post{},
			configEnabled: true,
			wantErrorMsg:  "",
			shouldHaveTTL: false,
		},
		{
			name:          "post with empty props",
			post:          &model.Post{Props: make(map[string]interface{})},
			configEnabled: true,
			wantErrorMsg:  "",
			shouldHaveTTL: false,
		},
		{
			name: "plugin disabled",
			post: &model.Post{
				Props: map[string]interface{}{
					"ttl": map[string]interface{}{
						"enabled":  true,
						"duration": "5m",
					},
				},
			},
			configEnabled: false,
			wantErrorMsg:  "",
			shouldHaveTTL: true,
		},
		{
			name: "valid TTL with 5m",
			post: &model.Post{
				Props: map[string]interface{}{
					"ttl": map[string]interface{}{
						"enabled":  true,
						"duration": "5m",
					},
				},
			},
			configEnabled:    true,
			allowedDurations: "",
			wantErrorMsg:     "",
			shouldHaveTTL:    true,
		},
		{
			name: "TTL enabled but no duration",
			post: &model.Post{
				Props: map[string]interface{}{
					"ttl": map[string]interface{}{
						"enabled": true,
					},
				},
			},
			configEnabled: true,
			wantErrorMsg:  "TTL duration is required when TTL is enabled",
			shouldHaveTTL: true,
		},
		{
			name: "TTL disabled",
			post: &model.Post{
				Props: map[string]interface{}{
					"ttl": map[string]interface{}{
						"enabled": false,
					},
				},
			},
			configEnabled: true,
			wantErrorMsg:  "",
			shouldHaveTTL: false,
		},
		{
			name: "duration not in allowed list",
			post: &model.Post{
				Props: map[string]interface{}{
					"ttl": map[string]interface{}{
						"enabled":  true,
						"duration": "1d",
					},
				},
			},
			configEnabled:    true,
			allowedDurations: "5m,15m,1h",
			wantErrorMsg:     "TTL duration '1d' is not allowed",
			shouldHaveTTL:    true,
		},
		{
			name: "duration in allowed list",
			post: &model.Post{
				Props: map[string]interface{}{
					"ttl": map[string]interface{}{
						"enabled":  true,
						"duration": "15m",
					},
				},
			},
			configEnabled:    true,
			allowedDurations: "5m,15m,1h",
			wantErrorMsg:     "",
			shouldHaveTTL:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Plugin{}
			p.configuration.Enabled = tt.configEnabled
			p.configuration.AllowedDurations = tt.allowedDurations

			returnedPost, errorMsg := p.MessageWillBePosted(&plugin.Context{}, tt.post)

			if tt.wantErrorMsg != "" {
				assert.Equal(t, tt.wantErrorMsg, errorMsg)
			} else {
				assert.Empty(t, errorMsg)
			}

			if tt.post != nil && tt.post.Props != nil {
				_, hasTTL := tt.post.Props["ttl"]
				assert.Equal(t, tt.shouldHaveTTL, hasTTL, "TTL presence mismatch")

				if tt.shouldHaveTTL && hasTTL && tt.wantErrorMsg == "" {
					ttl := tt.post.Props["ttl"]
					ttlMap, ok := ttl.(map[string]interface{})
					if ok && ttlMap["enabled"] == true {
						_, hasExpiresAt := ttlMap["expires_at"]
						if tt.configEnabled {
							assert.True(t, hasExpiresAt, "expires_at should be set for enabled TTL")
						}
					}
				}
			}

			assert.Same(t, tt.post, returnedPost, "Should return the same post instance")
		})
	}
}

func TestIsDurationAllowed(t *testing.T) {
	tests := []struct {
		name             string
		allowedDurations string
		duration         string
		want             bool
	}{
		{
			name:             "empty allowed list allows all",
			allowedDurations: "",
			duration:         "5m",
			want:             true,
		},
		{
			name:             "duration in allowed list",
			allowedDurations: "5m,15m,1h",
			duration:         "5m",
			want:             true,
		},
		{
			name:             "duration not in allowed list",
			allowedDurations: "5m,15m,1h",
			duration:         "1d",
			want:             false,
		},
		{
			name:             "all standard durations allowed",
			allowedDurations: "5m,15m,1h,1d",
			duration:         "1h",
			want:             true,
		},
		{
			name:             "handles whitespace",
			allowedDurations: "5m, 15m, 1h",
			duration:         "15m",
			want:             true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Plugin{}
			p.configuration.AllowedDurations = tt.allowedDurations
			if got := p.isDurationAllowed(tt.duration); got != tt.want {
				t.Errorf("Plugin.isDurationAllowed() = %v, want %v", got, tt.want)
			}
		})
	}
}
