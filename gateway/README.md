# Hyperledger Fabric chaincode kit (CCKit)

## Chaincode-as-service gateway generator

With gRPC we can define chaincode interface once in a .proto file and  API / SDK  will be automatically created for this chaincode.
We also get all the advantages of working with protocol buffers, including efficient serialization, a simple IDL, 
and easy interface updating.

Chaincode-as-service gateway generator allows to generate from gRPC service definition:
 
* Chaincode handlers interface 
* Chaincode gateway - service, can act as chaincode SDK or can be exposed as gRPC or REST service

### Install the generator

`GO111MODULE=on go install github.com/s7techlab/cckit/gateway/protoc-gen-cc-gateway`


## Example

### Commercial paper chaincode 

#### Data model

[schema.proto](../examples/cpaper_asservice/schema/schema.proto)

```proto
syntax = "proto3";

package schema;

import "google/protobuf/timestamp.proto";
import "github.com/mwitkow/go-proto-validators/validator.proto";

// Commercial Paper state entry
message CommercialPaper {

    enum State {
        ISSUED = 0;
        TRADING = 1;
        REDEEMED = 2;
    }

    // Issuer and Paper number comprises composite primary key of Commercial paper entry
    string issuer = 1;
    string paper_number = 2;

    string owner = 3;
    google.protobuf.Timestamp issue_date = 4;
    google.protobuf.Timestamp maturity_date = 5;
    int32 face_value = 6;
    State state = 7;

    // Additional unique field for entry
    string external_id = 8;
}

// CommercialPaperId identifier part
message CommercialPaperId {
    string issuer = 1;
    string paper_number = 2;
}

// ExternalId
message ExternalId {
    string id = 1;
}

// Container for returning multiple entities
message CommercialPaperList {
    repeated CommercialPaper items = 1;
}

// IssueCommercialPaper event
message IssueCommercialPaper {
    string issuer = 1 [(validator.field) = {string_not_empty : true}];
    string paper_number = 2 [(validator.field) = {string_not_empty : true}];
    google.protobuf.Timestamp issue_date = 3 [(validator.field) = {msg_exists : true}];
    google.protobuf.Timestamp maturity_date = 4 [(validator.field) = {msg_exists : true}];
    int32 face_value = 5 [(validator.field) = {int_gt : 0}];

    // external_id  - once more uniq id of state entry
    string external_id = 6 [(validator.field) = {string_not_empty : true}];
}

// BuyCommercialPaper event
message BuyCommercialPaper {
    string issuer = 1 [(validator.field) = {string_not_empty : true}];
    string paper_number = 2 [(validator.field) = {string_not_empty : true}];
    string current_owner = 3 [(validator.field) = {string_not_empty : true}];
    string new_owner = 4 [(validator.field) = {string_not_empty : true}];
    int32 price = 5 [(validator.field) = {int_gt : 0}];
    google.protobuf.Timestamp purchase_date = 6 [(validator.field) = {msg_exists : true}];
}

// RedeemCommercialPaper event
message RedeemCommercialPaper {
    string issuer = 1 [(validator.field) = {string_not_empty : true}];
    string paper_number = 2 [(validator.field) = {string_not_empty : true}];
    string redeeming_owner = 3 [(validator.field) = {string_not_empty : true}];
    google.protobuf.Timestamp redeem_date = 4 [(validator.field) = {msg_exists : true}];
}
```

#### Chaincode as service

Chaincode interface can be described with gRPC [service](../examples/cpaper_asservice/service/service.proto) notation.
Using `grpc-gateway` option we can define mapping for chaincode REST-API.

The `grpc-gateway` is a plugin of the Google protocol buffers compiler `protoc`. It reads protobuf service definitions and 
generates a reverse-proxy server which 'translates a RESTful HTTP API into gRPC. This server is generated according
 to the `google.api.http` annotations in your service definitions.

```proto
syntax = "proto3";

package service;

import "google/api/annotations.proto";
import "google/protobuf/empty.proto";
import "github.com/s7techlab/cckit/examples/cpaper_asservice/schema/schema.proto";

service CPaper {
    // List method returns all registered commercial papers
    rpc List (google.protobuf.Empty) returns (schema.CommercialPaperList) {
        option (google.api.http) = {
            get: "/cpaper"
        };
    }

    // Get method returns commercial paper data by id
    rpc Get (schema.CommercialPaperId) returns (schema.CommercialPaper) {
        option (google.api.http) = {
            get: "/cpaper/{issuer}/{paper_number}"
        };
    }

    // GetByExternalId
    rpc GetByExternalId  (schema.ExternalId) returns (schema.CommercialPaper) {
        option (google.api.http) = {
            get: "/cpaper/extid/{id}"
        };
    }

    // Issue commercial paper
    rpc Issue (schema.IssueCommercialPaper) returns (schema.CommercialPaper) {
        option (google.api.http) = {
            post : "/cpaper/issue"
        };
    }

    // Buy commercial paper
    rpc Buy (schema.BuyCommercialPaper) returns (schema.CommercialPaper) {
        option (google.api.http) = {
            post: "/cpaper/buy"
        };
    }

    // Redeem commercial paper
    rpc Redeem (schema.RedeemCommercialPaper) returns (schema.CommercialPaper) {
        option (google.api.http) = {
            post: "/cpaper/redeem"
        };
    }

    // Delete commercial paper
    rpc Delete (schema.CommercialPaperId) returns (schema.CommercialPaper) {
        option (google.api.http) = {
            delete: "/cpaper/{issuer}/{paper_number}"
        };
    }
}
```

####  Generator

[Makefile](../examples/cpaper_asservice/Makefile)

```makefile
.: generate

generate:
	@protoc --version
	@echo "commercial paper schema proto generation"
	@protoc -I=./schema/ \
	-I=../../vendor \
	--go_out=./schema/    \
	--govalidators_out=./schema/ \
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
	./service/service.proto
```

#### Chaincode implementation

Chaincode implementation must contain [state and event mappings](../examples/cpaper_asservice/service/service.pb.cc.go)

```go
type CPaperImpl struct {
}

func (cc *CPaperImpl) state(ctx router.Context) m.MappedState {
	return m.WrapState(ctx.State(), m.StateMappings{}.
		//  Create mapping for Commercial Paper entity
		Add(&schema.CommercialPaper{},
			m.PKeySchema(&schema.CommercialPaperId{}), // Key namespace will be <"CommercialPaper", Issuer, PaperNumber>
			m.List(&schema.CommercialPaperList{}),     // Structure of result for List method
			m.UniqKey("ExternalId"),                   // External Id is unique
		))
}

func (cc *CPaperImpl) event(ctx router.Context) state.Event {
	return m.WrapEvent(ctx.Event(), m.EventMappings{}.
		// Event name will be "IssueCommercialPaper", payload - same as issue payload
		Add(&schema.IssueCommercialPaper{}).
		// Event name will be "BuyCommercialPaper"
		Add(&schema.BuyCommercialPaper{}).
		// Event name will be "RedeemCommercialPaper"
		Add(&schema.RedeemCommercialPaper{}))
}
```

Chaincode service implementation must conform to generated from service definition
 [CPaperChaincode](../examples/cpaper_asservice/service/service.pb.cc.go) interface:

```go
func (cc *CPaperImpl) List(ctx router.Context, in *empty.Empty) (*schema.CommercialPaperList, error) {
	if res, err := cc.state(ctx).List(&schema.CommercialPaper{}); err != nil {
		return nil, err
	} else {
		return res.(*schema.CommercialPaperList), nil
	}
}

func (cc *CPaperImpl) Get(ctx router.Context, id *schema.CommercialPaperId) (*schema.CommercialPaper, error) {
	if res, err := cc.state(ctx).Get(id, &schema.CommercialPaper{}); err != nil {
		return nil, err
	} else {
		return res.(*schema.CommercialPaper), nil
	}
}

func (cc *CPaperImpl) GetByExternalId(ctx router.Context, id *schema.ExternalId) (*schema.CommercialPaper, error) {
	if res, err := cc.state(ctx).GetByUniqKey(
		&schema.CommercialPaper{}, "ExternalId", []string{id.Id}, &schema.CommercialPaper{}); err != nil {
		return nil, err
	} else {
		return res.(*schema.CommercialPaper), nil
	}
}

...

```

Then, chaincode service implementation can be embedded into chaincode method router with generated 
[RegisterCPaperChaincode](../examples/cpaper_asservice/service/service.pb.cc.go#L58) function:

```go 
func CCRouter(name string) (*router.Group, error) {
	r := router.New(name)
	// Store on the ledger the information about chaincode instantiation
	r.Init(owner.InvokeSetFromCreator)

	if err := service.RegisterCPaperChaincode(r, &CPaperImpl{}); err != nil {
		return nil, err
	}

	return r, nil
}
```