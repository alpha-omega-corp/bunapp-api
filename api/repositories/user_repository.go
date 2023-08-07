package repositories

import (
	"context"
	"github.com/alpha-omega-corp/bunapp-api/api/types"
	"github.com/uptrace/bun"
)

type IUserRepository interface {
	GetByEmail(email string, ctx context.Context) (*types.User, error)
	GetById(id int64, ctx context.Context) (*types.User, error)
	GetUser(column string, value any, ctx context.Context) (*types.User, error)
	GetAll(ctx context.Context) ([]types.User, error)
	CreateUser(u *types.UserRaw, ctx context.Context) (*types.User, error)
	DeleteUser(id int64, ctx context.Context) error
}

type UserRepository struct {
	IUserRepository
	db *bun.DB
}

func NewUserRepository(db *bun.DB) IUserRepository {
	return &UserRepository{
		db: db,
	}
}

func (ur *UserRepository) GetByEmail(email string, ctx context.Context) (*types.User, error) {
	return ur.GetUser("email", email, ctx)
}

func (ur *UserRepository) GetById(id int64, ctx context.Context) (*types.User, error) {
	return ur.GetUser("id", id, ctx)
}

func (ur *UserRepository) GetUser(column string, value any, ctx context.Context) (*types.User, error) {
	var rawUser types.UserRaw
	err := ur.db.NewSelect().Model(&rawUser).Where(column+" = ?", value).Scan(ctx, &rawUser)
	if err != nil {
		return nil, err
	}

	return rawUser.ToUser(), nil
}

func (ur *UserRepository) GetAll(ctx context.Context) ([]types.User, error) {
	var users []types.UserRaw
	if err := ur.db.NewSelect().Model(&users).Scan(ctx); err != nil {
		return nil, err
	}

	us := make([]types.User, len(users))
	for i, user := range users {
		us[i] = *user.ToUser()
	}

	return us, nil
}

func (ur *UserRepository) CreateUser(u *types.UserRaw, ctx context.Context) (*types.User, error) {
	_, err := ur.db.NewInsert().Model(u).Exec(ctx)
	if err != nil {
		return nil, err
	}

	user, err := ur.GetByEmail(u.Email, ctx)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (ur *UserRepository) DeleteUser(id int64, ctx context.Context) error {
	_, err := ur.db.NewDelete().Model(&types.UserRaw{}).Where("id = ?", id).Exec(ctx)

	if err != nil {
		return err
	}

	return nil
}
