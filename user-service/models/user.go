package models

type User struct {
	ID        int64  `json:"id"`
	Name      string `json:"name" form:"name"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}
