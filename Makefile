server:
	go run cmd/main.go -env=dev serve

build:
	go build -o bin/api cmd/main.go

db_create:
	docker-compose up -d

db_init:
	go run cmd/main.go -env=dev db init

db_reset:
	go run cmd/main.go -env=dev db reset


