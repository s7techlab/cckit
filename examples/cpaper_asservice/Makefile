.: generate

generate:
	@protoc --version
	@echo "commercial paper schema proto generation"
	@protoc -I=./schema/ \
	-I=../../vendor \
	--go_out=./schema/    \
	--govalidators_out=./schema/ \
    --doc_out=./schema/ --doc_opt=markdown,schema.md \
	./schema/schema.proto

	@echo "commercial paper service proto generation"
	@protoc -I=./service/ \
	-I=../../../../../ \
	-I=../../vendor \
	-I=../../third_party/googleapis \
	--go_out=plugins=grpc:./service/    \
 	--cc-gateway_out=logtostderr=true:./service/ \
	--grpc-gateway_out=logtostderr=true:./service/ \
    --swagger_out=logtostderr=true:./service/ \
    --doc_out=./service/ --doc_opt=markdown,service.md \
	./service/service.proto