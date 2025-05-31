package spotify_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/cadoween/decibel/internal/spotify"
	"github.com/cadoween/decibel/internal/spotify/ksqltest"
)

func TestSQLite_BulkInsertStreams(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		streams []spotify.Stream
		mock    func(*ksqltest.MockProvider)
		wantErr error
	}{
		{
			name: "successfully inserts streams in batches",
			streams: []spotify.Stream{
				{
					TS:                            time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					Username:                      "user1",
					Platform:                      "platform1",
					MSPlayed:                      1000,
					ConnCountry:                   "US",
					IPAddrDecrypted:               "1.1.1.1",
					UserAgentDecrypted:            "agent1",
					MasterMetadataTrackName:       "track1",
					MasterMetadataAlbumArtistName: "artist1",
					MasterMetadataAlbumAlbumName:  "album1",
					SpotifyTrackURI:               "uri1",
				},
				{
					TS:       time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
					Username: "user2",
				},
			},
			mock: func(m *ksqltest.MockProvider) {
				m.EXPECT().Exec(gomock.Any(), gomock.Any()).Return(nil, nil)
				m.EXPECT().Exec(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
			},
			wantErr: nil,
		},
		{
			name:    "handles create table error",
			streams: []spotify.Stream{},
			mock: func(m *ksqltest.MockProvider) {
				m.EXPECT().Exec(gomock.Any(), gomock.Any()).Return(nil, assert.AnError)
			},
			wantErr: assert.AnError,
		},
		{
			name: "handles batch insert error",
			streams: []spotify.Stream{{
				TS:       time.Now(),
				Username: "user1",
			}},
			mock: func(m *ksqltest.MockProvider) {
				m.EXPECT().Exec(gomock.Any(), gomock.Any()).Return(nil, nil)
				m.EXPECT().Exec(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, assert.AnError)
			},
			wantErr: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockProvider := ksqltest.NewMockProvider(ctrl)
			tt.mock(mockProvider)

			sqlite := spotify.NewSQLite(mockProvider)
			err := sqlite.BulkInsertStreams(context.Background(), tt.streams)

			if tt.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestSQLite_GetTopArtistsByPlayTime(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	tests := []struct {
		name    string
		mock    func(*ksqltest.MockProvider)
		want    []spotify.ArtistStats
		wantErr error
	}{
		{
			name: "successfully retrieves top artists",
			mock: func(m *ksqltest.MockProvider) {
				expected := []spotify.ArtistStats{
					{Artist: "Artist1", PlayCount: 100, TotalPlayTime: 5000},
					{Artist: "Artist2", PlayCount: 80, TotalPlayTime: 4000},
				}
				m.EXPECT().Query(ctx, gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, records any, _ string, _ ...any) error {
						*(records.(*[]spotify.ArtistStats)) = expected
						return nil
					})
			},
			want: []spotify.ArtistStats{
				{Artist: "Artist1", PlayCount: 100, TotalPlayTime: 5000},
				{Artist: "Artist2", PlayCount: 80, TotalPlayTime: 4000},
			},
		},
		{
			name: "handles query error",
			mock: func(m *ksqltest.MockProvider) {
				m.EXPECT().Query(ctx, gomock.Any(), gomock.Any()).Return(assert.AnError)
			},
			wantErr: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockProvider := ksqltest.NewMockProvider(ctrl)
			tt.mock(mockProvider)

			sqlite := spotify.NewSQLite(mockProvider)
			got, err := sqlite.GetTopArtistsByPlayTime(ctx)

			if tt.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestSQLite_GetMostSkippedTracks(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	tests := []struct {
		name    string
		mock    func(*ksqltest.MockProvider)
		want    []spotify.TrackSkipStats
		wantErr error
	}{
		{
			name: "successfully retrieves most skipped tracks",
			mock: func(m *ksqltest.MockProvider) {
				expected := []spotify.TrackSkipStats{
					{TrackName: "Track1", ArtistName: "Artist1", SkipCount: 8, SkipRate: 0.8},
					{TrackName: "Track2", ArtistName: "Artist2", SkipCount: 6, SkipRate: 0.6},
				}
				m.EXPECT().Query(ctx, gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, records any, _ string, _ ...any) error {
						*(records.(*[]spotify.TrackSkipStats)) = expected
						return nil
					})
			},
			want: []spotify.TrackSkipStats{
				{TrackName: "Track1", ArtistName: "Artist1", SkipCount: 8, SkipRate: 0.8},
				{TrackName: "Track2", ArtistName: "Artist2", SkipCount: 6, SkipRate: 0.6},
			},
		},
		{
			name: "handles query error",
			mock: func(m *ksqltest.MockProvider) {
				m.EXPECT().Query(ctx, gomock.Any(), gomock.Any()).Return(assert.AnError)
			},
			wantErr: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockProvider := ksqltest.NewMockProvider(ctrl)
			tt.mock(mockProvider)

			sqlite := spotify.NewSQLite(mockProvider)
			got, err := sqlite.GetMostSkippedTracks(ctx)

			if tt.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
