package postgres

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/DSorbon/effective-mobile-task/internal/models"
	"github.com/DSorbon/effective-mobile-task/internal/repository"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ repository.Song = (*SongRepository)(nil)

type SongRepository struct {
	pool *pgxpool.Pool
}

func NewSongRepository(pool *pgxpool.Pool) *SongRepository {
	return &SongRepository{
		pool: pool,
	}
}

func (r *SongRepository) List(ctx context.Context, filter *models.SongFilter) (*models.SongPagination, error) {
	var (
		limit  uint64 = 10
		offset uint64 = limit * uint64(filter.Page-1)
	)

	builderSelect := sq.Select().
		From("songs").
		PlaceholderFormat(sq.Dollar).
		Limit(limit).
		Offset(offset)

	if filter.Artist != "" {
		artist := fmt.Sprintf("%%%s%%", filter.Artist)
		builderSelect = builderSelect.Where(sq.Like{"artist": artist})
	}

	if filter.Group != "" {
		group := fmt.Sprintf("%%%s%%", filter.Group)
		builderSelect = builderSelect.Where(sq.Like{"s_group": group})
	}

	if filter.Title != "" {
		title := fmt.Sprintf("%%%s%%", filter.Title)
		builderSelect = builderSelect.Where(sq.Like{"title": title})
	}

	if filter.ReleaseDate != nil {
		builderSelect = builderSelect.Where(sq.Eq{"release_date": filter.ReleaseDate})
	}

	pagination, err := r.pagination(ctx, builderSelect, int(limit), filter.Page)
	if err != nil {
		return nil, err
	}

	builderSelect = builderSelect.Columns("id", "artist", "s_group", "title", "lyrics", "release_date", "created_at", "updated_at")

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	record := models.Song{}
	songs := []models.Song{}

	for rows.Next() {
		err := rows.Scan(&record.ID, &record.Artist, &record.Group, &record.Title, &record.Lyrics, &record.ReleaseDate, &record.CreatedAt, &record.UpdatedAt)
		if err != nil {
			return nil, err
		}

		songs = append(songs, record)
	}

	return &models.SongPagination{
		Data: songs,
		Page: *pagination,
	}, nil
}

func (r *SongRepository) Create(ctx context.Context, song *models.SongCreate) error {
	builderInsert := sq.Insert("songs").
		PlaceholderFormat(sq.Dollar).
		Columns("artist", "s_group", "title", "lyrics", "release_date").
		Values(song.Artist, song.Group, song.Title, song.Lyrics, song.ReleaseDate)

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return err
	}

	_, err = r.pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *SongRepository) Update(ctx context.Context, ID int64, song *models.SongUpdate) error {
	builderUpdate := sq.Update("songs").PlaceholderFormat(sq.Dollar)

	if song.Artist != nil {
		builderUpdate = builderUpdate.Set("artist", song.Artist)
	}

	if song.Group != nil {
		builderUpdate = builderUpdate.Set("s_group", song.Group)
	}

	if song.Title != nil {
		builderUpdate = builderUpdate.Set("title", song.Title)
	}

	if song.Lyrics != nil {
		builderUpdate = builderUpdate.Set("lyrics", song.Lyrics)
	}

	if song.ReleaseDate != nil {
		builderUpdate = builderUpdate.Set("release_date", song.ReleaseDate)
	}

	builderUpdate = builderUpdate.Set("updated_at", time.Now()).Where(sq.Eq{"id": ID})

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		return err
	}

	_, err = r.pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *SongRepository) Get(ctx context.Context, ID int64) (*models.Song, error) {
	builderSelectOne := sq.Select("id", "artist", "s_group", "title", "lyrics", "release_date", "created_at", "updated_at").
		From("songs").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": ID}).
		Limit(1)

	query, args, err := builderSelectOne.ToSql()
	if err != nil {
		return nil, err
	}

	var song models.Song
	err = r.pool.QueryRow(ctx, query, args...).Scan(&song.ID, &song.Artist, &song.Group, &song.Title, &song.Lyrics, &song.ReleaseDate, &song.CreatedAt, &song.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &song, nil
}

func (r *SongRepository) Delete(ctx context.Context, ID int64) error {
	builderInsert := sq.Delete("songs").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": ID})

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return err
	}

	_, err = r.pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *SongRepository) pagination(ctx context.Context, builder sq.SelectBuilder, limit, page int) (*models.Pagination, error) {
	var (
		tmpl        = models.Pagination{}
		recordCount int
	)

	builderCount := builder.Columns("count(id)")

	query, args, err := builderCount.ToSql()
	if err != nil {
		return nil, err
	}

	r.pool.QueryRow(ctx, query, args...).Scan(&recordCount)

	tmpl.TotalPage = int(math.Ceil(float64(recordCount) / float64(limit)))
	tmpl.CurrentPage = page
	tmpl.RecordPerPage = limit

	// Calculator the Next/Previous Page
	if page <= 0 {
		tmpl.Next = page + 1
	} else if page < tmpl.TotalPage {
		tmpl.Previous = page - 1
		tmpl.Next = page + 1
	} else if page == tmpl.TotalPage {
		tmpl.Previous = page - 1
		tmpl.Next = 0
	}

	return &tmpl, nil
}
