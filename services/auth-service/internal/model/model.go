package model

import "time"

type Role string

const (
	RoleStudent Role = "student"
	RoleCompany Role = "company"
	RoleAdmin   Role = "admin"
)

type User struct {
	ID           string
	Email        string
	PasswordHash string
	Role         Role
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
