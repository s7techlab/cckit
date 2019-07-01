# Hyperledger Fabric chaincode kit (CCKit)

## Chaincode gateway
 
With gRPC we can define chaincode interface once in a .proto file and  API / SDK  will be automatically created for this chaincode.
We also get all the advantages of working with protocol buffers, including efficient serialization, a simple IDL, 
and easy interface updating.

Chaincode-as-service gateway generator allows to generate from gRPC service definition:
 
* Chaincode handlers interface 
* Chaincode gateway - service, can act as chaincode SDK or can be exposed as gRPC or REST service

### Install the generator

`GO111MODULE=on go install github.com/s7techlab/cckit/gateway/protoc-gen-cc-gateway`




