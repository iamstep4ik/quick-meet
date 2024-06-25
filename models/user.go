package models

type User struct {
	Id           int    `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	HashPassword string `json:"hash_password"`
}
