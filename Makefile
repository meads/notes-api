# GO_BUILD_ENV := CGO_ENABLED=0 GOOS=linux GOARCH=amd64
# DOCKER_BUILD=$(shell pwd)/.docker_build
# DOCKER_CMD=$(DOCKER_BUILD)/notes-api
DATABASE_URL := $(shell grep -iR '^DATABASE_URL*' .env | cut -d= -f2-)

sqlc:
	@sqlc version
	@sqlc compile
	@sqlc generate

tidy:
	@go mod tidy

mockgen:
	@mockgen -package repository -destination ./internal/repository/querier_mock.go github.com/meads/notes-api/internal/db/sqlc Querier
	@mockgen -package security -destination ./internal/security/bcrypt_mock.go github.com/meads/notes-api/internal/security Hasher
	@mockgen -package security -destination ./internal/security/claims_mock.go github.com/meads/notes-api/internal/security Tokener

	@mockgen -package service -destination ./internal/service/note_repository_mock.go github.com/meads/notes-api/internal/domain NoteRepository
	@mockgen -package service -destination ./internal/service/session_repository_mock.go github.com/meads/notes-api/internal/domain SessionRepository
	@mockgen -package service -destination ./internal/service/user_repository_mock.go  github.com/meads/notes-api/internal/domain UserRepository

	@mockgen -package handler -destination ./internal/handler/note_service_mock.go github.com/meads/notes-api/internal/handler NoteServicer
	@mockgen -package handler -destination ./internal/handler/user_service_mock.go github.com/meads/notes-api/internal/handler UserServicer
	@mockgen -package handler -destination ./internal/handler/auth_service_mock.go github.com/meads/notes-api/internal/handler AuthServicer

generate: tidy sqlc mockgen

test:
	@go test -v ./... \
	| sed ''/PASS/s//$$(printf "\033[32mPASS\033[0m")/g'' | sed ''/FAIL/s//$$(printf "\033[31mFAIL\033[0m")/g''

test-coverage:
	@go test -v -coverprofile cover.out ./...
	@go tool cover -html=cover.out

verify: generate test

migrate-drop-recreate:
	@migrate -path db/migration -database $(DATABASE_URL) drop
	@migrate -path db/migration -database $(DATABASE_URL) up

local-db-shell:
	@docker exec -it firstly-api-db-1 /bin/bash

local-db-psql:
	@docker exec -it firstly-api-db-1 psql $(DATABASE_URL)

# scale-zero:
# 	@heroku ps:scale api=0

# login:
# 	@heroku login
# 	@heroku container:login

# login-docker:
# 	@docker login --username=$(DOCKER_USERNAME) --password=$$(heroku auth:token) registry.heroku.com

# clean:
# 	@rm -rf $(DOCKER_BUILD)
# 	@mkdir -p $(DOCKER_BUILD)

# build: clean
# 	$(GO_BUILD_ENV) go build -v -o $(DOCKER_CMD) .

deploy:
	@git push origin main
# 	@heroku container:push api
# 	@heroku container:release api
