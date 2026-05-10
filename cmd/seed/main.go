// cmd/seed imports manga.json into the SQLite database.
// Run once before starting servers: go run ./cmd/seed
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"mangahub/pkg/database"
	"mangahub/pkg/models"
	"mangahub/pkg/utils"
)

func main() {

	cfg := utils.LoadConfig()

	db, err := database.New(cfg.Database.Path)
	if err != nil{
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	//read json file
	data, err := os.ReadFile(cfg.DataPath)
	if err != nil{
		log.Fatalf("Failed to Read JSON File: %v", err)
	}

	var mangas []models.Manga
	if err := json.Unmarshal(data, &mangas); err != nil {
		log.Fatalf("Failed to Unmarshal JSON File: %v", err)
	}
	insertSql := `INSERT OR IGNORE INTO manga (id, title, author, genres, status, total_chapters, description, cover_url)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	stmt, err := db.Prepare(insertSql)
	if err != nil {
		log.Fatalf("Failed to Prepare Insert Statement: %v", err)
	}
	defer stmt.Close()

	count := 0 
	for _, manga := range mangas{
		genres, _ := json.Marshal(manga.Genres)
		_, err := stmt.Exec(
			manga.ID, 
			manga.Title, 
			manga.Author, 
			string(genres), 
			manga.Status, 
			manga.TotalChapters, 
			manga.Description, 
			manga.CoverURL,
		)
		if err != nil {
			log.Fatalf("Failed to Insert Manga: %v", err)
		}else {
			count ++
		}
	}
	fmt.Printf("Successfully Insert %d Manga\n", count)

}
