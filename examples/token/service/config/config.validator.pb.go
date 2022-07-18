// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: token/service/config/config.proto

package config

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	_ "github.com/mwitkow/go-proto-validators"
	github_com_mwitkow_go_proto_validators "github.com/mwitkow/go-proto-validators"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	_ "google.golang.org/protobuf/types/known/emptypb"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

func (this *CreateTokenTypeRequest) Validate() error {
	if this.Name == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("Name", fmt.Errorf(`value '%v' must not be an empty string`, this.Name))
	}
	if this.Symbol == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("Symbol", fmt.Errorf(`value '%v' must not be an empty string`, this.Symbol))
	}
	if !(this.Decimals < 9) {
		return github_com_mwitkow_go_proto_validators.FieldError("Decimals", fmt.Errorf(`value '%v' must be less than '9'`, this.Decimals))
	}
	if _, ok := TokenGroupType_name[int32(this.GroupType)]; !ok {
		return github_com_mwitkow_go_proto_validators.FieldError("GroupType", fmt.Errorf(`value '%v' must be a valid TokenGroupType field`, this.GroupType))
	}
	for _, item := range this.Meta {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Meta", err)
			}
		}
	}
	return nil
}
func (this *UpdateTokenTypeRequest) Validate() error {
	if this.Name == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("Name", fmt.Errorf(`value '%v' must not be an empty string`, this.Name))
	}
	if this.Symbol == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("Symbol", fmt.Errorf(`value '%v' must not be an empty string`, this.Symbol))
	}
	for _, item := range this.Meta {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Meta", err)
			}
		}
	}
	return nil
}
func (this *CreateTokenGroupRequest) Validate() error {
	if len(this.Name) < 1 {
		return github_com_mwitkow_go_proto_validators.FieldError("Name", fmt.Errorf(`value '%v' must contain at least 1 elements`, this.Name))
	}
	if this.TokenType == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("TokenType", fmt.Errorf(`value '%v' must not be an empty string`, this.TokenType))
	}
	for _, item := range this.Meta {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Meta", err)
			}
		}
	}
	return nil
}
func (this *Config) Validate() error {
	return nil
}
func (this *TokenId) Validate() error {
	return nil
}
func (this *TokenTypeId) Validate() error {
	return nil
}
func (this *TokenType) Validate() error {
	for _, item := range this.Meta {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Meta", err)
			}
		}
	}
	return nil
}
func (this *TokenTypes) Validate() error {
	for _, item := range this.Types {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Types", err)
			}
		}
	}
	return nil
}
func (this *TokenGroupId) Validate() error {
	return nil
}
func (this *TokenGroup) Validate() error {
	for _, item := range this.Meta {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Meta", err)
			}
		}
	}
	return nil
}
func (this *TokenGroups) Validate() error {
	for _, item := range this.Groups {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Groups", err)
			}
		}
	}
	return nil
}
func (this *TokenMetaRequest) Validate() error {
	if this.Key == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("Key", fmt.Errorf(`value '%v' must not be an empty string`, this.Key))
	}
	if this.Value == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("Value", fmt.Errorf(`value '%v' must not be an empty string`, this.Value))
	}
	return nil
}
func (this *TokenMeta) Validate() error {
	return nil
}
func (this *Token) Validate() error {
	if this.Type != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Type); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Type", err)
		}
	}
	if this.Group != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Group); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Group", err)
		}
	}
	return nil
}
func (this *TokenTypeCreated) Validate() error {
	return nil
}
func (this *TokenGroupCreated) Validate() error {
	return nil
}
