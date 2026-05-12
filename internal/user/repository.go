package user
import (
	"database/sql"
	"time"
	"mangahub/pkg/models"
)
type Repository struct {
	db *sql.DB
}
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}
func (r *Repository) AddOrUpdate(progress *models.UserProgress) error {
	query := `
		INSERT OR REPLACE INTO user_progress (user_id, manga_id, current_chapter, status, updated_at) 
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := r.db.Exec(query, 
		progress.UserID, 
		progress.MangaID, 
		progress.CurrentChapter, 
		progress.Status, 
		time.Now())
		
	return err
}
func (r *Repository) Get(userID, mangaID string) (*models.UserProgress, error) {
	var p models.UserProgress
	query := `SELECT user_id, manga_id, current_chapter, status, updated_at FROM user_progress WHERE user_id = ? AND manga_id = ?`
	
	err := r.db.QueryRow(query, userID, mangaID).Scan(
		&p.UserID, &p.MangaID, &p.CurrentChapter, &p.Status, &p.UpdatedAt)
		
	if err != nil {
		return nil, err 
	}
	return &p, nil
}
func (r *Repository) List(userID string, status string) ([]models.UserProgress, error) {
	query := `SELECT user_id, manga_id, current_chapter, status, updated_at FROM user_progress WHERE user_id = ?`
	var args []interface{}
	args = append(args, userID)
	if status != "" {
		query += ` AND status = ?`
		args = append(args, status)
	}
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.UserProgress
	for rows.Next() {
		var p models.UserProgress
		if err := rows.Scan(&p.UserID, &p.MangaID, &p.CurrentChapter, &p.Status, &p.UpdatedAt); err == nil {
			list = append(list, p)
		}
	}
	
	if list == nil {
		list = []models.UserProgress{} 
	}
	return list, nil
}
func (r *Repository) Delete(userID, mangaID string) error {
	_, err := r.db.Exec(`DELETE FROM user_progress WHERE user_id = ? AND manga_id = ?`, userID, mangaID)
	return err
}