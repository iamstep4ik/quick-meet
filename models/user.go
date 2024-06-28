package models

import "time"

type User struct {
	Id           int       `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	Email        string    `json:"email" db:"email"`
	HashPassword string    `json:"hash_password" db:"password_hash"`
	InsertedAt   time.Time `json:"inserted_at" db:"inserted_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}
