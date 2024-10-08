package repository

import (
	"context"

	"github.com/DSorbon/effective-mobile-task/internal/models"
)

type Song interface {
	List(ctx context.Context, filter *models.SongFilter) (*models.SongPagination, error)
	Create(ctx context.Context, song *models.SongCreate) error
	Update(ctx context.Context, ID int64, song *models.SongUpdate) error
	Get(ctx context.Context, ID int64) (*models.Song, error)
	Delete(ctx context.Context, ID int64) error
}
