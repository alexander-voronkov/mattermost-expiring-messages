package main

import (
	"strings"
)

type configuration struct {
	Enabled          bool   `json:"enabled"`
	AllowedDurations string `json:"allowed_durations"`
}

const (
	minServerVersion = "9.0.0"
	defaultDurations = "5m,15m,1h,1d"
)

func (p *Plugin) OnConfigurationChange() error {
	var cfg configuration
	if err := p.API.LoadPluginConfiguration(&cfg); err != nil {
		return err
	}

	p.configurationLock.Lock()
	defer p.configurationLock.Unlock()

	p.configuration = cfg
	return nil
}

func (p *Plugin) getAllowedDurations() []string {
	durations := p.configuration.AllowedDurations
	if durations == "" {
		durations = defaultDurations
	}

	parts := strings.Split(durations, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func (p *Plugin) isDurationAllowed(duration string) bool {
	allowed := p.getAllowedDurations()
	if len(allowed) == 0 {
		return true
	}

	for _, a := range allowed {
		if a == duration {
			return true
		}
	}
	return false
}
