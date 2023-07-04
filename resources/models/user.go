package models

import (
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID       int64  `bun:"id,pk"`
	Name     string `bun:"first_name"`
	LastName string `bun:"last_name"`
	Age      int    `bun:"age"`
	Email    string `bun:"email"`
	Password string `bun:"password"`
}
