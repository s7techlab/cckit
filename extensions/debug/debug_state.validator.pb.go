// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: debug/debug_state.proto

package debug

import (
	fmt "fmt"
	math "math"
	proto "github.com/golang/protobuf/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	github_com_mwitkow_go_proto_validators "github.com/mwitkow/go-proto-validators"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

func (this *Prefix) Validate() error {
	return nil
}
func (this *Prefixes) Validate() error {
	for _, item := range this.Prefixes {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Prefixes", err)
			}
		}
	}
	return nil
}
func (this *PrefixesMatchCount) Validate() error {
	// Validation of proto3 map<> fields is unsupported.
	return nil
}
func (this *CompositeKeys) Validate() error {
	for _, item := range this.Keys {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Keys", err)
			}
		}
	}
	return nil
}
func (this *CompositeKey) Validate() error {
	return nil
}
func (this *Value) Validate() error {
	return nil
}
