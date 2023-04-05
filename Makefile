GOFILES := $(shell find . -name '*.go' -not -path "./vendor/*" -not -path "./tests/*" | egrep -v "^\./\.go" | grep -v _test.go)
DEPS_HASHICORP = hashicorp hashicorp-agent
DEPS_HASHICORP_TLS = hashicorp-tls hashicorp-agent-tls
DEPS_POSTGRES = postgres
DEPS_POSTGRES_TLS = postgres-ssl
PACKAGES ?= $(shell go list ./... | egrep -v "tests|e2e|mocks|mock" )
KEY_MANAGER_SERVICES = key-manager

UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Linux)
	OPEN = xdg-open
endif
ifeq ($(UNAME_S),Darwin)
	OPEN = open
endif

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

.PHONY: all lint lint-ci integration-tests swagger-tool

hashicorp:
	@docker-compose -f deps/hashicorp/docker-compose.yml up --build -d $(DEPS_HASHICORP)
	@sleep 2 # Sleep couple seconds to wait token to be created

hashicorp-tls:
	@docker-compose -f deps/hashicorp/docker-compose.yml up --build -d $(DEPS_HASHICORP_TLS)
	@sleep 2 # Sleep couple seconds to wait token to be created

hashicorp-down:
	@docker-compose -f deps/hashicorp/docker-compose.yml down --volumes --timeout 0

networks:
	@docker network create --driver=bridge hashicorp || true
	@docker network create --driver=bridge --subnet=172.16.237.0/24 besu || true
	@docker network create --driver=bridge --subnet=172.16.238.0/24 quorum || true

down-networks:
	@docker network rm quorum || true
	@docker network rm besu || true
	@docker network rm hashicorp || true

postgres:
	@docker-compose -f deps/docker-compose.yml up -d $(DEPS_POSTGRES)

postgres-tls:
	@docker-compose -f deps/docker-compose.yml up -d $(DEPS_POSTGRES_TLS)

postgres-down:
	@docker-compose -f deps/docker-compose.yml down --volumes --timeout 0

deps: networks hashicorp postgres

deps-tls: networks hashicorp-tls postgres-tls

down-deps: postgres-down hashicorp-down down-networks

gobuild:
	@GOOS=linux GOARCH=amd64 go build -o ./build/bin/key-manager

gobuild-dbg:
	CGO_ENABLED=1 go build -gcflags=all="-N -l" -i -o ./build/bin/key-manager

qkm: gobuild
	@docker-compose -f ./docker-compose.dev.yml up --force-recreate --build -d $(KEY_MANAGER_SERVICES)

dev: deps qkm

up-tls: deps-tls gobuild
	@docker-compose -f ./docker-compose.dev.yml up --build -d $(KEY_MANAGER_SERVICES)

down-dev:
	@docker-compose -f ./docker-compose.dev.yml down --volumes --timeout 0

run: gobuild
	@./build/bin/key-manager run

run-dbg: gobuild-dbg
	@dlv --listen=:2345 --headless=true --api-version=2 --accept-multiclient exec ./build/bin/key-manager run

sync: gobuild
	@docker-compose -f ./docker-compose.dev.yml up sync

deploy-remote-env:
	@bash ./scripts/deploy-remote-env.sh

run-server-dev:
	go run main.go run
