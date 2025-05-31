package seeder

import "github.com/urfave/cli/v3"

var Commands = []*cli.Command{
	{
		Name:        "run",
		Usage:       "Run Spotify streaming history data seeder",
		Description: "Reads Spotify Extended Streaming History JSON files and seeds them into the SQLite database",
		Action:      runAction,
		Flags:       runFlags,
	},
}
