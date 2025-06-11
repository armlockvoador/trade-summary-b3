gen-mock:
	@echo "clearing existing mocks..."
	@rm -rf ./mocks/*
	@echo "generating project mocks..."
	mockery --all


gen-db:
	sqlc generate

gen-db-deps:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

tidy:
	go mod tidy -compat=1.22

lint:
	golangci-lint run -v
