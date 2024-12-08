include .env

run:
	@go run ./cmd/server/main.go

migrate_create:
	@migrate create -ext sql -dir db/migrations $(filter-out $@,$(MAKECMDGOALS))

migrate_up:
	@migrate -path db/migrations -database $(DB_DSN) up $(filter-out $@,$(MAKECMDGOALS))

migrate_down:
	@migrate -path db/migrations -database $(DB_DSN) down $(filter-out $@,$(MAKECMDGOALS))

migrate_version:
	@migrate -path db/migrations -database $(DB_DSN) version

migrate_force:
	@migrate -path db/migrations -database $(DB_DSN) force $(filter-out $@,$(MAKECMDGOALS))

sqlc_generate:
	@sqlc generate

docker_create_db:
	@docker run -p 5432:5432 -d --name golang-final  -e POSTGRES_PASSWORD=Sadasa2015 -e POSTGRES_USER=kakimbekn -e POSTGRES_DB=golang-final postgres

docker_start_db:
	@docker start golang-final-db
	
docker_stop_db:
	@docker stop golang-final-db

docker-compose_up:
	@docker-compose up -d --build

docker-compose_stop:
	@docker-compose stop

docker-build:
	@docker build -t e-commerce-gin .

# neccessary to use arguments in makefile
%:
	@: