package app

import (
	"github.com/alpha-omega-corp/bunapp-api/api/repositories"
	"github.com/uptrace/bun"
)

type RepoManager struct {
	userRepo repositories.IUserRepository
}

func (app *App) initRepositories() {
	app.repoManager = NewRepoManager(app.Database())
}

func NewRepoManager(db *bun.DB) *RepoManager {
	return &RepoManager{
		userRepo: repositories.NewUserRepository(db),
	}
}

func (r *RepoManager) User() repositories.IUserRepository {
	return r.userRepo
}
