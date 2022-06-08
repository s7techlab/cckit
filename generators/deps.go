//go:build tools
// +build tools

package generators

import (
	// chaincode gateway
	_ "github.com/s7techlab/cckit/gateway/protoc-gen-cc-gateway"
	// proto/grpc
	_ "github.com/golang/protobuf/protoc-gen-go"
	// json gateway
	_ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway"
	// docs
	_ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger"
	// validation schema
	_ "github.com/mwitkow/go-proto-validators/protoc-gen-govalidators"
	// protoc docs
	_ "github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc"
)
