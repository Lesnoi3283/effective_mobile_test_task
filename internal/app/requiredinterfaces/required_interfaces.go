package requiredinterfaces

import (
	"context"
	"musiclib/internal/app/entities"
)

//go:generate mockgen -source=required_interfaces.go -destination=./mocks/mocks.go -package=mocks

type ExtraDataProvider interface {
	GetExtraSongData(song entities.Song) (releaseDate, text, link string, err error)
}

type SongStorage interface {
	SaveSong(ctx context.Context, song entities.Song) (id uint64, err error)
	GetSongList(ctx context.Context, filter entities.Song, offset int, limit int) ([]entities.Song, error)
	GetSongLyrics(ctx context.Context, id uint64) (string, error)
	RemoveSong(ctx context.Context, id uint64) error
	UpdateSong(ctx context.Context, song entities.Song) error
}
