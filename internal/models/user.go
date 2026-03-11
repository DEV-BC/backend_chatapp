package models

import (
	"database/sql"
	"errors"
	"time"

	"github.com/DEV-BC/backend_chatapp/internal/db"
)

type User struct {
	ID                   int64      `json:"id"`
	Name                 string     `json:"name"`
	Email                string     `json:"email"`
	Password             string     `json:"-"`
	RefreshTokenWeb      *string    `json:"-"`
	RefreshTokenWebAt    *time.Time `json:"-"`
	RefreshTokenMobile   *string    `json:"-"`
	RefreshTokenMobileAt *time.Time `json:"-"`
	CreatedAt            time.Time  `json:"created_at"`
}

func GetUserByEmail(email string) (*User, error) {
	u := &User{}
	row := db.DB.QueryRow(
		`SELECT id, name, email, password, refresh_token_web, refresh_token_web_at, refresh_token_mobile, refresh_token_mobile_at, created_at 
					FROM users 
					WHERE email = ?`, email,
	)
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.RefreshTokenWeb, &u.RefreshTokenWebAt, &u.RefreshTokenMobile, &u.RefreshTokenMobileAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

func CreateUserByEmail(name, email, hashedPassword string) (*User, error) {
	res, err := db.DB.Exec("INSERT INTO users (name, email, password) VALUES (?, ?, ?)", name, email, hashedPassword)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()

	createdAt := time.Now()

	return &User{
		ID:        id,
		Name:      name,
		Email:     email,
		CreatedAt: createdAt,
	}, nil
}
