package main

import (
	"fmt"
	"github.com/alpha-omega-corp/bunapp-api/api"
	"github.com/alpha-omega-corp/bunapp-api/api/types"
	"github.com/alpha-omega-corp/bunapp-api/app"
	"github.com/alpha-omega-corp/bunapp-api/app/httputils"
	"github.com/uptrace/bun/migrate"
	"github.com/urfave/cli/v2"
	"log"
	"net/http"
	"os"
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
			newDBCommand(),
		},
	}

	if err := appCli.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

var serverCommand = &cli.Command{
	Name:  "server",
	Usage: "manage application server",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "addr",
			Value: "localhost:8001",
			Usage: "serve address",
		},
	},
	Subcommands: []*cli.Command{
		{
			Name:  "start",
			Usage: "start the http server",
			Action: func(c *cli.Context) error {
				api.Bootstrap()
				ctx, appInstance, err := app.Start(c.Context, "bunapp-api", c.String("env"))

				if err != nil {
					return err
				}
				defer appInstance.Stop()

				var handler http.Handler
				handler = appInstance.Router()
				handler = httputils.ExitOnPanicHandler{Next: handler}

				srv := &http.Server{
					Addr:         c.String("addr"),
					ReadTimeout:  60 * time.Second,
					WriteTimeout: 60 * time.Second,
					IdleTimeout:  60 * time.Second,
					Handler:      handler,
				}

				go func() {
					if err := srv.ListenAndServe(); err != nil && !isServerClosed(err) {
						log.Printf("ListenAndServe failed: %s", err)
					}
				}()

				fmt.Printf("listening on http://%s\n", srv.Addr)
				fmt.Println(app.WaitExitSignal())

				return srv.Shutdown(ctx)
			},
		},
	},
}

func newDBCommand() *cli.Command {
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

					migrator := migrate.NewMigrator(appInstance.Database(), migrate.NewMigrations())
					return migrator.Init(ctx)
				},
			},
			{
				Name:  "reset",
				Usage: "migrate database",
				Action: func(c *cli.Context) error {
					ctx, appInstance, err := app.StartCLI(c)
					if err != nil {
						return err
					}
					defer appInstance.Stop()

					db := appInstance.Database()
					return db.ResetModel(ctx, (*types.UserRaw)(nil))

				},
			},
		},
	}
}

func isServerClosed(err error) bool {
	return err.Error() == "http: Server closed"
}
