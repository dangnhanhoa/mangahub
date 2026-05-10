package auth

import (
	"database/sql"
	"errors"
	"time"
	"log"

	"mangahub/pkg/models"
	"mangahub/pkg/utils"	

	"golang.org/x/crypto/bcrypt"
)

type Service struct{
	db *sql.DB

}

func NewService(db *sql.DB) *Service {
	return &Service{db:db}
}

func (s *Service) Register(username, password string) (*models.User, error){

	var exitID string
	err := s.db.QueryRow("SELECT id FROM users WHERE username = ?",username).Scan(&exitID)
	if err == nil {
		log.Printf("[AUTH] Error: %s", err)
		return nil, errors.New("Username already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("[AUTH] Error: %s", err)
		return nil, errors.New("Fail to create account")
	}
	user := models.User{
		Username: username,
		PasswordHash: string(hashedPassword),
		CreatedAt: time.Now(),
	}

	_, err = s.db.Exec("INSERT INTO users (id, username, password_hash, created_at) VALUES(?,?,?,?)",
		user.ID, user.Username, user.PasswordHash, user.CreatedAt)
	if err != nil{
		log.Printf("[AUTH] Error: %s", err)
		return nil, errors.New("Fail to create account")
	}

	return &user, nil
}

func (s *Service ) Login(username, password string)(string, error){
	var user models.User

	err := s.db.QueryRow("SELECT id, username, password_hash FROM users WHERE username = ?",
		username).Scan(&user.ID, &user.Username, &user.PasswordHash)
	if err != nil {
		log.Printf("[AUTH] Error: %s", err)
		return "", errors.New("Fail to Login: Invalid username")
	}

	if err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	); err != nil {
		log.Printf("[AUTH] Error: %s", err)
		return "", errors.New("Fail to Login: Wrong password")
	}

	cfg := utils.LoadConfig()
	token, err := utils.GenerateToken(
		user.ID,
		user.Username,
		cfg.JWT.Secret,
	)
	if err != nil {
		log.Printf("[AUTH] Error: %s", err)
		return "", errors.New("Fail to Login: Internal server error")
	}

	return token, nil
}