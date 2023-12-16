package roles

import "time"

type Role struct {
	Id int64 `json:"id"`
	Name string `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func (Role) TableName() string {
	return "roles"
}