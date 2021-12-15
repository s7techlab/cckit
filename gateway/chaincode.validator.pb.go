// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: chaincode.proto

package gateway

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	_ "github.com/hyperledger/fabric-protos-go/peer"
	github_com_mwitkow_go_proto_validators "github.com/mwitkow/go-proto-validators"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

func (this *ChaincodeLocator) Validate() error {
	return nil
}
func (this *ChaincodeInput) Validate() error {
	if this.Chaincode != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Chaincode); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Chaincode", err)
		}
	}
	// Validation of proto3 map<> fields is unsupported.
	return nil
}
func (this *ChaincodeExec) Validate() error {
	if this.Input != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Input); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Input", err)
		}
	}
	return nil
}
func (this *BlockRange) Validate() error {
	return nil
}
func (this *ChaincodeEventsRequest) Validate() error {
	if this.Chaincode != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Chaincode); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Chaincode", err)
		}
	}
	if this.Block != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Block); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Block", err)
		}
	}
	return nil
}
func (this *ChaincodeInstanceInput) Validate() error {
	// Validation of proto3 map<> fields is unsupported.
	return nil
}
func (this *ChaincodeInstanceExec) Validate() error {
	if this.Input != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Input); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Input", err)
		}
	}
	return nil
}
func (this *ChaincodeInstanceEventsRequest) Validate() error {
	if this.Block != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Block); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Block", err)
		}
	}
	return nil
}
