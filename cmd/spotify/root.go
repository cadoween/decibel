package spotify

import (
	"github.com/urfave/cli/v3"

	"github.com/cadoween/decibel/cmd/spotify/seeder"
	"github.com/cadoween/decibel/cmd/spotify/stats"
)

var Commands = []*cli.Command{
	{
		Name:        "seeder",
		Usage:       "Spotify streaming history data seeder",
		Description: "Seed your Spotify streaming history data in the local database",
		Commands:    seeder.Commands,
	},

	{
		Name:        "stats",
		Usage:       "View Spotify listening statistics",
		Description: "Analyze your Spotify listening history and view various statistics",
		Commands:    stats.Commands,
	},
}
