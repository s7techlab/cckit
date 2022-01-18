// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: chaincode.proto

package gateway

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	_ "github.com/hyperledger/fabric-protos-go/peer"
	_ "github.com/mwitkow/go-proto-validators"
	github_com_mwitkow_go_proto_validators "github.com/mwitkow/go-proto-validators"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

func (this *ChaincodeLocator) Validate() error {
	if this.Chaincode == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("Chaincode", fmt.Errorf(`value '%v' must not be an empty string`, this.Chaincode))
	}
	if this.Channel == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("Channel", fmt.Errorf(`value '%v' must not be an empty string`, this.Channel))
	}
	return nil
}
func (this *ChaincodeInput) Validate() error {
	if nil == this.Chaincode {
		return github_com_mwitkow_go_proto_validators.FieldError("Chaincode", fmt.Errorf("message must exist"))
	}
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
func (this *BlockLimit) Validate() error {
	return nil
}
func (this *ChaincodeEventsStreamRequest) Validate() error {
	if nil == this.Chaincode {
		return github_com_mwitkow_go_proto_validators.FieldError("Chaincode", fmt.Errorf("message must exist"))
	}
	if this.Chaincode != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Chaincode); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Chaincode", err)
		}
	}
	if this.FromBlock != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.FromBlock); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("FromBlock", err)
		}
	}
	if this.ToBlock != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.ToBlock); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("ToBlock", err)
		}
	}
	return nil
}
func (this *ChaincodeEventsRequest) Validate() error {
	if nil == this.Chaincode {
		return github_com_mwitkow_go_proto_validators.FieldError("Chaincode", fmt.Errorf("message must exist"))
	}
	if this.Chaincode != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Chaincode); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Chaincode", err)
		}
	}
	if this.FromBlock != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.FromBlock); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("FromBlock", err)
		}
	}
	if this.ToBlock != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.ToBlock); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("ToBlock", err)
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
func (this *ChaincodeInstanceEventsStreamRequest) Validate() error {
	if this.FromBlock != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.FromBlock); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("FromBlock", err)
		}
	}
	if this.ToBlock != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.ToBlock); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("ToBlock", err)
		}
	}
	return nil
}
func (this *ChaincodeInstanceEventsRequest) Validate() error {
	if this.FromBlock != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.FromBlock); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("FromBlock", err)
		}
	}
	if this.ToBlock != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.ToBlock); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("ToBlock", err)
		}
	}
	return nil
}
func (this *ChaincodeEvents) Validate() error {
	if this.Chaincode != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Chaincode); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Chaincode", err)
		}
	}
	if this.FromBlock != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.FromBlock); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("FromBlock", err)
		}
	}
	if this.ToBlock != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.ToBlock); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("ToBlock", err)
		}
	}
	for _, item := range this.Items {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Items", err)
			}
		}
	}
	return nil
}
func (this *RawJson) Validate() error {
	return nil
}
func (this *ChaincodeEvent) Validate() error {
	if this.Event != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Event); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Event", err)
		}
	}
	if this.Payload != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Payload); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Payload", err)
		}
	}
	return nil
}
