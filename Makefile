main:
	go run main.go server.go

docker:
	docker run --name petitions -p 5555:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:16-alpine

create_db:
	psql -h localhost -p 5555 -U root -c "CREATE DATABASE petitions_db;"

mgup:
	migrate -database "postgres://root:secret@localhost:5555/petitions_db?sslmode=disable" -path ./schema up

mgdown:
	migrate -database "postgres://root:secret@localhost:5555/petitions_db?sslmode=disable" -path ./schema down -force

sqlc:
	sqlc generate

Phony:main, docker, create_db, mgup, mgdown, sqlc