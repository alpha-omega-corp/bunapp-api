package app

import (
	"github.com/alpha-omega-corp/bunapp-api/repository"
	"github.com/uptrace/bun"
)

type RepoManager struct {
	userRepo repository.IUserRepository
}

func (app *App) initRepositories() {
	app.repoManager = NewRepoManager(app.Database())
}

func NewRepoManager(db *bun.DB) *RepoManager {
	return &RepoManager{
		userRepo: repository.NewUserRepository(db),
	}
}

func (r *RepoManager) User() repository.IUserRepository {
	return r.userRepo
}
