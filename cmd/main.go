package main

import (
	"context"
	"os"

	"github.com/rs/zerolog"
	"github.com/urfave/cli/v3"

	"github.com/cadoween/decibel/cmd/spotify"
)

func main() {
	cmd := &cli.Command{
		Name:        "decibel",
		Usage:       "Analyze and manage your music listening history",
		Description: "A command-line tool for processing and analyzing music streaming history data, providing insights into your listening habits across different platforms.",
		Commands: []*cli.Command{
			{
				Name:        "spotify",
				Usage:       "Analyze and manage your Spotify listening history",
				Description: "A command-line tool for processing and analyzing Spotify streaming history data, providing insights into your listening habits.",
				Commands:    spotify.Commands,
			},
		},
	}

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	ctx := context.Background()

	if err := cmd.Run(logger.WithContext(ctx), os.Args); err != nil {
		logger.Error().Err(err).Msg("Failed to run decibel command")
	}
}
