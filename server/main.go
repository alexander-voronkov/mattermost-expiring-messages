package main

import (
	"github.com/alexander-voronkov/mattermost-expiring-messages/server"
	"github.com/mattermost/mattermost/server/v8/plugin"
)

func main() {
	plugin.Main(server.Plugin{})
}
