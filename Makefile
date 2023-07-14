run:
	go run cmd/main.go -env=dev serve

db_create:
	cd storage/docker/ && docker-compose up -d

db_init:
	go run cmd/main.go -env=dev db init

db_migrate:
	go run cmd/main.go -env=dev db migrate

build:
	go build -o storage/bin/api cmd/main.go
