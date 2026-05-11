package manga

import (
	"database/sql"
	"encoding/json"

	"mangahub/pkg/models"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository{
	return &Repository{db: db}
}

func (r *Repository) GetByID(id string)(*models.Manga, error){
	var manga models.Manga
	var genresJSON string
	
	query := `SELECT id, title, author, genres, description, total_chapters, status, cover_url`
	err := r.db.QueryRow(query,id).Scan(
		&manga.ID,
		&manga.Title,
		&manga.Author,
		&genresJSON,
		&manga.Description,
		&manga.TotalChapters,
		&manga.Status,
		&manga.CoverURL,
	)
	if err != nil {
		return nil, err;
	}

	if err = json.Unmarshal([]byte(genresJSON),&manga.Genres); err != nil{
		manga.Genres = []string{}
	} 
	return &manga, nil
}

func (r* Repository) List (filter models.SearchFilters, limit, offset int)([]models.Manga, error){
	query := `SELECT id, title, author, genres, status, total_chapters, description, cover_url FROM manga WHERE 1=1`

	var args [] interface {}

	if filter.Query != "" {
		query += ` AND title LIKE ?`
		args = append(args, "%" + filter.Query + "%")
	}

	if filter.Status != ""{
		query += ` AND status = ?`
		args = append(args, filter.Status)
	}

	query += " ORDER BY id ASC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil { 
		return nil, err 
	}
	defer rows.Close()

	var mangas []models.Manga
	for rows.Next(){
		var manga models.Manga
		var genresJSON string

		if err := rows.Scan(
			&manga.ID,
			&manga.Title,
			&manga.Author,
			&genresJSON,
			&manga.Status,
			&manga.TotalChapters,
			&manga.Description,
			&manga.CoverURL,
		); err != nil {
			continue
		}

		json.Unmarshal([]byte(genresJSON),& manga.Genres)
		mangas = append(mangas, manga)
	}

	if mangas == nil {
		mangas = []models.Manga{}
	}

	return mangas, nil
}