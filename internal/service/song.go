package service

import (
	"context"

	"github.com/DSorbon/effective-mobile-task/internal/models"
	"github.com/DSorbon/effective-mobile-task/internal/repository"
)

var _ Song = (*SongService)(nil)

type SongService struct {
	songRepository repository.Song
}

func NewSongService(songRepository repository.Song) *SongService {
	return &SongService{
		songRepository: songRepository,
	}
}

func (s *SongService) List(ctx context.Context, filter *models.SongFilter) (*models.SongPagination, error) {
	return s.songRepository.List(ctx, filter)
}

func (s *SongService) Create(ctx context.Context, song *models.SongCreate) error {
	return s.songRepository.Create(ctx, song)
}

func (s *SongService) Update(ctx context.Context, ID int64, song *models.SongUpdate) error {
	return s.songRepository.Update(ctx, ID, song)
}

func (s *SongService) Get(ctx context.Context, ID int64) (*models.Song, error) {
	return s.songRepository.Get(ctx, ID)
}

func (s *SongService) Delete(ctx context.Context, ID int64) error {
	return s.songRepository.Delete(ctx, ID)
}
