package user

import "time"

type Role string

const (
	RoleAdmin Role = "ADMIN"
)

type User struct {
	ID           int64
	Login        string
	PasswordHash string
	Role         Role
	CreatedAt    time.Time
}
