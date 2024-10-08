package models

import "time"

type Song struct {
	ID          int64      `json:"id"`
	Artist      string     `json:"artist"`
	Group       string     `json:"group"`
	Title       string     `json:"title"`
	Lyrics      string     `json:"lyrics"`
	ReleaseDate time.Time  `json:"release_date"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

type SongCreate struct {
	Artist      string
	Group       string
	Title       string
	Lyrics      string
	ReleaseDate time.Time
}

type SongUpdate struct {
	Artist      *string
	Group       *string
	Title       *string
	Lyrics      *string
	ReleaseDate *time.Time
}

type SongFilter struct {
	Artist      string
	Group       string
	Title       string
	Lyrics      string
	ReleaseDate *time.Time
	Page        int
}

type SongPagination struct {
	Data []Song
	Page Pagination
}
