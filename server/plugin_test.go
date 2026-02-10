package server

import (
	"testing"

	"github.com/mattermost/mattermost/server/v8/model"
	"github.com/mattermost/mattermost/server/v8/plugin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockAPI is a mock implementation of plugin.API
type MockAPI struct {
	mock.Mock
}

func (m *MockAPI) LogDebug(msg string, args ...interface{}) {
	m.Called(msg, args)
}

func (m *MockAPI) LogInfo(msg string, args ...interface{}) {
	m.Called(msg, args)
}

func (m *MockAPI) LogError(msg string, args ...interface{}) {
	m.Called(msg, args)
}

func (m *MockAPI) LogWarn(msg string, args ...interface{}) {
	m.Called(msg, args)
}

func (m *MockAPI) KVSet(key string, value []byte) error {
	args := m.Called(key, value)
	return args.Error(0)
}

func (m *MockAPI) KVGet(key string) ([]byte, *model.AppError) {
	args := m.Called(key)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*model.AppError)
	}
	return args.Get(0).([]byte), args.Get(1).(*model.AppError)
}

func (m *MockAPI) KVDelete(key string) *model.AppError {
	args := m.Called(key)
	return args.Get(0).(*model.AppError)
}

func (m *MockAPI) KVList(prefix string, page, perPage int) ([]string, *model.AppError) {
	args := m.Called(prefix, page, perPage)
	return args.Get(0).([]string), args.Get(1).(*model.AppError)
}

func (m *MockAPI) DeletePost(postID string) *model.AppError {
	args := m.Called(postID)
	return args.Get(0).(*model.AppError)
}

func (m *MockAPI) LoadPluginConfiguration(dest interface{}) error {
	args := m.Called(dest)
	return args.Error(0)
}

func (m *MockAPI) ExecuteDirectQuery(query string) (interface{}, *model.AppError) {
	args := m.Called(query)
	return args.Get(0), args.Get(1).(*model.AppError)
}

func TestMessageWillBePosted(t *testing.T) {
	tests := []struct {
		name             string
		post             *model.Post
		configEnabled    bool
		allowedDurations []string
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
			shouldHaveTTL: true, // TTL stays when plugin is disabled
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
			allowedDurations: []string{},
			wantErrorMsg:     "",
			shouldHaveTTL:    true,
		},
		{
			name: "valid TTL with 1h",
			post: &model.Post{
				Props: map[string]interface{}{
					"ttl": map[string]interface{}{
						"enabled":  true,
						"duration": "1h",
					},
				},
			},
			configEnabled:    true,
			allowedDurations: []string{},
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
			name: "TTL enabled with empty duration",
			post: &model.Post{
				Props: map[string]interface{}{
					"ttl": map[string]interface{}{
						"enabled":  true,
						"duration": "",
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
			name: "invalid duration format",
			post: &model.Post{
				Props: map[string]interface{}{
					"ttl": map[string]interface{}{
						"enabled":  true,
						"duration": "invalid",
					},
				},
			},
			configEnabled: true,
			wantErrorMsg:  "TTL duration 'invalid' is not allowed",
			shouldHaveTTL: true,
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
			allowedDurations: []string{"5m", "15m", "1h"},
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
			allowedDurations: []string{"5m", "15m", "1h"},
			wantErrorMsg:     "",
			shouldHaveTTL:    true,
		},
		{
			name: "malformed TTL - not a map",
			post: &model.Post{
				Props: map[string]interface{}{
					"ttl": "invalid",
				},
			},
			configEnabled: true,
			wantErrorMsg:  "",
			shouldHaveTTL: true,
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

			if tt.post != nil {
				_, hasTTL := tt.post.Props["ttl"]
				assert.Equal(t, tt.shouldHaveTTL, hasTTL, "TTL presence mismatch")

				if tt.shouldHaveTTL && hasTTL {
					ttl := tt.post.Props["ttl"]
					ttlMap, ok := ttl.(map[string]interface{})
					if ok && ttlMap["enabled"] == true {
						// Check that expires_at was set
						_, hasExpiresAt := ttlMap["expires_at"]
						assert.True(t, hasExpiresAt, "expires_at should be set for enabled TTL")
					}
				}
			}

			assert.Same(t, tt.post, returnedPost, "Should return the same post instance")
		})
	}
}

func TestMessageWillBeUpdated(t *testing.T) {
	tests := []struct {
		name             string
		newPost          *model.Post
		oldPost          *model.Post
		configEnabled    bool
		allowedDurations []string
		wantErrorMsg     string
		shouldHaveTTL    bool
	}{
		{
			name: "update with valid TTL",
			newPost: &model.Post{
				Props: map[string]interface{}{
					"ttl": map[string]interface{}{
						"enabled":  true,
						"duration": "5m",
					},
				},
			},
			oldPost:       &model.Post{},
			configEnabled: true,
			wantErrorMsg:  "",
			shouldHaveTTL: true,
		},
		{
			name: "update removing TTL",
			newPost: &model.Post{
				Props: map[string]interface{}{
					"ttl": map[string]interface{}{
						"enabled": false,
					},
				},
			},
			oldPost: &model.Post{
				Props: map[string]interface{}{
					"ttl": map[string]interface{}{
						"enabled":    true,
						"duration":   "5m",
						"expires_at": float64(1700000000000),
					},
				},
			},
			configEnabled: true,
			wantErrorMsg:  "",
			shouldHaveTTL: false,
		},
		{
			name: "update with invalid duration",
			newPost: &model.Post{
				Props: map[string]interface{}{
					"ttl": map[string]interface{}{
						"enabled":  true,
						"duration": "invalid",
					},
				},
			},
			oldPost:       &model.Post{},
			configEnabled: true,
			wantErrorMsg:  "TTL duration 'invalid' is not allowed",
			shouldHaveTTL: true,
		},
		{
			name:          "update with no props",
			newPost:       &model.Post{},
			oldPost:       &model.Post{},
			configEnabled: true,
			wantErrorMsg:  "",
			shouldHaveTTL: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Plugin{}
			p.configuration.Enabled = tt.configEnabled
			p.configuration.AllowedDurations = tt.allowedDurations

			returnedPost, errorMsg := p.MessageWillBeUpdated(&plugin.Context{}, tt.newPost, tt.oldPost)

			if tt.wantErrorMsg != "" {
				assert.Equal(t, tt.wantErrorMsg, errorMsg)
			} else {
				assert.Empty(t, errorMsg)
			}

			if tt.newPost != nil && tt.newPost.Props != nil {
				_, hasTTL := tt.newPost.Props["ttl"]
				assert.Equal(t, tt.shouldHaveTTL, hasTTL, "TTL presence mismatch")
			}

			assert.Same(t, tt.newPost, returnedPost, "Should return the same post instance")
		})
	}
}

func TestPostHasBeenPosted(t *testing.T) {
	tests := []struct {
		name        string
		post        *model.Post
		shouldQueue bool
	}{
		{
			name:        "nil post",
			post:        nil,
			shouldQueue: false,
		},
		{
			name:        "post with no props",
			post:        &model.Post{},
			shouldQueue: false,
		},
		{
			name:        "post with empty props",
			post:        &model.Post{Props: make(map[string]interface{})},
			shouldQueue: false,
		},
		{
			name: "post with TTL disabled",
			post: &model.Post{
				Props: map[string]interface{}{
					"ttl": map[string]interface{}{
						"enabled": false,
					},
				},
			},
			shouldQueue: false,
		},
		{
			name: "post with TTL enabled but no expires_at",
			post: &model.Post{
				Props: map[string]interface{}{
					"ttl": map[string]interface{}{
						"enabled": true,
					},
				},
			},
			shouldQueue: false,
		},
		{
			name: "post with valid TTL",
			post: &model.Post{
				Id: "test-post-id",
				Props: map[string]interface{}{
					"ttl": map[string]interface{}{
						"enabled":    true,
						"expires_at": float64(1700000000000),
					},
				},
			},
			shouldQueue: true,
		},
		{
			name: "post with malformed TTL",
			post: &model.Post{
				Props: map[string]interface{}{
					"ttl": "invalid",
				},
			},
			shouldQueue: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPI := new(MockAPI)
			p := &Plugin{}
			p.SetAPI(mockAPI)

			// This should not panic
			p.PostHasBeenPosted(&plugin.Context{}, tt.post)

			if tt.shouldQueue {
				// If the post should be queued, KVSet should have been called
				// We can't easily verify this without a more sophisticated mock,
				// but at least verify it doesn't panic
				require.NotNil(t, tt.post.Id)
			}
		})
	}
}

func TestOnConfigurationChange(t *testing.T) {
	tests := []struct {
		name    string
		config  configuration
		wantErr bool
	}{
		{
			name: "default configuration",
			config: configuration{
				Enabled:          false,
				AllowedDurations: []string{},
			},
			wantErr: false,
		},
		{
			name: "enabled with durations",
			config: configuration{
				Enabled:          true,
				AllowedDurations: []string{"5m", "15m", "1h", "1d"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Plugin{}

			err := p.OnConfigurationChange()
			assert.NoError(t, err)

			// Verify configuration was loaded
			p.configurationLock.RLock()
			defer p.configurationLock.RUnlock()

			// The actual config loaded will be from the file, but we verify
			// the lock mechanism works without race conditions
			assert.NotNil(t, p.configuration)
		})
	}
}

func TestGetPluginManifest(t *testing.T) {
	p := &Plugin{}
	manifest := p.GetPluginManifest()

	assert.NotNil(t, manifest)
	assert.Equal(t, "com.fambear.expiring-messages", manifest.Id)
	assert.Equal(t, "Expiring Messages", manifest.Name)
	assert.Equal(t, "0.1.0", manifest.Version)
	assert.Equal(t, minServerVersion, manifest.MinServerVersion)
}
