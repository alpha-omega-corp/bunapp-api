package main

import (
	"chadgpt-api/app"
	"chadgpt-api/app/httputils"
	"chadgpt-api/cmd/migrations"
	"chadgpt-api/resources/handlers"
	"chadgpt-api/resources/models"
	"fmt"
	"github.com/uptrace/bun/migrate"
	"github.com/urfave/cli/v2"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	appCli := &cli.App{
		Name:  "app",
		Usage: "bootstrap the application",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "env",
				Value: "dev",
				Usage: "environment",
			},
		},
		Commands: []*cli.Command{
			serverCommand,
			newDBCommand(migrations.Migrations),
		},
	}

	if err := appCli.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

var serverCommand = &cli.Command{
	Name:  "serve",
	Usage: "start http server",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "addr",
			Value: "localhost:8000",
			Usage: "serve address",
		},
	},
	Action: func(c *cli.Context) error {
		ctx, appInstance, err := app.Start(c.Context, "resources", c.String("env"))
		if err != nil {
			return err
		}
		defer appInstance.Stop()

		var handler http.Handler
		handler = appInstance.Router()
		handler = httputils.ExitOnPanicHandler{Next: handler}

		srv := &http.Server{
			Addr:         c.String("addr"),
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  60 * time.Second,
			Handler:      handler,
		}

		api := appInstance.ApiRouter()
		userHandler := handlers.NewUserHandler(appInstance)

		api.GET("/users", userHandler.List)
		api.GET("/users/:id", userHandler.Get)

		go func() {
			if err := srv.ListenAndServe(); err != nil && !isServerClosed(err) {
				log.Printf("ListenAndServe failed: %s", err)
			}
		}()

		fmt.Printf("listening on http://%s\n", srv.Addr)
		fmt.Println(app.WaitExitSignal())

		return srv.Shutdown(ctx)
	},
}

func newDBCommand(migrations *migrate.Migrations) *cli.Command {
	return &cli.Command{
		Name:  "db",
		Usage: "manage database migrations",
		Subcommands: []*cli.Command{
			{
				Name:  "init",
				Usage: "create migration tables",
				Action: func(c *cli.Context) error {
					ctx, appInstance, err := app.StartCLI(c)
					if err != nil {
						return err
					}
					defer appInstance.Stop()

					migrator := migrate.NewMigrator(appInstance.Database(), migrations)
					return migrator.Init(ctx)
				},
			},
			{
				Name:  "migrate",
				Usage: "migrate database",
				Action: func(c *cli.Context) error {
					ctx, appInstance, err := app.StartCLI(c)
					if err != nil {
						return err
					}
					defer appInstance.Stop()

					return appInstance.Database().ResetModel(ctx,
						(*models.User)(nil),
					)
				},
			},
			{
				Name:  "create_go",
				Usage: "create Go migration",
				Action: func(c *cli.Context) error {
					ctx, appInstance, err := app.StartCLI(c)
					if err != nil {
						return err
					}
					defer appInstance.Stop()

					migrator := migrate.NewMigrator(appInstance.Database(), migrations)

					name := strings.Join(c.Args().Slice(), "_")
					mf, err := migrator.CreateGoMigration(ctx, name)
					if err != nil {
						return err
					}
					fmt.Printf("created migration %s (%s)\n", mf.Name, mf.Path)

					return nil
				},
			},
		},
	}
}

func isServerClosed(err error) bool {
	return err.Error() == "http: Server closed"
}
