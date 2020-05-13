GOFLAGS ?= -mod=vendor

test:
	@echo "go test -mod vendor ./..."
	@go test -mod=vendor ./...

refresh-deps:
	@echo "go mod tidy"
	@GOFLAGS='' GONOSUMDB=github.com/hyperledger/fabric go mod tidy
	@echo "go mod vendor"
	@GOFLAGS='' GONOSUMDB=github.com/hyperledger/fabric go mod vendor