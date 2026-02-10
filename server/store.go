package main

// ensureIndex is a no-op for now.
// The plugin uses KV store buckets for tracking expiring posts,
// which doesn't require database indexes.
// If direct database access is needed in the future,
// consider using the Mattermost plugin API's KV store or
// requesting the feature through proper channels.
func (p *Plugin) ensureIndex() error {
	// KV store based approach doesn't need SQL indexes
	p.API.LogInfo("Expiring Messages plugin initialized")
	return nil
}
