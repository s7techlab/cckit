.: generate

generate:
	@echo "chaincode gateway service"
	@protoc -I=. \
	 -I=../../  \
	 -I=../vendor \
	 -I=../third_party/googleapis \
	 -I=../third_party/hyperledger \
	 --grpc-gateway-b2bchain_out=paths=source_relative:. \
     --go-b2bchain_out=paths=source_relative,plugins=grpc:.  \
     --govalidators-b2bchain_out=paths=source_relative:. \
     --swagger-b2bchain_out=logtostderr=true:. \
	 chaincode.proto