include .env

run:
	@go run ./cmd/server/main.go

migrate_create:
	@migrate create -ext sql -dir db/migrations $(filter-out $@,$(MAKECMDGOALS))

migrate_up:
	@migrate -path db/migrations -database $(DB_DSN)?sslmode=disable up

migrate_down:
	@migrate -path db/migrations -database $(DB_DSN)?sslmode=disable down

sqlc_generate:
	@sqlc generate

docker_run_db:
	@docker run -p 5432:5432 -d --name golang-final-db  -e POSTGRES_PASSWORD=Sadasa2015 -e POSTGRES_USER=kakimbekn -e POSTGRES_DB=golang-final postgres

# neccessary to use arguments in makefile
%:
	@: