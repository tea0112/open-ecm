package users

import (
	"open-ecm/roles"
	"time"
)

type User struct {
	Id        *int64     `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Password  string     `json:"password"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"-"`
	Role      roles.Role
}

func (u User) TableName() string {
	return "users"
}
