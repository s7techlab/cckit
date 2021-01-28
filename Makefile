GOFLAGS ?= -mod=vendor
PROTO_PACKAGES_CC := examples/cpaper_asservice/service
PROTO_PACKAGES_GW := examples/cpaper_asservice/service examples/cpaper_asservice/schema examples/cpaper_extended/schema examples/payment/schema gateway/events gateway/service state/schema


test:
	@echo "go test -mod vendor ./..."
	@go test ./...

refresh-deps:
	@echo "go mod tidy"
	@GOFLAGS='' GONOSUMDB=github.com/hyperledger/fabric go mod tidy
	@echo "go mod vendor"
	@GOFLAGS='' GONOSUMDB=github.com/hyperledger/fabric go mod vendor

proto:
	@for pkg in $(PROTO_PACKAGES_GW) ;do echo $$pkg && buf generate -v --template buf.gen.gw.yaml --path $$pkg -o ./$$(echo $$pkg | cut -d "/" -f1); done
	@for pkg in $(PROTO_PACKAGES_CC) ;do echo $$pkg && buf generate --template buf.gen.cc.yaml --path $$pkg -o ./$$pkg; done