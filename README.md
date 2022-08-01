# Hyperledger Fabric chaincode kit (CCKit)

[![Go Report Card](https://goreportcard.com/badge/github.com/s7techlab/cckit)](https://goreportcard.com/report/github.com/s7techlab/cckit)
![Build](https://api.travis-ci.org/s7techlab/cckit.svg?branch=master)
[![Coverage Status](https://coveralls.io/repos/github/s7techlab/cckit/badge.svg?branch=master)](https://coveralls.io/github/s7techlab/cckit?branch=master)

## Overview

A [smart contract](https://hyperledger-fabric.readthedocs.io/en/latest/glossary.html#smart-contract) is code, 
invoked by a client application external to the blockchain network – that manages access and modifications to a set of
key-value pairs in the World State.  In Hyperledger Fabric, smart contracts are referred to as chaincode.

**CCKit** is a **programming toolkit** for developing and testing Hyperledger Fabric golang chaincodes. It enhances
the development experience while providing developers components for creating more readable and secure
smart contracts.

## Chaincode examples

There are several chaincode "official" examples available: 

* [Commercial paper](https://hyperledger-fabric.readthedocs.io/en/release-1.4/developapps/smartcontract.html) from official [Hyperledger Fabric documentation](https://hyperledger-fabric.readthedocs.io)
* [Blockchain insurance application](https://github.com/IBM/build-blockchain-insurance-app) (testing tutorial: how to [write tests for "insurance" chaincode](examples/insurance))

and [others](docs/chaincode-examples.md)

**Main problems** with existing examples are: 

* Working with chaincode state at very low level
* Lots of code duplication (JSON marshalling / unmarshalling, validation, access control, etc)
* Chaincode methods routing appeared only in HLF 1.4 and only in Node.Js chaincode
* Uncompleted testing tools (MockStub)

### CCKit features 

* [Chaincode method router](router) with invocation handlers and middleware capabilities 
* [Chaincode state modeling](state) using [protocol buffers](examples/cpaper_extended) / [golang struct to json marshalling](examples/cars), with [private data support](examples/private_cars)
* Designing chaincode in [gRPC service notation](gateway) with code generation of chaincode SDK, gRPC and REST-API
* [MockStub testing](testing), allowing to immediately receive test results
* [Data encryption](extensions/encryption) on application level
* Chaincode method [access control](extensions/owner)

### Publications with usage examples 

* [Service-oriented Hyperledger Fabric application development using gRPC definitions](https://medium.com/coinmonks/service-oriented-hyperledger-fabric-application-development-32e66f578f9a)
* [Hyperledger Fabric smart contract data model: protobuf to chaincode state mapping](https://medium.com/coinmonks/hyperledger-fabric-smart-contract-data-model-protobuf-to-chaincode-state-mapping-191cdcfa0b78)
* [Hyperledger Fabric chaincode test driven development (TDD) with unit testing](https://medium.com/coinmonks/test-driven-hyperledger-fabric-golang-chaincode-development-dbec4cb78049)
* [ERC20 token as Hyperledger Fabric Golang chaincode](https://medium.com/@viktornosov/erc20-token-as-hyperledger-fabric-golang-chaincode-d09dfd16a339)
* [CCKit: Routing and middleware for Hyperledger Fabric Golang chaincode](https://medium.com/@viktornosov/routing-and-middleware-for-developing-hyperledger-fabric-chaincode-written-in-go-90913951bf08)
* [Developing and testing Hyperledger Fabric smart contracts](https://habr.com/post/426705/) [RUS]

## Examples based on CCKit

* [Cars](examples/cars) - car registration chaincode, *simplest* example
* [Commercial paper, service-oriented approach](https://github.com/s7techlab/hyperledger-fabric-samples) - 
  recommended way to start new application. Code generation radically simplifies building on-chain and off-chain applications.


* [Commercial paper](examples/cpaper) - faithful reimplementation of the official example 
* [Commercial paper extended example](examples/cpaper_extended) - with protobuf chaincode state schema and other features
* [ERC-20](examples/erc20) - tokens smart contract, implementing ERC-20 interface
* [Cars private](examples/private_cars) - car registration chaincode with private data
* [Payment](examples/payment) - a few examples of chaincodes with encrypted state 
 
## Installation

CCKit requires Go 1.11+ with modules support

### Standalone
 
`git clone git@github.com:s7techlab/cckit.git`

`go mod vendor`

### As dependency

`go get -u github.com/s7techlab/cckit`

## Example - Commercial Paper chaincode

### Scenario

[Commercial paper](https://hyperledger-fabric.readthedocs.io/en/release-1.4/developapps/scenario.html) 
scenario from official documentation describes a Hyperledger Fabric network, aimed to issue, buy and redeem 
commercial paper.

![commercial paper network](examples/cpaper/img/cpaper-network.png)

### 5 steps to develop chaincode

Chaincode is a domain specific program which relates to specific business process. The job of a smart
contract developer is to take an existing business process and express it as a smart contract in a
programming language. Steps of chaincode development:

1. Define chaincode model - schema for state entries, transaction payload and events
2. Define chaincode interface
3. Implement chaincode instantiate method
4. Implement chaincode methods with business logic
5. Create tests


### Define chaincode model

With protocol buffers, you write a `.proto` description of the data structure you wish to store.
From that, the protocol buffer compiler creates a golang struct that implements automatic encoding
and parsing of the protocol buffer data with an efficient binary format (or json).

Code generation can be simplified with a short [Makefile](examples/cpaper_extended/schema):

```makefile
.: generate

generate:
	@echo "schema"
	@protoc -I=./ --go_out=./ ./*.proto
```

#### Chaincode state

The following file shows how to define the world state schema using protobuf.

[examples/cpaper_extended/schema/state.proto](examples/cpaper_extended/schema/state.proto)

```proto
syntax = "proto3";

package cckit.examples.cpaper_extended.schema;
option go_package = "schema";

import "google/protobuf/timestamp.proto";

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

// Container for returning multiple entities
message CommercialPaperList {
    repeated CommercialPaper items = 1;
}
```

#### Chaincode transaction and events payload

This file defines the data payload used in the business logic methods.
In this example transaction and event payloads are exactly the same for the sake of brevity, but you could
create a different schema for each type of payload.

[examples/cpaper_extended/schema/payload.proto](examples/cpaper_extended/schema/payload.proto)

```proto
// IssueCommercialPaper event
syntax = "proto3";

package cckit.examples.cpaper_extended.schema;
option go_package = "schema";

import "google/protobuf/timestamp.proto";
import "github.com/mwitkow/go-proto-validators/validator.proto";

// IssueCommercialPaper event
message IssueCommercialPaper {
    string issuer = 1;
    string paper_number = 2;
    google.protobuf.Timestamp issue_date = 3;
    google.protobuf.Timestamp maturity_date = 4;
    int32 face_value = 5;

    // external_id - another unique constraint
    string external_id = 6;
}

// BuyCommercialPaper event
message BuyCommercialPaper {
    string issuer = 1;
    string paper_number = 2;
    string current_owner = 3;
    string new_owner = 4;
    int32 price = 5;
    google.protobuf.Timestamp purchase_date = 6;
}

// RedeemCommercialPaper event
message RedeemCommercialPaper {
    string issuer = 1;
    string paper_number = 2;
    string redeeming_owner = 3;
    google.protobuf.Timestamp redeem_date = 4;
}
```

### Define chaincode interface

In [examples/cpaper_extended/chaincode.go](examples/cpaper_extended/chaincode.go) file we will define the mappings,
chaincode initialization method and business logic in the transaction methods.
For brevity, we will only display snippets of the code here, please refer to the original file for full example.

Firstly we define mapping rules. These specify the struct used to hold a specific chaincode state, it's primary key, list mapping, unique keys, etc.
Then we define the schemas used for emitting events.

```go
var (
	// State mappings
	StateMappings = m.StateMappings{}.
		// Create mapping for Commercial Paper entity
		Add(&schema.CommercialPaper{},
			// Key namespace will be <"CommercialPaper", Issuer, PaperNumber>
			m.PKeySchema(&schema.CommercialPaperId{}),
			// Structure of result for List method
			m.List(&schema.CommercialPaperList{}),
			// External Id is unique
			m.UniqKey("ExternalId"),
		)

	// EventMappings
	EventMappings = m.EventMappings{}.
		// Event name will be "IssueCommercialPaper", payload - same as issue payload
		Add(&schema.IssueCommercialPaper{}).
		// Event name will be "BuyCommercialPaper"
		Add(&schema.BuyCommercialPaper{}).
		// Event name will be "RedeemCommercialPaper"
		Add(&schema.RedeemCommercialPaper{})
)
```

CCKit uses [router](router) to define rules about how to map chaincode invocation to a particular handler,
as well as what kind of middleware needs to be used during a request, for example how to convert incoming argument from
[]byte to target type (string, struct, etc).

```go
func NewCC() *router.Chaincode {

	r := router.New(`commercial_paper`)

	// Mappings for chaincode state
	r.Use(m.MapStates(StateMappings))

	// Mappings for chaincode events
	r.Use(m.MapEvents(EventMappings))

	// Store in chaincode state information about chaincode first instantiator
	r.Init(owner.InvokeSetFromCreator)

	// Method for debug chaincode state
	debug.AddHandlers(r, `debug`, owner.Only)

	r.
		// read methods
		Query(`list`, cpaperList).

		Query(`get`, cpaperGet, defparam.Proto(&schema.CommercialPaperId{})).

		// txn methods
		Invoke(`issue`, cpaperIssue, defparam.Proto(&schema.IssueCommercialPaper{})).
		Invoke(`buy`, cpaperBuy, defparam.Proto(&schema.BuyCommercialPaper{})).
		Invoke(`redeem`, cpaperRedeem, defparam.Proto(&schema.RedeemCommercialPaper{})).
		Invoke(`delete`, cpaperDelete, defparam.Proto(&schema.CommercialPaperId{}))

	return router.NewChaincode(r)
}
```

### Implement chaincode `init` method

In many cases during chaincode instantiation we need to define permissions for chaincode functions -
"who is allowed to do this thing", incredibly important in the world of smart contracts.
The most common and basic form of access control is the concept of `ownership`: there's one account (combination
of MSP and certificate identifiers) that is the owner and can do administrative tasks on contracts. This 
approach is perfectly reasonable for contracts that only have a single administrative user.

CCKit provides `owner` extension for implementing ownership and access control in Hyperledger Fabric chaincodes.
In the previous snippet, as an `init` method, we used [owner.InvokeSetFromCreator](extensions/owner/handler.go), storing information
which stores the information about who is the owner into the world state upon chaincode instantiation.

### Implement business rules as chaincode methods

Now we have to define the actual business logic which will modify the world state when a transaction occurs.
In this example we will show only the `buy` method for brevity.
Please refer to [examples/cpaper_extended/chaincode.go](examples/cpaper_extended/chaincode.go) for full implementation.

```go
func invokeCPaperBuy(c router.Context) (interface{}, error) {
	var (
		cpaper *schema.CommercialPaper

		// Buy transaction payload
		buyData = c.Param().(*schema.BuyCommercialPaper)

		// Get the current commercial paper state
		cp, err = c.State().Get(
			&schema.CommercialPaperId{Issuer: buyData.Issuer, PaperNumber: buyData.PaperNumber},
			&schema.CommercialPaper{})
	)

	if err != nil {
		return nil, errors.Wrap(err, "not found")
	}

	cpaper = cp.(*schema.CommercialPaper)

	// Validate current owner
	if cpaper.Owner != buyData.CurrentOwner {
		return nil, fmt.Errorf(
			"paper %s %s is not owned by %s",
			cpaper.Issuer, cpaper.PaperNumber, buyData.CurrentOwner)
	}

	// First buyData moves state from ISSUED to TRADING
	if cpaper.State == schema.CommercialPaper_ISSUED {
		cpaper.State = schema.CommercialPaper_TRADING
	}

	// Check paper is not already REDEEMED
	if cpaper.State == schema.CommercialPaper_TRADING {
		cpaper.Owner = buyData.NewOwner
	} else {
		return nil, fmt.Errorf(
			"paper %s %s is not trading.current state = %s",
			cpaper.Issuer, cpaper.PaperNumber, cpaper.State)
	}

	if err = c.Event().Set(buyData); err != nil {
		return nil, err
	}

	return cpaper, c.State().Put(cpaper)
}
```

### Test chaincode functionality

And finally we should write tests to ensure our business logic is behaving as it should.
Again, for brevity, we omitted most of the code from [examples/cpaper_extended/chaincode_test.go](examples/cpaper_extended/chaincode_test.go).
CCKit support chaincode testing with [MockStub](testing).

```go
var _ = Describe(`CommercialPaper`, func() {
	paperChaincode := testcc.NewMockStub(`commercial_paper`, NewCC())

	BeforeSuite(func() {
		// Init chaincode with admin identity
		expectcc.ResponseOk(
			paperChaincode.
				From(testdata.GetTestIdentity(MspName, path.Join("testdata", "admin", "admin.pem"))).
				Init())
	})

	Describe("Commercial Paper lifecycle", func() {
		// ...

		It("Allow buyer to buy commercial paper", func() {
			buyTransactionData := &schema.BuyCommercialPaper{
				Issuer:       IssuerName,
				PaperNumber:  "0001",
				CurrentOwner: IssuerName,
				NewOwner:     BuyerName,
				Price:        95000,
				PurchaseDate: ptypes.TimestampNow(),
			}

			expectcc.ResponseOk(paperChaincode.Invoke(`buy`, buyTransactionData))

			queryResponse := paperChaincode.Query("get", &schema.CommercialPaperId{
				Issuer:      IssuerName,
				PaperNumber: "0001",
			})

			paper := expectcc.PayloadIs(queryResponse, &schema.CommercialPaper{}).(*schema.CommercialPaper)

			Expect(paper.Owner).To(Equal(BuyerName))
			Expect(paper.State).To(Equal(schema.CommercialPaper_TRADING))

			Expect(<-paperChaincode.ChaincodeEventsChannel).To(BeEquivalentTo(&peer.ChaincodeEvent{
				EventName: `BuyCommercialPaper`,
				Payload:   testcc.MustProtoMarshal(buyTransactionData),
			}))

			paperChaincode.ClearEvents()
		})

		// ...

	})
})
```
