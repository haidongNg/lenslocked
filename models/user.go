package models

import (
	"database/sql"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int
	Email        string
	PasswordHash string
}

type UserService struct {
	DB *sql.DB
}

type NewUser struct {
	Email    string
	Password string
}

func (us *UserService) Create(email string, password string) (*User, error) {
	email = strings.ToLower(email)
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("Create user: %w", err)
	}

	passwordHash := string(hashedBytes)

	user := User{
		Email:        email,
		PasswordHash: passwordHash,
	}
	// Insert
	row, err := us.DB.Exec(`INSERT INTO sys_users(email, password_hash) VALUES(?,?)`, email, passwordHash)
	if err != nil {
		return nil, fmt.Errorf("Create user: %w", err)
	}
	// Get last ID after insert
	lastId, err := row.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("Create user: %w", err)
	}
	user.ID = int(lastId)
	return &user, nil
}

func (us *UserService) Authenticate(email string, password string) (*User, error) {
	email = strings.ToLower(email)
	user := User{
		Email: email,
	}

	row := us.DB.QueryRow("SELECT id, password_hash FROM sys_users WHERE email = ?", email)

	err := row.Scan(&user.ID, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("Authenticate: %v", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("Authenticate: %w", err)
	}
	return &user, nil
}
