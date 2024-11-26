package gormpostgres

import (
	"context"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"musiclib/internal/app/entities"
)

type GormDB struct {
	db *gorm.DB
}

// NewGormDB opens a new connection to a postgresql database.
func NewGormDB(dsn string) (GormDB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return GormDB{}, fmt.Errorf("failed to connect to database: %w", err)
	}
	return GormDB{db: db}, nil
}

// SaveSong saves a new song and returns its ID.
func (g *GormDB) SaveSong(ctx context.Context, song entities.Song) (uint64, error) {
	if err := g.db.WithContext(ctx).Create(&song).Error; err != nil {
		return 0, err
	}
	return song.ID, nil
}

// GetSongList returns list of songs.
func (g *GormDB) GetSongList(ctx context.Context, filter entities.Song, offset int, limit int) ([]entities.Song, error) {
	var songs []entities.Song
	query := g.db.WithContext(ctx).Model(&entities.Song{})

	//filter
	if filter.ID != 0 {
		query = query.Where("id = ?", filter.ID)
	}
	if filter.Song != "" {
		query = query.Where("song ILIKE ?", "%"+filter.Song+"%")
	}
	if filter.Group != "" {
		query = query.Where("group ILIKE ?", "%"+filter.Group+"%")
	}
	if filter.ReleaseDate != "" {
		query = query.Where("release_date = ?", filter.ReleaseDate)
	}

	// get songs
	err := query.Offset(offset).Limit(limit).Find(&songs).Error
	return songs, err
}

// GetSongLyrics returns song`s lyrics.
func (g *GormDB) GetSongLyrics(ctx context.Context, id int) (string, error) {
	var song entities.Song
	err := g.db.WithContext(ctx).Select("text").First(&song, id).Error
	if err != nil {
		return "", err
	}
	return song.Text, nil
}

// RemoveSong removes the song.
func (g *GormDB) RemoveSong(ctx context.Context, id int) error {
	return g.db.WithContext(ctx).Delete(&entities.Song{}, id).Error
}

// UpdateSong updates the song.
func (g *GormDB) UpdateSong(ctx context.Context, song entities.Song) error {
	return g.db.WithContext(ctx).Save(&song).Error
}
