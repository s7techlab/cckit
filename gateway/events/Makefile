
.: generate

GOFLAGS ?= -mod=vendor

test:
	@go test -ginkgo.failFast

generate: clean
	@echo "event service proto generation"
	@protoc -I=. \
	-I=../../../  \
    -I=../../vendor \
    -I=../../third_party/googleapis \
	--go_out=plugins=grpc:.  \
	--grpc-gateway_out=logtostderr=true:. \
    --swagger_out=logtostderr=true:. \
    --doc_out=markdown,event.pb.md:. \
	./*.proto

clean:
	@rm -f *.swagger.json *.pb.md *.validator.pb.go *.pb.go *.pb.gw.go *.pb.cc.go
