package spotify

import (
	"context"
	"fmt"
	"slices"
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

	insertQuery := `
		INSERT INTO spotify_streams (
			ts, username, platform, ms_played, conn_country, ip_addr_decrypted,
			user_agent_decrypted, master_metadata_track_name, master_metadata_album_artist_name,
			master_metadata_album_album_name, spotify_track_uri, episode_name,
			episode_show_name, spotify_episode_uri, reason_start, reason_end,
			shuffle, skipped, offline, offline_timestamp, incognito_mode
		) VALUES 
	`

	// SQLite has a limit of 999 parameters, each stream has 21 parameters so
	// we'll use batches of 45 rows (945 parameters) to stay safely under the
	// limit.
	for batch := range slices.Chunk(streams, 45) {
		args := make([]any, 0, len(batch)*21)
		valueStrings := make([]string, len(batch))

		for j := range batch {
			valueStrings[j] = "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
			args = append(args,
				batch[j].TS, batch[j].Username, batch[j].Platform, batch[j].MSPlayed, batch[j].ConnCountry,
				batch[j].IPAddrDecrypted, batch[j].UserAgentDecrypted, batch[j].MasterMetadataTrackName,
				batch[j].MasterMetadataAlbumArtistName, batch[j].MasterMetadataAlbumAlbumName,
				batch[j].SpotifyTrackURI, batch[j].EpisodeName, batch[j].EpisodeShowName,
				batch[j].SpotifyEpisodeURI, batch[j].ReasonStart, batch[j].ReasonEnd, batch[j].Shuffle,
				batch[j].Skipped, batch[j].Offline, batch[j].OfflineTimestamp, batch[j].IncognitoMode,
			)
		}

		batchQuery := insertQuery + strings.Join(valueStrings, ",")
		if _, err := s.sqlProvider.Exec(ctx, batchQuery, args...); err != nil {
			return fmt.Errorf("s.sqlProvider.Exec: %w", err)
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
