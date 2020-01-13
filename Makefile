GOFLAGS ?= -mod=vendor

test:
	@echo "go test -mod vendor ./..."
	@go test ./...