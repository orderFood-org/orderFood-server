package ds

import (
	"gorm.io/gorm"
)

// type DbTime struct {
// 	CreatedAt time.Time `json:"created_at"`
// 	UpdatedAt time.Time `json:"updated_at"`
// }

// type IdType struct {
// 	Id uint64 `json:"id"`
// }

type User struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}
