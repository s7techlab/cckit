.: generate

generate:
	@echo "chaincode invoke service"
	@protoc -I=. \
	 -I=../../../  \
	 -I=../../vendor \
	 -I=../../third_party/googleapis \
	 --go_out=plugins=grpc:.  \
	 --govalidators_out=./ \
	 --grpc-gateway_out=logtostderr=true:. \
	 chaincode.proto