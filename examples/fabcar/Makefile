.: generate

generate: clean
	@echo "fabcar service proto generation"
	@protoc -I=. \
	-I=../../vendor \
	-I=../../third_party/googleapis \
	-I=../../../  \
	--govalidators_out=. \
	--cc-gateway_out=paths=source_relative,embed_swagger,service_name_method_prefix,service_resolver:. \
	--go_out=plugins=grpc:. \
	--grpc-gateway_out=logtostderr=true:. \
    --swagger_out=logtostderr=true:. \
    --doc_out=../../../../docs/tools/proto-gen.tmpl,fabcar.pb.md:. \
	./*.proto

clean:
	@rm -f *.pb.md *.swagger.json *.validator.pb.go *.pb.go *.pb.gw.go *.pb.cc.go