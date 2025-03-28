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


create-mig: ## create new migrations with name by $name. https://github.com/golang-migrate/migrate/tree/master/cmd/migrate
	migrate create -ext sql -dir migrations -seq $(name)
.PHONY: create-mig
