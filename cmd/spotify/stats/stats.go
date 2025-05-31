package stats

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/urfave/cli/v3"
	"github.com/vingarcia/ksql"
	ksqlite "github.com/vingarcia/ksql/adapters/modernc-ksqlite"

	"github.com/cadoween/decibel/internal/spotify"
	"github.com/cadoween/decibel/pkg/iox"
)

func artistsAction(ctx context.Context, c *cli.Command) error {
	logger := zerolog.Ctx(ctx)
	dbPath := c.String("db")
	verbose := c.Bool("verbose")

	if verbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	logger.Debug().
		Str("db_path", dbPath).
		Msg("Connecting to database")

	db, err := ksqlite.New(ctx, dbPath, ksql.Config{})
	if err != nil {
		return fmt.Errorf("ksqlite.New: %w", err)
	}
	defer iox.Close(db, logger)

	spotifySQLite := spotify.NewSQLite(db)

	artists, err := spotifySQLite.GetTopArtistsByPlayTime(ctx)
	if err != nil {
		return fmt.Errorf("spotifySQLite.GetTopArtistsByPlayTime: %w", err)
	}

	_, _ = fmt.Printf("\nTop Artists by Play Time:\n\n")
	_, _ = fmt.Printf("%-30s %-12s %-15s\n", "Artist", "Play Count", "Total Time")
	_, _ = fmt.Printf("%s\n", strings.Repeat("-", 60))

	for _, artist := range artists {
		duration := time.Duration(artist.TotalPlayTime) * time.Millisecond
		hours := int(duration.Hours())
		minutes := int(duration.Minutes()) % 60

		_, _ = fmt.Printf("%-30s %-12d %dh %dm\n",
			truncateString(artist.Artist, 30),
			artist.PlayCount,
			hours,
			minutes,
		)
	}

	return nil
}

func tracksAction(ctx context.Context, c *cli.Command) error {
	logger := zerolog.Ctx(ctx)
	dbPath := c.String("db")
	verbose := c.Bool("verbose")

	if verbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	logger.Debug().
		Str("db_path", dbPath).
		Msg("Connecting to database")

	db, err := ksqlite.New(ctx, dbPath, ksql.Config{})
	if err != nil {
		return fmt.Errorf("ksqlite.New: %w", err)
	}
	defer iox.Close(db, logger)

	spotifySQLite := spotify.NewSQLite(db)

	tracks, err := spotifySQLite.GetTopTracksByPlayTime(ctx)
	if err != nil {
		return fmt.Errorf("spotifySQLite.GetTopTracksByPlayTime: %w", err)
	}

	_, _ = fmt.Printf("\nTop Tracks by Play Count:\n\n")
	_, _ = fmt.Printf("%-40s %-30s %-12s %-15s\n", "Track", "Artist", "Play Count", "Total Time")
	_, _ = fmt.Printf("%s\n", strings.Repeat("-", 100))

	for _, track := range tracks {
		duration := time.Duration(track.TotalPlayTimeMS) * time.Millisecond
		hours := int(duration.Hours())
		minutes := int(duration.Minutes()) % 60

		_, _ = fmt.Printf("%-40s %-30s %-12d %dh %dm\n",
			truncateString(track.Track, 40),
			truncateString(track.Artist, 30),
			track.PlayCount,
			hours,
			minutes,
		)
	}

	return nil
}

// truncateString cuts a string if it's longer than maxLen and adds "..." at the
// end.
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
