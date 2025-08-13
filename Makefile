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

# Testing
.PHONY: test
test:
	go clean -testcache
	go test ./... -covermode count -coverpkg=chat/auth/internal/service/...,chat/auth/internal/api/...,chat/chat_server/internal/service/...,chat/chat_server/internal/api/... -count 5

.PHONY: test-coverage
test-coverage:
	cd auth && make test-coverage
	cd chat_server && make test-coverage
	@echo "Merging coverage reports..."
	@echo "mode: count" > coverage.out
	@grep -h -v "^mode:" auth/coverage.out chat_server/coverage.out >> coverage.out 2>/dev/null || true
	@echo "Generating combined HTML coverage report..."
	@go tool cover -html=coverage.out
	@echo "Combined coverage summary:"
	@go tool cover -func=./coverage.out | grep "total" || echo "No coverage data found"
	@grep -sqFx "/coverage.out" .gitignore || echo "/coverage.out" >> .gitignore