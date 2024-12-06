include .env

run:
	@go run ./cmd/server/main.go

migrate_create:
	@migrate create -ext sql -dir db/migrations $(filter-out $@,$(MAKECMDGOALS))

migrate_up:
	@migrate -path db/migrations -database $(DB_DSN)?sslmode=disable up $(filter-out $@,$(MAKECMDGOALS))

migrate_down:
	@migrate -path db/migrations -database $(DB_DSN)?sslmode=disable down $(filter-out $@,$(MAKECMDGOALS))

migrate_version:
	@migrate -path db/migrations -database $(DB_DSN)?sslmode=disable version

migrate_force:
	@migrate -path db/migrations -database $(DB_DSN)?sslmode=disable force $(filter-out $@,$(MAKECMDGOALS))

sqlc_generate:
	@sqlc generate

docker_create_db:
	@docker run -p 5432:5432 -d --name golang-final-db  -e POSTGRES_PASSWORD=Sadasa2015 -e POSTGRES_USER=kakimbekn -e POSTGRES_DB=golang-final postgres

docker_start_db:
	@docker start golang-final-db
	
docker_stop_db:
	@docker stop golang-final-db

# neccessary to use arguments in makefile
%:
	@: