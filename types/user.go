package types

import (
	"github.com/alpha-omega-corp/bunapp-api/httputils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"
	"os"
)

type UserRaw struct {
	bun.BaseModel     `bun:"table:users,alias:u"`
	ID                int64  `bun:"id,pk,autoincrement"`
	FirstName         string `bun:"first_name"`
	LastName          string `bun:"last_name"`
	Email             string `bun:"email,unique"`
	Age               int    `bun:"age"`
	EncryptedPassword string `bun:"encrypted_password"`
}

type User struct {
	Id                int64  `json:"id"`
	FirstName         string `json:"firstName"`
	LastName          string `json:"lastName"`
	Age               int    `json:"age"`
	Email             string `json:"email"`
	EncryptedPassword string `json:"-"`
}

func (u *User) Verify(pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(pw)) == nil
}

func (u *User) CreateToken() (string, error) {
	claims := &jwt.MapClaims{
		"expiresAt": 15000,
		"email":     u.Email,
	}

	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func (u *User) Claims(c jwt.MapClaims) error {
	if u.Email != c["email"].(string) {
		return httputils.ErrForbidden
	}

	return nil
}

func (u *UserRaw) ToUser() *User {
	return &User{
		Id:                u.ID,
		FirstName:         u.FirstName,
		LastName:          u.LastName,
		Age:               u.Age,
		Email:             u.Email,
		EncryptedPassword: u.EncryptedPassword,
	}
}

type CreateUserRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Age       int    `json:"age"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}
