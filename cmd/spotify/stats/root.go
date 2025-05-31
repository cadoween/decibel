package stats

import "github.com/urfave/cli/v3"

var Commands = []*cli.Command{
	{
		Name:        "top-artists",
		Usage:       "Get top artists by play time",
		Description: "Show your most listened artists sorted by total play time",
		Action:      topArtistsAction,
		Flags:       sharedFlags,
	},
	{
		Name:        "top-tracks",
		Usage:       "Get top tracks by play time",
		Description: "Show your most played tracks sorted by play count",
		Action:      topTracksAction,
		Flags:       sharedFlags,
	},
	{
		Name:        "top-albums",
		Usage:       "Get top albums by play count",
		Description: "Show your most played albums sorted by play count",
		Action:      topAlbumsAction,
		Flags:       sharedFlags,
	},
	{
		Name:        "most-skipped-tracks",
		Usage:       "Get most skipped tracks",
		Description: "Show tracks that are most frequently skipped (minimum 5 plays)",
		Action:      mostSkippedTracksAction,
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
