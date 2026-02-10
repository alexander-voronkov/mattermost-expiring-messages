package main

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/mattermost/mattermost/server/public/model"
)

var durationRegex = regexp.MustCompile(`^(\d+)([mhd])$`)

func calculateExpiresAt(duration string) int64 {
	milliseconds, err := parseDuration(duration)
	if err != nil {
		return model.GetMillis() + (5 * 60 * 1000)
	}
	return model.GetMillis() + milliseconds
}

func parseDuration(duration string) (int64, error) {
	matches := durationRegex.FindStringSubmatch(duration)
	if matches == nil {
		return 0, fmt.Errorf("invalid duration format: %s", duration)
	}

	value, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, fmt.Errorf("invalid duration value: %s", matches[1])
	}

	unit := matches[2]

	switch unit {
	case "m":
		return int64(value) * int64(time.Minute/time.Millisecond), nil
	case "h":
		return int64(value) * int64(time.Hour/time.Millisecond), nil
	case "d":
		return int64(value) * int64(24*time.Hour/time.Millisecond), nil
	default:
		return 0, fmt.Errorf("invalid duration unit: %s", unit)
	}
}
