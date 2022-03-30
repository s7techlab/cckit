.: generate

generate: clean
	@echo "cflat service proto generation"
	@protoc-b2bchain -I=. \
	-I=../../../../vendor \
	-I=../../../../third_party/googleapis \
	-I=../../../../../  \
	--govalidators-b2bchain_out=. \
	--cc-gateway-b2bchain_out=paths=source_relative,embed_swagger,service_name_method_prefix,service_resolver:. \
	--go-b2bchain_out=plugins=grpc:. \
	--grpc-gateway-b2bchain_out=logtostderr=true:. \
    --swagger-b2bchain_out=logtostderr=true:. \
    --doc-b2bchain_out=../../../../docs/tools/proto-gen.tmpl,cflat.pb.md:. \
	./*.proto

clean:
	@rm -f *.pb.md *.swagger.json *.validator.pb.go *.pb.go *.pb.gw.go *.pb.cc.go