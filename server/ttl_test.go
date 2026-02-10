package server

import (
	"testing"
	"time"

	"github.com/mattermost/mattermost/server/v8/model"
)

func TestParseDuration(t *testing.T) {
	tests := []struct {
		name      string
		duration  string
		want      int64
		wantError bool
	}{
		{
			name:      "5 minutes",
			duration:  "5m",
			want:      5 * 60 * 1000,
			wantError: false,
		},
		{
			name:      "15 minutes",
			duration:  "15m",
			want:      15 * 60 * 1000,
			wantError: false,
		},
		{
			name:      "1 hour",
			duration:  "1h",
			want:      60 * 60 * 1000,
			wantError: false,
		},
		{
			name:      "2 hours",
			duration:  "2h",
			want:      2 * 60 * 60 * 1000,
			wantError: false,
		},
		{
			name:      "1 day",
			duration:  "1d",
			want:      24 * 60 * 60 * 1000,
			wantError: false,
		},
		{
			name:      "3 days",
			duration:  "3d",
			want:      3 * 24 * 60 * 60 * 1000,
			wantError: false,
		},
		{
			name:      "invalid format - no unit",
			duration:  "5",
			want:      0,
			wantError: true,
		},
		{
			name:      "invalid format - invalid unit",
			duration:  "5x",
			want:      0,
			wantError: true,
		},
		{
			name:      "invalid format - negative number",
			duration:  "-5m",
			want:      0,
			wantError: true,
		},
		{
			name:      "invalid format - empty string",
			duration:  "",
			want:      0,
			wantError: true,
		},
		{
			name:      "invalid format - non-numeric",
			duration:  "am",
			want:      0,
			wantError: true,
		},
		{
			name:      "large value minutes",
			duration:  "999m",
			want:      999 * 60 * 1000,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseDuration(tt.duration)
			if (err != nil) != tt.wantError {
				t.Errorf("parseDuration() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if got != tt.want {
				t.Errorf("parseDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculateExpiresAt(t *testing.T) {
	// Set a fixed mock time
	mockTime := int64(1700000000000) // 2023-11-15 04:26:40 UTC

	tests := []struct {
		name     string
		duration string
	}{
		{"5 minutes", "5m"},
		{"15 minutes", "15m"},
		{"1 hour", "1h"},
		{"1 day", "1d"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We can't easily mock model.GetMillis(), so we just verify the result is reasonable
			result := calculateExpiresAt(tt.duration)
			now := model.GetMillis()

			// Result should be in the future
			if result <= now {
				t.Errorf("calculateExpiresAt() result %v is not in the future (now: %v)", result, now)
			}

			// Result should not be more than the duration plus some tolerance
			maxExpected := now + (30 * 24 * time.Hour / time.Millisecond) // 30 days max
			if result > maxExpected {
				t.Errorf("calculateExpiresAt() result %v is too far in the future (max: %v)", result, maxExpected)
			}
		})
	}
}

func TestGetExpirationBucketKey(t *testing.T) {
	tests := []struct {
		name string
		time time.Time
		want string
	}{
		{
			name: "epoch time",
			time: time.Unix(0, 0),
			want: "expiration_bucket_0_",
		},
		{
			name: "1 minute after epoch",
			time: time.Unix(60, 0),
			want: "expiration_bucket_1_",
		},
		{
			name: "100 minutes after epoch",
			time: time.Unix(6000, 0),
			want: "expiration_bucket_100_",
		},
		{
			name: "arbitrary time",
			time: time.Unix(1700000000, 0),
			want: "expiration_bucket_28333333_",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getExpirationBucketKey(tt.time); got != tt.want {
				t.Errorf("getExpirationBucketKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsDurationAllowed(t *testing.T) {
	tests := []struct {
		name             string
		allowedDurations []string
		duration         string
		want             bool
	}{
		{
			name:             "empty allowed list allows all",
			allowedDurations: []string{},
			duration:         "5m",
			want:             true,
		},
		{
			name:             "duration in allowed list",
			allowedDurations: []string{"5m", "15m", "1h"},
			duration:         "5m",
			want:             true,
		},
		{
			name:             "duration not in allowed list",
			allowedDurations: []string{"5m", "15m", "1h"},
			duration:         "1d",
			want:             false,
		},
		{
			name:             "all standard durations allowed",
			allowedDurations: []string{"5m", "15m", "1h", "1d"},
			duration:         "1h",
			want:             true,
		},
		{
			name:             "custom duration not in list",
			allowedDurations: []string{"5m", "1h"},
			duration:         "30m",
			want:             false,
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
