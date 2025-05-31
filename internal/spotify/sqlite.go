package spotify

import (
	"context"
	"fmt"
	"strings"

	"github.com/vingarcia/ksql"
)

type SQLite struct {
	sqlProvider ksql.Provider
}

func NewSQLite(sqlProvider ksql.Provider) *SQLite {
	return &SQLite{
		sqlProvider: sqlProvider,
	}
}

func (s *SQLite) BulkInsertStreams(ctx context.Context, streams []Stream) error {
	createTableQuery := `
		CREATE TABLE IF NOT EXISTS spotify_streams (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			ts TIMESTAMP,
			username TEXT,
			platform TEXT,
			ms_played INTEGER,
			conn_country TEXT,
			ip_addr_decrypted TEXT,
			user_agent_decrypted TEXT,
			master_metadata_track_name TEXT,
			master_metadata_album_artist_name TEXT,
			master_metadata_album_album_name TEXT,
			spotify_track_uri TEXT,
			episode_name TEXT,
			episode_show_name TEXT,
			spotify_episode_uri TEXT,
			reason_start TEXT,
			reason_end TEXT,
			shuffle BOOLEAN,
			skipped BOOLEAN,
			offline BOOLEAN,
			offline_timestamp INTEGER,
			incognito_mode BOOLEAN
		)
	`

	if _, err := s.sqlProvider.Exec(ctx, createTableQuery); err != nil {
		return fmt.Errorf("s.sqlProvider.Exec: %w", err)
	}

	// Each stream has 21 parameters, and SQLite has a limit of 999 parameters
	// So we'll use batches of 45 rows (945 parameters) to stay safely under the
	// limit.
	batchSize := 45
	total := len(streams)

	insertQuery := `
		INSERT INTO spotify_streams (
			ts, username, platform, ms_played, conn_country, ip_addr_decrypted,
			user_agent_decrypted, master_metadata_track_name, master_metadata_album_artist_name,
			master_metadata_album_album_name, spotify_track_uri, episode_name,
			episode_show_name, spotify_episode_uri, reason_start, reason_end,
			shuffle, skipped, offline, offline_timestamp, incognito_mode
		) VALUES 
	`

	for i := 0; i < total; i += batchSize {
		end := min(i+batchSize, total)
		batch := streams[i:end]
		args := make([]any, 0, len(batch)*21)
		valueStrings := make([]string, len(batch))

		for j, stream := range batch {
			valueStrings[j] = "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
			args = append(args,
				stream.TS, stream.Username, stream.Platform, stream.MSPlayed, stream.ConnCountry, stream.IPAddrDecrypted,
				stream.UserAgentDecrypted, stream.MasterMetadataTrackName, stream.MasterMetadataAlbumArtistName,
				stream.MasterMetadataAlbumAlbumName, stream.SpotifyTrackURI, stream.EpisodeName,
				stream.EpisodeShowName, stream.SpotifyEpisodeURI, stream.ReasonStart, stream.ReasonEnd,
				stream.Shuffle, stream.Skipped, stream.Offline, stream.OfflineTimestamp, stream.IncognitoMode,
			)
		}

		batchQuery := insertQuery + strings.Join(valueStrings, ",")
		if _, err := s.sqlProvider.Exec(ctx, batchQuery, args...); err != nil {
			return fmt.Errorf("failed to insert batch %d-%d: %w", i, end, err)
		}
	}

	return nil
}

func (s *SQLite) GetTopArtistsByPlayTime(ctx context.Context) ([]ArtistStats, error) {
	query := `
		SELECT 
			master_metadata_album_artist_name,
			COUNT(*) as play_count,
			SUM(ms_played) as total_play_time_ms
		FROM spotify_streams
		WHERE master_metadata_album_artist_name IS NOT NULL
		GROUP BY master_metadata_album_artist_name
		ORDER BY total_play_time_ms DESC
		LIMIT 10
	`

	var results []ArtistStats
	if err := s.sqlProvider.Query(ctx, &results, query); err != nil {
		return nil, fmt.Errorf("s.sqlProvider.Query: %w", err)
	}

	return results, nil
}

func (s *SQLite) GetTopTracksByPlayTime(ctx context.Context) ([]TrackStats, error) {
	query := `
		SELECT
			master_metadata_track_name,
			master_metadata_album_artist_name,
			COUNT(*) AS play_count,
			SUM(ms_played) AS total_play_time_ms
		FROM spotify_streams
		WHERE master_metadata_track_name IS NOT NULL
		GROUP BY master_metadata_track_name, master_metadata_album_artist_name
		ORDER BY total_play_time_ms DESC
		LIMIT 10
	`

	var results []TrackStats
	if err := s.sqlProvider.Query(ctx, &results, query); err != nil {
		return nil, fmt.Errorf("s.sqlProvider.Query: %w", err)
	}

	return results, nil
}
