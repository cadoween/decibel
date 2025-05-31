package stats

import "github.com/urfave/cli/v3"

var Commands = []*cli.Command{
	{
		Name:        "artists",
		Usage:       "Get top artists by play time",
		Description: "Show your most listened artists sorted by total play time",
		Action:      artistsAction,
		Flags:       sharedFlags,
	},
	{
		Name:        "tracks",
		Usage:       "Get top tracks by play time",
		Description: "Show your most played tracks sorted by play count",
		Action:      tracksAction,
		Flags:       sharedFlags,
	},
}

var sharedFlags = []cli.Flag{
	&cli.StringFlag{
		Name:     "db",
		Usage:    "Path to the SQLite database file",
		Required: true,
	},
	&cli.BoolFlag{
		Name:    "verbose",
		Usage:   "Enable verbose logging",
		Value:   false,
		Aliases: []string{"v"},
	},
}
