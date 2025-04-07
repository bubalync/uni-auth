# todo заменить ямл на env
#include .env
#export

run-local:
	go run cmd/app/main.go --config ./config/local.yaml
.PHONY: local

migrate-up:
	go run cmd/migrator/migrator.go --databaseURL "postgres://postgres:postgres@localhost:5432/uni-auth"

swag: ## swag init
	swag init -g internal/api/http/router.go
.PHONY: swag

migrate-create: ## create new migrations with name by $name. https://github.com/golang-migrate/migrate/tree/master/cmd/migrate
	migrate create -ext sql -dir migrations -seq $(name)
.PHONY: create-mig

mockgen: ### generate mock
	mockgen -source=internal/service/service.go -destination=internal/mocks/servicemocks/service.go -package=servicemocks
	mockgen -source=pkg/hasher/password.go      -destination=internal/mocks/utilmocks/hasher.go     -package=utilmocks
	mockgen -source=internal/repo/repo.go       -destination=internal/mocks/repomocks/repo.go       -package=repomocks
.PHONY: mockgen

test: ### run test
	go test -v ./...
.PHONY: test

cover-html: ### run test with coverage and open html report
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out
.PHONY: coverage-html

cover: ### run test with coverage
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	rm coverage.out
.PHONY: coverage

