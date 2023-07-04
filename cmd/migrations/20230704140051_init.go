package migrations

import (
	"chadgpt-api/app"
	"chadgpt-api/resources/models"
	"context"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dbfixture"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		db.RegisterModel(ctx, (*models.User)(nil))
		fixture := dbfixture.New(db, dbfixture.WithRecreateTables())
		return fixture.Load(ctx, app.FS(), "fixture.yaml")
	}, func(ctx context.Context, db *bun.DB) error {
		return db.ResetModel(ctx, (*models.User)(nil))
	})
}
