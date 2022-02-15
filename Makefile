GOFLAGS ?= -mod=vendor

PROTO_PACKAGES_GO := state
PROTO_PACKAGES_GW := gateway
PROTO_PACKAGES_CC_WITHSERVICE_PREFIX := extensions
PROTO_PACKAGES_CC := examples

test:
	@echo "go test -mod vendor ./..."
	@go test ./...

refresh-deps:
	@echo "go mod tidy"
	@GOFLAGS='' GONOSUMDB=github.com/hyperledger/fabric go mod tidy
	@echo "go mod vendor"
	@GOFLAGS='' GONOSUMDB=github.com/hyperledger/fabric go mod vendor

proto: clean
	@for pkg in $(PROTO_PACKAGES_CC) ;do echo $$pkg && buf generate --template buf.gen.cc.yaml $$pkg -o ./$$(echo $$pkg | cut -d "/" -f1); done
	@for pkg in $(PROTO_PACKAGES_CC_WITHSERVICE_PREFIX) ;do echo $$pkg && buf generate --template buf.gen.cc-with-service-prefix.yaml $$pkg -o ./$$(echo $$pkg | cut -d "/" -f1); done
	@for pkg in $(PROTO_PACKAGES_GW) ;do echo $$pkg && buf generate --template buf.gen.gw.yaml $$pkg -o ./$$(echo $$pkg | cut -d "/" -f1); done
	@for pkg in $(PROTO_PACKAGES_GO) ;do echo $$pkg && buf generate --template buf.gen.go.yaml $$pkg -o ./$$(echo $$pkg | cut -d "/" -f1); done
clean:
	@for pkg in $(PROTO_PACKAGES_CC); do find $$pkg \( -name '*.pb.go' -or -name '*.pb.cc.go' -or -name '*.pb.gw.go' -or -name '*.swagger.json' -or -name '*.pb.md' \) -delete;done
	@for pkg in $(PROTO_PACKAGES_CC_WITHSERVICE_PREFIX); do find $$pkg \( -name '*.pb.go' -or -name '*.pb.cc.go' -or -name '*.pb.gw.go' -or -name '*.swagger.json' -or -name '*.pb.md' \) -delete;done
	@for pkg in $(PROTO_PACKAGES_GW); do find $$pkg \( -name '*.pb.go' -or -name '*.pb.gw.go' -or -name '*.swagger.json' -or -name '*.pb.md' \) -delete;done
	@for pkg in $(PROTO_PACKAGES_GO); do find $$pkg \( -name '*.pb.go' -or -name '*.pb.md' \) -delete;done