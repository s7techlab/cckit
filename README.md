# Hyperledger Fabric chaincode kit (CCKit)

[![Go Report Card](https://goreportcard.com/badge/github.com/s7techlab/cckit)](https://goreportcard.com/report/github.com/s7techlab/cckit)
![Build](https://api.travis-ci.org/s7techlab/cckit.svg?branch=master)
[![Coverage Status](https://coveralls.io/repos/github/s7techlab/cckit/badge.svg?branch=master)](https://coveralls.io/github/s7techlab/cckit?branch=master)

## Overview

A [smart contract](https://hyperledger-fabric.readthedocs.io/en/latest/glossary.html#smart-contract) is code, 
invoked by a client application external to the blockchain network – that manages access and modifications to a set of
key-value pairs in the World State.  In Hyperledger Fabric, smart contracts are referred to as chaincode.

**CCKit** is a **programming toolkit** for developing and testing hyperledger fabric golang chaincodes.  It enhances
 the development experience while providing developers components for creating more readable and secure 
 smart contracts.

## Chaincode examples

There are several chaincode "official" examples available: 

* [Commercial paper](https://hyperledger-fabric.readthedocs.io/en/release-1.4/developapps/smartcontract.html) from official [Hyperledger Fabric documentation](https://hyperledger-fabric.readthedocs.io)
* [Blockchain insurance application](https://github.com/IBM/build-blockchain-insurance-app) ( testing tutorial:  how to [write tests for "insurance" chaincode](examples/insurance) )

and [others](docs/chaincode-examples.md)

**Main problems** with existing examples are: 

* Working with chaincode state at very low level
* Lots of code duplication (json marshalling / unmarshalling, validation, access control etc)
* Chaincode methods routing appeared only in HLF 1.4 and only in Node.Js chaincode
* Uncompleted testing tools (MockStub)

### CCKit features 

* [Centralized chaincode invocation handling](router) with methods routing and middleware capabilities 
* [Chaincode state modelling](state) using [protocol buffers](examples/cpaper) / [golang struct to json marshalling](examples/cars), with private data support
* [MockStub testing](testing), allowing to immediately receive test results
* [Data encryption](extensions/encryption) on application level
* Chaincode method [access control](extensions/owner)

### Publications with usage examples 

* [Hyperledger Fabric smart contract data model: protobuf to chaincode state mapping](https://medium.com/coinmonks/hyperledger-fabric-smart-contract-data-model-protobuf-to-chaincode-state-mapping-191cdcfa0b78)
* [ERC20 token as Hyperledger Fabric Golang chaincode](https://medium.com/@viktornosov/erc20-token-as-hyperledger-fabric-golang-chaincode-d09dfd16a339)
* [CCKit: Routing and middleware for Hyperledger Fabric Golang chaincode](https://medium.com/@viktornosov/routing-and-middleware-for-developing-hyperledger-fabric-chaincode-written-in-go-90913951bf08)
* [Developing and testing Hyperledger Fabric smart contracts](https://habr.com/post/426705/) [RUS]

## Examples based on CCKit

* [Cars](examples/cars) - car registration chaincode, *simplest* example

* [Commercial paper](examples/cpaper) - faithful reimplementation of the official example 
* [Commercial paper extended example](examples/cpaper_extended) with protobuf chaincode state schema and other features
* [ERC-20](examples/erc20) - tokens smart contract, implementing ERC-20 interface
* [Cars private](examples/private_cars) -  car registration chaincode with private data
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
contract developer is to take an existing business process  and express it as a smart contract in a 
programming language.  Steps of chaincode development:

1. Define chaincode model - schema for state entries, input payload and events
2. Define chaincode interface
3. Implement chaincode instantiate method
4. Implement chaincode methods with business logic
5. Create tests


### Define chaincode model

With protocol buffers, you write a `.proto` description of the data structure you wish to store. 
From that, the protocol buffer compiler creates a golang struct that implements automatic encoding 
and parsing of the protocol buffer data with an efficient binary format (or json). 

Code generation can be simplified with short [Makefile](examples/cpaper/schema):

```makefile
.: generate

generate:
	@echo "schema"
	@protoc -I=./ --go_out=./ ./*.proto
```

#### Chaincode state

[state.proto](examples/cpaper/schema/state.proto)

```proto
syntax = "proto3";
package schema;

import "google/protobuf/timestamp.proto";

// CommercialPaper state entry
message CommercialPaper {

    enum State {
        ISSUED = 0;
        TRADING = 1;
        REDEEMED = 2;
    }

    // issuer and paper number comprises primary key of commercial paper entry
    string issuer = 1;
    string paper_number = 2;

    string owner = 3;
    google.protobuf.Timestamp issue_date = 4;
    google.protobuf.Timestamp maturity_date = 5;
    int32 face_value = 6;
    State state = 7;
}

// CommercialPaperId identifier part
message CommercialPaperId {
    string issuer = 1;
    string paper_number = 2;
}

message CommercialPaperList {
    repeated CommercialPaper items = 1;
}
```

#### Chaincode input payload and events

[payload.proto](examples/cpaper/schema/payload.proto)

```proto
// IssueCommercialPaper event
message IssueCommercialPaper {
    string issuer = 1;
    string paper_number = 2;
    google.protobuf.Timestamp issue_date = 3;
    google.protobuf.Timestamp maturity_date = 4;
    int32 face_value = 5;
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

CCKit uses [router](router) to define rules about how to map chaincode invocation to particular handler, 
as well as what kind of middleware needs to be used during request, for example how to convert incoming argument from []byte 
to target type (string, struct etc).

Also we can define mapping rules for creating chaincode state entries keys for protobuf structures.


```go

// State mappings
StateMappings = m.StateMappings{}.Add(
    &schema.CommercialPaper{}, // define mapping for this structure
    m.PKeySchema(&schema.CommercialPaperId{}),  // key  will be <`CommercialPaper`, Issuer, PaperNumber>
    m.List(&schema.CommercialPaperList{})) // structure-result of list method

// EventMappings
EventMappings = m.EventMappings{}.
    Add(&schema.IssueCommercialPaper{}).// event name will be `IssueCommercialPaper`,  payload - same as issue payload
    Add(&schema.BuyCommercialPaper{}).
    Add(&schema.RedeemCommercialPaper{})

func NewCC() *router.Chaincode {

	r := router.New(`commercial_paper`)

	// Mappings for chaincode state
	r.Use(m.MapStates(StateMappings))

	// Mappings for chaincode events
	r.Use(m.MapEvents(EventMappings))

	// store in chaincode state information about chaincode first instantiator
	r.Init(owner.InvokeSetFromCreator)

	// method for debug chaincode state
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

In many cases during chaincode instantiating we need to define permissions for chaincode functions -
"who is allowed to do this thing", incredibly important in the world of smart contracts.
The most common and basic form of access control is the concept of `ownership`: there's one account (combination
of MSP  and certificate identifiers) that is the owner and can do administrative tasks on contracts. This 
approach is perfectly reasonable for contracts that only have a single administrative user.

CCKit provides `owner` extension for implementing ownership and access control in Hyperledger Fabric chaincodes.
In this example we use as a `init` method [owner.InvokeSetFromCreator](extensions/owner/handler.go), storing information about owner in the
chaincode state.

### Implement business rules as chaincode methods

```go
package cpaper

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/s7techlab/cckit/examples/cpaper/schema"
	"github.com/s7techlab/cckit/router"
)

func cpaperList(c router.Context) (interface{}, error) {
	// commercial paper key is composite key <`CommercialPaper`>, {Issuer}, {PaperNumber} >
	// where `CommercialPaper` - namespace of this type
	// list method retrieves entries from chaincode state
	// using GetStateByPartialCompositeKey method, then unmarshal received from state bytes via proto.Ummarshal method
	// and creates slice of *schema.CommercialPaper
	return c.State().List(&schema.CommercialPaper{})
}

func cpaperIssue(c router.Context) (interface{}, error) {
	var (
		issue  = c.Param().(*schema.IssueCommercialPaper) //default parameter
		cpaper = &schema.CommercialPaper{
			Issuer:       issue.Issuer,
			PaperNumber:  issue.PaperNumber,
			Owner:        issue.Issuer,
			IssueDate:    issue.IssueDate,
			MaturityDate: issue.MaturityDate,
			FaceValue:    issue.FaceValue,
			State:        schema.CommercialPaper_ISSUED, // initial state
		}
		err error
	)

	if err = c.Event().Set(issue); err != nil {
		return nil, err
	}

	return cpaper, c.State().Insert(cpaper)
}

func cpaperBuy(c router.Context) (interface{}, error) {

	var (
		cpaper *schema.CommercialPaper

		// but tx payload
		buy = c.Param().(*schema.BuyCommercialPaper)

		// current commercial paper state
		cp, err = c.State().Get(
			&schema.CommercialPaperId{Issuer: buy.Issuer, PaperNumber: buy.PaperNumber},
			&schema.CommercialPaper{})
	)

	if err != nil {
		return nil, errors.Wrap(err, `not found`)
	}
	cpaper = cp.(*schema.CommercialPaper)

	// Validate current owner
	if cpaper.Owner != buy.CurrentOwner {
		return nil, fmt.Errorf(`paper %s %s is not owned by %s`, cpaper.Issuer, cpaper.PaperNumber, buy.CurrentOwner)
	}

	// First buy moves state from ISSUED to TRADING
	if cpaper.State == schema.CommercialPaper_ISSUED {
		cpaper.State = schema.CommercialPaper_TRADING
	}

	// Check paper is not already REDEEMED
	if cpaper.State == schema.CommercialPaper_TRADING {
		cpaper.Owner = buy.NewOwner
	} else {
		return nil, fmt.Errorf(`paper %s %s is not trading.current state = %s`, cpaper.Issuer, cpaper.PaperNumber, cpaper.State)
	}

	if err = c.Event().Set(buy); err != nil {
		return nil, err
	}

	return cpaper, c.State().Put(cpaper)
}

func cpaperRedeem(c router.Context) (interface{}, error) {
	// implement me
	return nil, nil
}

func cpaperGet(c router.Context) (interface{}, error) {
	return c.State().Get(c.Param().(*schema.CommercialPaperId))
}

func cpaperDelete(c router.Context) (interface{}, error) {
	return nil, c.State().Delete(c.Param().(*schema.CommercialPaperId))
}
```

### Test chaincode functionality

CCKit support chaincode testing with [Mockstub](testing)

```go
package mapping_test

import (
	"testing"

	"github.com/hyperledger/fabric/protos/peer"

	"github.com/golang/protobuf/ptypes"

	"github.com/golang/protobuf/proto"
	"github.com/s7techlab/cckit/examples/cpaper/schema"
	"github.com/s7techlab/cckit/examples/cpaper/testdata"
	"github.com/s7techlab/cckit/state"

	"github.com/s7techlab/cckit/examples/cpaper"

	examplecert "github.com/s7techlab/cckit/examples/cert"
	"github.com/s7techlab/cckit/identity"
	testcc "github.com/s7techlab/cckit/testing"
	expectcc "github.com/s7techlab/cckit/testing/expect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestState(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "State suite")
}

var (
	actors   identity.Actors
	cPaperCC *testcc.MockStub
	err      error
)
var _ = Describe(`Mapping`, func() {

	BeforeSuite(func() {
		actors, err = identity.ActorsFromPemFile(`SOME_MSP`, map[string]string{
			`owner`: `s7techlab.pem`,
		}, examplecert.Content)

		Expect(err).To(BeNil())

		//Create commercial papers chaincode mock - protobuf based schema
		cPaperCC = testcc.NewMockStub(`cpapers`, cpaper.NewCC())
		cPaperCC.From(actors[`owner`]).Init()

	})

	Describe(`Protobuf based schema`, func() {
		It("Allow to add data to chaincode state", func(done Done) {

			events := cPaperCC.EventSubscription()
			expectcc.ResponseOk(cPaperCC.Invoke(`issue`, &testdata.CPapers[0]))

			Expect(<-events).To(BeEquivalentTo(&peer.ChaincodeEvent{
				EventName: `IssueCommercialPaper`,
				Payload:   testcc.MustProtoMarshal(&testdata.CPapers[0]),
			}))

			expectcc.ResponseOk(cPaperCC.Invoke(`issue`, &testdata.CPapers[1]))
			expectcc.ResponseOk(cPaperCC.Invoke(`issue`, &testdata.CPapers[2]))

			close(done)
		}, 0.2)

		It("Disallow to insert entries with same keys", func() {
			expectcc.ResponseError(cPaperCC.Invoke(`issue`, &testdata.CPapers[0]))
		})

		It("Allow to get entry list", func() {
			cpapers := expectcc.PayloadIs(cPaperCC.Query(`list`), &[]schema.CommercialPaper{}).([]schema.CommercialPaper)
			Expect(len(cpapers)).To(Equal(3))
			Expect(cpapers[0].Issuer).To(Equal(testdata.CPapers[0].Issuer))
			Expect(cpapers[0].PaperNumber).To(Equal(testdata.CPapers[0].PaperNumber))
		})

		It("Allow to get entry raw protobuf", func() {
			cp := testdata.CPapers[0]
			cpaperProtoFromCC := cPaperCC.Query(`get`, &schema.CommercialPaperId{Issuer: cp.Issuer, PaperNumber: cp.PaperNumber}).Payload

			stateCpaper := &schema.CommercialPaper{
				Issuer:       cp.Issuer,
				PaperNumber:  cp.PaperNumber,
				Owner:        cp.Issuer,
				IssueDate:    cp.IssueDate,
				MaturityDate: cp.MaturityDate,
				FaceValue:    cp.FaceValue,
				State:        schema.CommercialPaper_ISSUED, // initial state
			}
			cPaperProto, _ := proto.Marshal(stateCpaper)
			Expect(cpaperProtoFromCC).To(Equal(cPaperProto))
		})

		It("Allow update data in chaincode state", func() {
			cp := testdata.CPapers[0]
			expectcc.ResponseOk(cPaperCC.Invoke(`buy`, &schema.BuyCommercialPaper{
				Issuer:       cp.Issuer,
				PaperNumber:  cp.PaperNumber,
				CurrentOwner: cp.Issuer,
				NewOwner:     `some-new-owner`,
				Price:        cp.FaceValue - 10,
				PurchaseDate: ptypes.TimestampNow(),
			}))

			cpaperFromCC := expectcc.PayloadIs(
				cPaperCC.Query(`get`, &schema.CommercialPaperId{Issuer: cp.Issuer, PaperNumber: cp.PaperNumber}),
				&schema.CommercialPaper{}).(*schema.CommercialPaper)

			// state is updated
			Expect(cpaperFromCC.State).To(Equal(schema.CommercialPaper_TRADING))
			Expect(cpaperFromCC.Owner).To(Equal(`some-new-owner`))
		})

		It("Allow to delete entry", func() {

			cp := testdata.CPapers[0]
			toDelete := &schema.CommercialPaperId{Issuer: cp.Issuer, PaperNumber: cp.PaperNumber}

			expectcc.ResponseOk(cPaperCC.Invoke(`delete`, toDelete))
			cpapers := expectcc.PayloadIs(cPaperCC.Invoke(`list`), &[]schema.CommercialPaper{}).([]schema.CommercialPaper)

			Expect(len(cpapers)).To(Equal(2))
			expectcc.ResponseError(cPaperCC.Invoke(`get`, toDelete), state.ErrKeyNotFound)
		})
	})

})
```