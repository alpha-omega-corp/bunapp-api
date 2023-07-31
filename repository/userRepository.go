package repository

import (
	"chadgpt-api/app"
	"context"
	"github.com/uptrace/bun"
)

type UserRepository struct {
	db *bun.DB
}

func NewUserRepository(db *bun.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (ur *UserRepository) GetUserByEmail(email string, ctx context.Context) (*app.User, error) {
	return ur.GetUser(ctx, "email", email)
}

func (ur *UserRepository) GetUser(ctx context.Context, column string, value any) (*app.User, error) {
	var user app.User
	err := ur.db.NewSelect().Model(&user).Where(column+" = ?", value).Scan(ctx, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
