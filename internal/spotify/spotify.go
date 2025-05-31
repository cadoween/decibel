package spotify

import "time"

type Stream struct {
	TS                            time.Time `json:"ts"`
	Username                      string    `json:"username"`
	Platform                      string    `json:"platform"`
	MSPlayed                      int       `json:"ms_played"`
	ConnCountry                   string    `json:"conn_country"`
	IPAddrDecrypted               string    `json:"ip_addr_decrypted"`
	UserAgentDecrypted            string    `json:"user_agent_decrypted"`
	MasterMetadataTrackName       string    `json:"master_metadata_track_name"`
	MasterMetadataAlbumArtistName string    `json:"master_metadata_album_artist_name"`
	MasterMetadataAlbumAlbumName  string    `json:"master_metadata_album_album_name"`
	SpotifyTrackURI               string    `json:"spotify_track_uri"`
	EpisodeName                   *string   `json:"episode_name"`
	EpisodeShowName               *string   `json:"episode_show_name"`
	SpotifyEpisodeURI             *string   `json:"spotify_episode_uri"`
	ReasonStart                   string    `json:"reason_start"`
	ReasonEnd                     string    `json:"reason_end"`
	Shuffle                       bool      `json:"shuffle"`
	Skipped                       bool      `json:"skipped"`
	Offline                       bool      `json:"offline"`
	OfflineTimestamp              int64     `json:"offline_timestamp"`
	IncognitoMode                 bool      `json:"incognito_mode"`
}

type ArtistStats struct {
	Artist        string `ksql:"master_metadata_album_artist_name"`
	PlayCount     int64  `ksql:"play_count"`
	TotalPlayTime int64  `ksql:"total_play_time_ms"`
}

type TrackStats struct {
	Track           string `ksql:"master_metadata_track_name"`
	Artist          string `ksql:"master_metadata_album_artist_name"`
	PlayCount       int64  `ksql:"play_count"`
	TotalPlayTimeMS int64  `ksql:"total_play_time_ms"`
}
