# Load only Make-related env
ifneq (,$(wildcard .env.make))
	include .env.make
	export DATABASE_URL MIGRATIONS_PATH
endif

.PHONY: migrate-create migrate-up migrate-down migrate-force migrate-version migrate-drop seed generate-docs test

# helpers
require-migrations-path:
	@test -n "$(MIGRATIONS_PATH)" || (echo "❌ MIGRATIONS_PATH is not set" && exit 1)

require-database-url:
	@test -n "$(DATABASE_URL)" || (echo "❌ DATABASE_URL is not set" && exit 1)

# migrate targets
migrate-create: require-migrations-path
	@migrate create -seq -ext sql -dir "$(MIGRATIONS_PATH)" \
	$(filter-out $@,$(MAKECMDGOALS))

migrate-up: require-migrations-path require-database-url
	@migrate -path="$(MIGRATIONS_PATH)" -database="$(DATABASE_URL)" up

migrate-down: require-migrations-path require-database-url
	@migrate -path="$(MIGRATIONS_PATH)" -database="$(DATABASE_URL)" \
	down $(filter-out $@,$(MAKECMDGOALS))

migrate-force: require-migrations-path require-database-url
	@migrate -path="$(MIGRATIONS_PATH)" -database="$(DATABASE_URL)" \
	force $(filter-out $@,$(MAKECMDGOALS))

migrate-version: require-migrations-path require-database-url
	@migrate -path="$(MIGRATIONS_PATH)" -database="$(DATABASE_URL)" version

migrate-drop: require-migrations-path require-database-url
	@migrate -path="$(MIGRATIONS_PATH)" -database="$(DATABASE_URL)" drop

# seed
seed:
	@go run ./cmd/migrate/seed/main.go

# generate docs
generate-docs:
	@swag init -d . -g ./internal/docs/swagger.go && swag fmt
# 	@swag init -d ./internal/docs,./internal,./cmd -g swagger.go --quiet && swag fmt

test:
	@go test -v ./...

# allow numeric args only (migration versions / steps)
%:
	@echo "$@" | grep -Eq '^[0-9]+$$' || (echo "❌ Unknown target: $@" && exit 1)
