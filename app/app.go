package app

import (
	"context"
	"database/sql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"github.com/uptrace/bunrouter"
	"github.com/urfave/cli/v2"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

type appCtxKey struct{}

func ContextWithApp(ctx context.Context, app *App) context.Context {
	ctx = context.WithValue(ctx, appCtxKey{}, app)
	return ctx
}

type App struct {
	ctx    context.Context
	config *Config

	stopping uint32
	stopCh   chan struct{}

	onStop      appHooks
	onAfterStop appHooks

	router    *bunrouter.Router
	apiRouter *bunrouter.Group

	// lazy init
	dbOnce sync.Once
	db     *bun.DB
}

func New(ctx context.Context, c *Config) *App {
	app := &App{
		config: c,
		stopCh: make(chan struct{}),
	}
	app.ctx = ContextWithApp(ctx, app)
	app.initRouter()
	return app
}

func StartCLI(c *cli.Context) (context.Context, *App, error) {
	return Start(c.Context, c.Command.Name, c.String("env"))
}

func Start(ctx context.Context, service, envName string) (context.Context, *App, error) {
	cfg, err := ReadConfig(FS(), service, envName)
	if err != nil {
		return nil, nil, err
	}
	return StartConfig(ctx, cfg)
}

func StartConfig(ctx context.Context, cfg *Config) (context.Context, *App, error) {
	rand.Seed(time.Now().UnixNano())

	app := New(ctx, cfg)
	if err := onStart.Run(ctx, app); err != nil {
		return nil, nil, err
	}
	return app.ctx, app, nil
}

func (app *App) Stop() {
	_ = app.onStop.Run(app.ctx, app)
	_ = app.onAfterStop.Run(app.ctx, app)
}

func (app *App) OnStop(name string, fn HookFunc) {
	app.onStop.Add(newHook(name, fn))
}

func (app *App) OnAfterStop(name string, fn HookFunc) {
	app.onAfterStop.Add(newHook(name, fn))
}

func (app *App) Context() context.Context {
	return app.ctx
}

func (app *App) Config() *Config {
	return app.config
}

func (app *App) Running() bool {
	return !app.Stopping()
}

func (app *App) Stopping() bool {
	return atomic.LoadUint32(&app.stopping) == 1
}

func (app *App) IsDebug() bool {
	return app.config.Debug
}

func (app *App) Router() *bunrouter.Router {
	return app.router
}

func (app *App) ApiRouter() *bunrouter.Group {
	return app.apiRouter
}

func (app *App) Database() *bun.DB {

	app.dbOnce.Do(func() {
		dbConf := app.config.DB
		driverOptions := pgdriver.NewConnector(
			pgdriver.WithAddr(dbConf.ADDR),
			pgdriver.WithDatabase(dbConf.NAME),
			pgdriver.WithUser(dbConf.USER),
			pgdriver.WithPassword(dbConf.PASS),
			pgdriver.WithTLSConfig(nil))

		conn := sql.OpenDB(driverOptions)
		db := bun.NewDB(conn, pgdialect.New())

		db.AddQueryHook(bundebug.NewQueryHook(
			bundebug.WithEnabled(true),
			bundebug.WithVerbose(true),
		))

		app.OnStop("db.Close", func(ctx context.Context, _ *App) error {
			return db.Close()
		})

		app.db = db
	})

	return app.db
}

func WaitExitSignal() os.Signal {
	ch := make(chan os.Signal, 3)
	signal.Notify(
		ch,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)
	return <-ch
}
