package request

type SongCreate struct {
	Artist      string `json:"artist" validate:"required"`
	Group       string `json:"group" validate:"required"`
	Title       string `json:"title" validate:"required"`
	Lyrics      string `json:"lyrics" validate:"required"`
	ReleaseDate string `json:"release_date" validate:"required|date|gt:2"`
}

type SongUpdate struct {
	Artist      *string `json:"artist,omitempty"`
	Group       *string `json:"group,omitempty"`
	Title       *string `json:"title,omitempty"`
	Lyrics      *string `json:"lyrics,omitempty"`
	ReleaseDate *string `json:"release_date,omitempty" validate:"date"`
}
