package spotify

import "time"

type Stream struct {
	TS                            time.Time `json:"ts"`
	SpotifyEpisodeURI             *string   `json:"spotify_episode_uri"`
	EpisodeShowName               *string   `json:"episode_show_name"`
	EpisodeName                   *string   `json:"episode_name"`
	SpotifyTrackURI               string    `json:"spotify_track_uri"`
	Username                      string    `json:"username"`
	UserAgentDecrypted            string    `json:"user_agent_decrypted"`
	MasterMetadataTrackName       string    `json:"master_metadata_track_name"`
	MasterMetadataAlbumArtistName string    `json:"master_metadata_album_artist_name"`
	MasterMetadataAlbumAlbumName  string    `json:"master_metadata_album_album_name"`
	ConnCountry                   string    `json:"conn_country"`
	ReasonEnd                     string    `json:"reason_end"`
	Platform                      string    `json:"platform"`
	IPAddrDecrypted               string    `json:"ip_addr_decrypted"`
	ReasonStart                   string    `json:"reason_start"`
	MSPlayed                      int       `json:"ms_played"`
	OfflineTimestamp              int64     `json:"offline_timestamp"`
	Shuffle                       bool      `json:"shuffle"`
	Skipped                       bool      `json:"skipped"`
	Offline                       bool      `json:"offline"`
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
