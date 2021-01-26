GOFLAGS ?= -mod=vendor
PROTO_PACKAGES := examples/cpaper_asservice/schema examples/cpaper_asservice/service examples/cpaper_extended examples/payment/schema gateway/events gateway/service state/

test:
	@echo "go test -mod vendor ./..."
	@go test ./...

refresh-deps:
	@echo "go mod tidy"
	@GOFLAGS='' GONOSUMDB=github.com/hyperledger/fabric go mod tidy
	@echo "go mod vendor"
	@GOFLAGS='' GONOSUMDB=github.com/hyperledger/fabric go mod vendor

proto:
	@for pkg in $(PROTO_PACKAGES) ;do echo $$pkg && buf generate --path $$pkg -o $$(echo $$pkg | cut -d "/" -f1); done