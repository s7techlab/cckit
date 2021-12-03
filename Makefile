GOFLAGS ?= -mod=vendor

PROTO_PACKAGES_GO := examples/cpaper_extended/schema state/schema examples/cpaper_asservice/schema
PROTO_PACKAGES_CCGW := extensions/debug extensions/owner examples/cpaper_asservice/service

test:
	@echo "go test -mod vendor ./..."
	@go test ./...

refresh-deps:
	@echo "go mod tidy"
	@GOFLAGS='' GONOSUMDB=github.com/hyperledger/fabric go mod tidy
	@echo "go mod vendor"
	@GOFLAGS='' GONOSUMDB=github.com/hyperledger/fabric go mod vendor

proto:
	@for pkg in $(PROTO_PACKAGES_GO) ;do echo $$pkg && buf generate --template buf.gen.go.yaml --path $$pkg -o ./$$(echo $$pkg | cut -d "/" -f1); done
	@for pkg in $(PROTO_PACKAGES_CCGW) ;do echo $$pkg && buf generate --template buf.gen.gw.yaml --path $$pkg -o ./$$(echo $$pkg | cut -d "/" -f1); done

clean:
	@for pkg in $(PROTO_PACKAGES_GO); do find $$pkg \( -name '*.pb.go' -or -name '*.pb.cc.go' -or -name '*.pb.gw.go' -or -name '*.swagger.json' \) -delete;done
	@for pkg in $(PROTO_PACKAGES_CCGW); do find $$pkg \( -name '*.pb.go' -or -name '*.pb.cc.go' -or -name '*.pb.gw.go' -or -name '*.swagger.json' \) -delete;done