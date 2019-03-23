# Hyperledger Fabric chaincode kit (CCKit)

[![Go Report Card](https://goreportcard.com/badge/github.com/s7techlab/cckit)](https://goreportcard.com/report/github.com/s7techlab/cckit)
![Build](https://api.travis-ci.org/s7techlab/cckit.svg?branch=master)
[![Coverage Status](https://coveralls.io/repos/github/s7techlab/cckit/badge.svg?branch=master)](https://coveralls.io/github/s7techlab/cckit?branch=master)

**CCKit** is a **programming toolkit** for developing and testing hyperledger fabric golang chaincodes.  It enhances
 the development experience while providing developers components for creating more readable and secure 
 smart contracts (chaincodes).

## Overview

A [smart contract](https://hyperledger-fabric.readthedocs.io/en/latest/glossary.html#smart-contract) is code – 
invoked by a client application external to the blockchain network – that manages access and modifications to a set of
key-value pairs in the World State.  In Hyperledger Fabric, smart contracts are referred to as chaincode.

### CCKit features 

* [Centralized chaincode invocation handling](router) with methods routing and middleware capabilities 
* [Chaincode state modelling](state) using protocol buffers / json marshalling
* [MockStub testing](testing), allowing to immediately receive test results
* [Data encryption](extensions/encryption) on application level
* Chaincode method access control


### Problems with existing chaincode examples

There are several chaincode examples available: 

* [Commercial paper](https://hyperledger-fabric.readthedocs.io/en/release-1.4/developapps/smartcontract.html) from official [Hyperledger Fabric documentation](https://hyperledger-fabric.readthedocs.io)
* [Blockchain insurance application](https://github.com/IBM/build-blockchain-insurance-app) ( testing tutorial:  how to [write tests for "insurance" chaincode](examples/insurance) )
* [Marbles from hyperledger](https://github.com/hyperledger/fabric/blob/release-1.1/examples/chaincode/go/marbles02/marbles_chaincode.go)
* [Marbles from IBM-Blockchain](https://github.com/IBM-Blockchain/marbles/blob/master/chaincode/src/marbles/marbles.go)

Main problems: 

* Chaincode methods routing appeared only in HLF 1.4 and only in Node.Js chaincode
* Lots of code duplication (json marshalling / unmarshalling, validation, access control etc)
* Uncompleted testing tools (MockStub)

### Publications

* [Hyperledger Fabric smart contract data model: protobuf to chaincode state mapping](https://medium.com/coinmonks/hyperledger-fabric-smart-contract-data-model-protobuf-to-chaincode-state-mapping-191cdcfa0b78)
* [ERC20 token as Hyperledger Fabric Golang chaincode](https://medium.com/@viktornosov/erc20-token-as-hyperledger-fabric-golang-chaincode-d09dfd16a339)
* [CCKit: Routing and middleware for Hyperledger Fabric Golang chaincode](https://medium.com/@viktornosov/routing-and-middleware-for-developing-hyperledger-fabric-chaincode-written-in-go-90913951bf08)
* [Developing and testing Hyperledger Fabric smart contracts](https://habr.com/post/426705/) [RUS]


## Installation

CCKit requires Go 1.11+ with modules support

### Standalone
 
`git clone git@github.com:s7techlab/cckit.git`

`go mod vendor`

### As dependency

`go get -u github.com/s7techlab/cckit`


## Example based on CCKit
 
### Chaincode "Cars" 

Car registration chaincode use simple golang structure with json marshalling. Example with protobuf model is [here](state).

[source code](examples/cars/cars.go),  [tests](examples/cars/cars_test.go)

```go
// Simple CRUD chaincode for store information about cars
package main

import (
	"errors"
	"time"

	"github.com/s7techlab/cckit/extensions/owner"
	"github.com/s7techlab/cckit/router"
	p "github.com/s7techlab/cckit/router/param"
)

var (
	ErrCarAlreadyExists = errors.New(`car already exists`)
)

const CarEntity = `CAR`
const CarRegisteredEvent = `CAR_REGISTERED`

// CarPayload chaincode method argument
type CarPayload struct {
	Id    string
	Title string
	Owner string
}

// Car struct for chaincode state
type Car struct {
	Id    string
	Title string
	Owner string

	UpdatedAt time.Time // set by chaincode method
}

// Key for car entry in chaincode state
func (c Car) Key() ([]string, error) {
	return []string{CarEntity, c.Id}, nil
}

func New() *router.Chaincode {
	r := router.New(`cars`) // also initialized logger with "cars" prefix

	r.Init(invokeInit)

	r.Group(`car`).
		// everyone  can view information about registered cars.
		Query(`List`, queryCars).                                             // chain code method name is carList
		Query(`Get`, queryCar, p.String(`id`)).                               // chain code method name is carGet, method has 1 string argument "id"
		
		Invoke(`Register`, invokeCarRegister, p.Struct(`car`, &CarPayload{}), // 1 struct argument
			owner.Only) // allow access to method only for chaincode owner (authority)
			            // Only authority can register car information

	return router.NewChaincode(r)
}

// ======= Init ==================
func invokeInit(c router.Context) (interface{}, error) {
	return owner.SetFromCreator(c)
}

// ======= Chaincode methods =====

// car get info chaincode method handler
func queryCar(c router.Context) (interface{}, error) {
	// get state entry by composite key using CarKeyPrefix and car.Id
	//  and unmarshal from []byte to Car struct
	return c.State().Get(&Car{Id: c.ParamString(`id`)})
}

// cars car list chaincode method handler
func queryCars(c router.Context) (interface{}, error) {
	return c.State().List(
		CarEntity, // get list of state entries of type CarKeyPrefix
		&Car{})    // unmarshal from []byte and append to []Car slice
}

// carRegister car register chaincode method handler
func invokeCarRegister(c router.Context) (interface{}, error) {
	// arg name defined in router method definition
	p := c.Param(`car`).(CarPayload)

	t, _ := c.Time() // tx time
	car := &Car{     // data for chaincode state
		Id:        p.Id,
		Title:     p.Title,
		Owner:     p.Owner,
		UpdatedAt: t,
	}

	// trigger event
	c.Event().Set(CarRegisteredEvent, car)

	return car, // peer.Response payload will be json serialized car data
		//put json serialized data to state
		// create composite key using CarKeyPrefix and car.Id
		c.State().Insert(car)
}
```

### Test for chaincode

Tests are based on a modified [MockStub](testing/mockstub.go)

```go
package main

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	examplecert "github.com/s7techlab/cckit/examples/cert"
	"github.com/s7techlab/cckit/extensions/owner"
	testcc "github.com/s7techlab/cckit/testing"
	expectcc "github.com/s7techlab/cckit/testing/expect"
)

func TestCars(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cars Suite")
}

var _ = Describe(`Cars`, func() {

	//Create chaincode mock
	cc := testcc.NewMockStub(`cars`, New())

	// load actor certificates
	actors, err := examplecert.Actors(map[string]string{
		`authority`: `s7techlab.pem`,
		`someone`:   `victor-nosov.pem`,
	})
	if err != nil {
		panic(err)
	}

	// cars fixtures
	car1 := &Car{
		Id:    `A777MP77`,
		Title: `BMW`,
		Owner: `victor-nosov`,
	}

	car2 := &Car{
		Id:    `O888OO77`,
		Title: `TOYOTA`,
		Owner: `alexander`,
	}

	BeforeSuite(func() {
		// init chaincode
		expectcc.ResponseOk(cc.From(actors[`authority`]).Init()) // init chaincode from authority
	})

	Describe("Car", func() {

		It("Allow authority to add information about car", func() {
			//invoke chaincode method from authority actor
			expectcc.ResponseOk(cc.From(actors[`authority`]).Invoke(`carRegister`, car1))
		})

		It("Disallow non authority to add information about car", func() {
			//invoke chaincode method from non authority actor
			expectcc.ResponseError(
				cc.From(actors[`someone`]).Invoke(`carRegister`, car1),
				owner.ErrOwnerOnly) // expect "only owner" error
		})

		It("Disallow authority to add duplicate information about car", func() {
			expectcc.ResponseError(
				cc.From(actors[`authority`]).Invoke(`carRegister`, car1),
				ErrCarAlreadyExists) //expect already exists
		})

		It("Allow everyone to retrieve car information", func() {
			car := expectcc.PayloadIs(cc.Invoke(`carGet`, car1.Id),
				&Car{}).(Car)

			Expect(car.Title).To(Equal(car1.Title))
			Expect(car.Id).To(Equal(car1.Id))
		})

		It("Allow everyone to get car list", func() {
			//  &[]Car{} - declares target type for unmarshalling from []byte received from chaincode
			cars := expectcc.PayloadIs(cc.Invoke(`carList`), &[]Car{}).([]Car)

			Expect(len(cars)).To(Equal(1))
			Expect(cars[0].Id).To(Equal(car1.Id))
		})

		It("Allow authority to add more information about car", func() {
			// register second car
			expectcc.ResponseOk(cc.Invoke(`carRegister`, car2))
			cars := expectcc.PayloadIs(
				cc.From(actors[`authority`]).Invoke(`carList`),
				&[]Car{}).([]Car)

			Expect(len(cars)).To(Equal(2))
		})
	})
})
```

