run:
	go run cmd/main.go -env=dev serve

init-migrations:
	go run cmd/main.go -env=dev db init

run-migrations:
	go run cmd/main.go -env=dev db migrate

create-migration:
	go run cmd/main.go -env=dev db create_go $(name)



test:
	go test -v ./...

build:
	go build -o bin/api
