package server

import (
	"github.com/mattermost/mattermost/server/v8/model"
)

type configuration struct {
	Enabled          bool     `json:"enabled"`
	AllowedDurations []string `json:"allowed_durations"`
}

const (
	minServerVersion = "9.0.0"
)

func (p *Plugin) OnConfigurationChange() error {
	var configuration configuration
	if err := p.API.LoadPluginConfiguration(&configuration); err != nil {
		return err
	}

	p.configurationLock.Lock()
	defer p.configurationLock.Unlock()

	p.configuration = configuration
	return nil
}

func (p *Plugin) isDurationAllowed(duration string) bool {
	if len(p.configuration.AllowedDurations) == 0 {
		return true
	}

	for _, allowed := range p.configuration.AllowedDurations {
		if allowed == duration {
			return true
		}
	}
	return false
}

func (p *Plugin) GetPluginManifest() *model.Manifest {
	return &model.Manifest{
		Id:               "com.fambear.expiring-messages",
		Name:             "Expiring Messages",
		Version:          "0.1.0",
		MinServerVersion: minServerVersion,
	}
}
