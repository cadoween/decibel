package seeder

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/urfave/cli/v3"
	"github.com/vingarcia/ksql"
	ksqlite "github.com/vingarcia/ksql/adapters/modernc-ksqlite"

	"github.com/cadoween/decibel/internal/spotify"
	"github.com/cadoween/decibel/pkg/iox"
)

var runFlags = []cli.Flag{
	&cli.StringFlag{
		Name:     "db",
		Usage:    "Path to the SQLite database file",
		Required: true,
	},

	&cli.StringFlag{
		Name:     "dir",
		Usage:    "Directory containing Spotify Extended Streaming History",
		Required: true,
	},

	&cli.BoolFlag{
		Name:    "verbose",
		Usage:   "Enable verbose logging",
		Value:   false,
		Aliases: []string{"v"},
	},
}

func runAction(ctx context.Context, c *cli.Command) error {
	logger := zerolog.Ctx(ctx)
	dbPath := c.String("db")
	dataDir := c.String("dir")
	verbose := c.Bool("verbose")

	if verbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	logger.Debug().
		Str("db_path", dbPath).
		Str("data_dir", dataDir).
		Msg("Initializing database connection and data import")

	db, err := ksqlite.New(ctx, dbPath, ksql.Config{})
	if err != nil {
		return fmt.Errorf("ksqlite.New: %w", err)
	}
	defer iox.Close(db, logger)

	logger.Debug().Str("path", dataDir).Msg("Reading streaming history")

	spotifyJSONReader := spotify.NewJSONReader()
	streams, err := spotifyJSONReader.ReadStreamsFromFolder(ctx, dataDir)
	if err != nil {
		return fmt.Errorf("spotifyJSONReader.ReadStreamsFromFolder: %w", err)
	}
	logger.Info().Int("count", len(streams)).Msg("Found streams in history files")

	logger.Debug().Msg("Starting data import")

	spotifySQLite := spotify.NewSQLite(db)
	if err := spotifySQLite.BulkInsertStreams(ctx, streams); err != nil {
		return fmt.Errorf("spotifySQLite.BulkInsertStreams: %w", err)
	}

	logger.Info().
		Int("total_streams", len(streams)).
		Str("database", dbPath).
		Msg("Successfully imported streams into database")

	return nil
}
