package app

import (
	"chadgpt-api/types"
	"context"
	"database/sql"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"
	"os"
)

func ValidateJwt(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
}

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID                int64  `bun:"id,pk,autoincrement"`
	FirstName         string `bun:"first_name"`
	LastName          string `bun:"last_name"`
	Email             string `bun:"email,unique"`
	Age               int    `bun:"age"`
	EncryptedPassword string `bun:"encrypted_password"`
}

func (app *App) NewUser(ctx context.Context, data *types.CreateUserRequest) (sql.Result, error) {
	encPw, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return app.Database().NewInsert().Model(&User{
		FirstName:         data.FirstName,
		LastName:          data.LastName,
		Age:               data.Age,
		Email:             data.Email,
		EncryptedPassword: string(encPw),
	}).Exec(ctx)
}

func (u *User) IsValid(pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(pw)) == nil
}

func (u *User) CreateJwt() (string, error) {
	claims := &jwt.MapClaims{
		"expiresAt": 15000,
		"email":     u.Email,
	}

	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func (u *User) ToResponse() types.UserResponse {
	return types.UserResponse{
		Id:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Age:       u.Age,
		Email:     u.Email,
	}
}
