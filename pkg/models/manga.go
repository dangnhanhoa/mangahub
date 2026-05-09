package models

type Manga struct {
	ID            string   `json:"id" db:"id"`
	Title         string   `json:"title" db:"title"`
	Author        string   `json:"author" db:"author"`
	Genres        []string `json:"genres"`      // serialized as JSON text in SQLite
	Status        string   `json:"status" db:"status"` // "ongoing" | "completed"
	TotalChapters int      `json:"total_chapters" db:"total_chapters"`
	Description   string   `json:"description" db:"description"`
	CoverURL      string   `json:"cover_url,omitempty" db:"cover_url"`
}

type SearchFilters struct {
	Query  string
	Genre  string
	Status string
	Page   int
	Limit  int
}
