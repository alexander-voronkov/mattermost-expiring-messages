package server

import (
	"fmt"

	"github.com/mattermost/mattermost/server/v8/plugin"
)

func (p *Plugin) ensureIndex() error {
	query := `
		CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_posts_ttl_expires_at
		ON posts (((props->'ttl'->>'expires_at')::bigint))
		WHERE (props->'ttl'->>'enabled')::boolean = true;
	`

	_, appErr := p.API.ExecuteDirectQuery(query)
	if appErr != nil {
		return fmt.Errorf("failed to create TTL index: %w", appErr)
	}

	p.API.LogInfo("TTL index created or already exists")
	return nil
}
