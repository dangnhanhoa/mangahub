package models

type Manga struct {
	ID            string   `json:"id" db:"id"`
	Title         string   `json:"title" db:"title"`
	Author        string   `json:"author" db:"author"`
	Genres        []string `json:"genres"`
	Status        string   `json:"status" db:"status"`
	TotalChapters int      `json:"total_chapters" db:"total_chapters"`
	Description   string   `json:"description" db:"description"`
	CoverURL      string   `json:"cover_url,omitempty" db:"cover_url"`
}

type SearchFilters struct {
	Query       string   `form:"q"`
	Genres      []string `form:"genres"`
    Status      string   `form:"status"`
    YearRange   [2]int   `form:"year_range"`
    Rating      float64  `form:"rating"`
    SortBy      string   `form:"sort_by"` 
}
