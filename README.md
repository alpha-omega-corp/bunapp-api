# BUNAPP-API

This repository serves as an example of a RESTful API written in `golang`

## MODULES
- [bunorm](https://github.com/uptrace/bun) for database connections
- [bunrouter](https://github.com/uptrace/bunrouterbun) for http routing
- [golang-jwt/v5](https://github.com/golang-jwt/jwt) for authentication
- [urfavcli/v2](https://github.com/urfave/cli) for command line interface integration


## SETUP ENVIRONMENT

Install `golang` and `make` binaries.
```bash
sudo apt update && sudo apt upgrade
sudo apt install golang makefile
```

Pull the code to local directory `~/GolandProjects/bunapp-api`.
```bash
git clone https://github.com/alpha-omega-corp/bunapp-api.git ~/GolandProjects/bunapp-api
```

Open the project in `Goland` and install dependencies.
```bash
go mod download
go mod tidy
```

Create a `dev` config file from the `test` config file in `app/config`.
```bash
cp ./app/config/test.yaml ./app/config/dev.yaml
```

Create a `dev` database from the `docker-compose` file.
```bash
make db_create
```

Create migrations table for `bun`.
```bash
make db_init
```

Run the database reset command once to create required tables.
```bash
make db_reset
```

Start the application using the `Makefile`.
```bash
make server
```

Or manually using `go run`.
```bash
go run cmd/main.go -env=dev server
```


## FEATURES

### 1. Command line interface

[Documentation](https://cli.urfave.org/v2/getting-started/) | [GitHub](https://github.com/urfave/cli)


`main.go` handles the application's `cli` commands.

**Top levels commands:**
```bash
go run cmd/main.go -env=dev server
go run cmd/main.go -env=dev db
```

**Database commands:**
```bash
go run cmd/main.go -env=dev db init
go run cmd/main.go -env=dev db reset
```

**Server commands:**
```bash
go run cmd/main.go -env=dev server start
```

### 2. Bun application package

[Documentation](https://bun.uptrace.dev/guide/starter-kit.html) | [GitHub](https://github.com/go-bun/bun-starter-kit)

The `app` directory contains the code responsible for creating an `app *App` instance.

1. Create a directory `example` inside `github.com/alpha-omega-corp/bunapp-api`

2. From `example` create a `init.go` file with the following code:

```go
package example

import (
  "context"
  "github.com/alpha-omega-corp/bunapp-api/app"
)

func Bootstrap() {
  app.OnStart("example.init", func (ctx context.Context, app *app.App) error {
    //app.Router()
    //app.ApiRouter()
    //app.Database()
    //app.Repositories()
    //app.GptClient()
    //app.PromptManager()
  })
}
```
<u>**Example**</u>: [api-app](https://github.com/alpha-omega-corp/bunapp-api/blob/production/api/init.go)

The `callback` function gives you access to the [`app *App`](https://github.com/alpha-omega-corp/bunapp-api/blob/production/app/app.go) instance

To start the `example` application call the public `Bootstrap` function from the `cmd/main.go` file.

```golang
package main

import (
    "github.com/alpha-omega-corp/bunapp-api/example"
)

func main() {
	//...
    Action:
        func(c *cli.Context) error {
            example.Bootstrap()
        } //...
}

```
<u>**Example**</u>: [cmd/main.go](https://github.com/alpha-omega-corp/bunapp-api/blob/production/cmd/main.go)








