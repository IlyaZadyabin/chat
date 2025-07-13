include docker-compose.env
export
	
LOCAL_BIN:=$(CURDIR)/bin

# Database
db-up:
	docker-compose up -d

db-down:
	docker-compose down

# Migrations
install-deps:
	cd auth && make install-deps
	cd chat_server && make install-deps

local-migration-up:
	cd auth && make local-migration-up
	cd chat_server && make local-migration-up

local-migration-down:
	cd auth && make local-migration-down
	cd chat_server && make local-migration-down

local-migration-status:
	cd auth && make local-migration-status
	cd chat_server && make local-migration-status

# Service migrations
local-migration-auth:
	cd auth && make local-migration-up

local-migration-chat:
	cd chat_server && make local-migration-up

# Linter commands
install-golangci-lint:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0

lint:
	$(LOCAL_BIN)/golangci-lint run ./... --config .golangci.pipeline.yaml