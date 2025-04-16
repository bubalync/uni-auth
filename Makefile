# todo заменить ямл на env
#include .env
#export

LOCAL_BIN:=$(CURDIR)/bin
PROTO_DIR=api/uni-auth-proto/proto/auth/v1
OUT_PROTO_GEN_DIR=internal/proto/v1

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
.PHONY: migrate-create

mockgen: ### generate mock
	mockgen -source=internal/service/service.go -destination=internal/mocks/servicemocks/service.go -package=servicemocks
	mockgen -source=pkg/hasher/password.go      -destination=internal/mocks/utilmocks/hasher.go     -package=utilmocks
	mockgen -source=internal/lib/jwtgen/jwt.go  -destination=internal/mocks/utilmocks/jwt.go        -package=utilmocks
	mockgen -source=pkg/redis/redis.go          -destination=internal/mocks/redismocks/redis.go     -package=redismocks
	mockgen -source=internal/repo/repo.go       -destination=internal/mocks/repomocks/repo.go       -package=repomocks
.PHONY: mockgen

test: ### run test
	go test -v ./...
.PHONY: test

cover-html: ### run test with coverage and open html report
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out
.PHONY: cover-html

cover: ### run test with coverage
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	rm coverage.out
.PHONY: cover

#bin-deps: ### install tools
#    GOBIN=$(LOCAL_BIN) go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
#    GOBIN=$(LOCAL_BIN) go install github.com/golang/mock/mockgen@latest
#    GOBIN=$(LOCAL_BIN) go install go.uber.org/mock/mockgen@latest
#    GOBIN=$(LOCAL_BIN) go install github.com/swaggo/swag/cmd/swag@latest
#    GOBIN=$(LOCAL_BIN) go install github.com/daixiang0/gci@latest
#    GOBIN=$(LOCAL_BIN) go install mvdan.cc/gofumpt@latest
#    GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
#    GOBIN=$(LOCAL_BIN) go install golang.org/x/vuln/cmd/govulncheck@latest
#.PHONY: bin-deps

protoc: # gen from submodule: uni-auth-proto
	protoc \
		--proto_path=$(PROTO_DIR) \
		--go_out=$(OUT_PROTO_GEN_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(OUT_PROTO_GEN_DIR) --go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/auth.proto
.PHONY: protoc