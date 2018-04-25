# Hyperledger Fabric chaincode kit (CCKit)

[![Go Report Card](https://goreportcard.com/badge/github.com/s7techlab/cckit)](https://goreportcard.com/report/github.com/s7techlab/cckit)


**CCkit** is a **programming toolkit** for developing and testing hyperledger fabric chaincode

## Overview

### Features 

* Centralized chaincode invocation handling
* Middleware support
* Chaincode method access control
* Automatic json marshalling / unmarshalling
* MockStub testing

### Existing chaincode examples

* https://github.com/hyperledger/fabric/blob/release-1.1/examples/chaincode/go/marbles02/marbles_chaincode.go
* https://github.com/IBM-Blockchain/marbles/blob/master/chaincode/src/marbles/marbles.go
* https://github.com/IBM-Blockchain-Archive/car-lease-demo/blob/master/Chaincode/src/vehicle_code/vehicles.go

## Example based on CCKit

### Chaincode "Cars" 

[source code](examples/cars/cars.go),  [tests](examples/cars/cars_test.go)

```go
// Simple CRUD chaincode for store information about cars
package main

import (
	"errors"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/s7techlab/cckit/extensions/owner"
	"github.com/s7techlab/cckit/router"
	p "github.com/s7techlab/cckit/router/param"
)

var (
	ErrCarAlreadyExists = errors.New(`car already exists`)
)

const CarKeyPrefix = `CAR`

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

type Chaincode struct {
	router *router.Group
}

func New() *Chaincode {
	r := router.New(`cars`) // also initialized logger with "cars" prefix

	r.Group(`car`).
		Query(`List`, cars).                                                        // chain code method name is carList
		Query(`Get`, car, p.String(`id`)).                                          // chain code method name is carGet
		Invoke(`Register`, carRegister, p.Struct(`car`, &CarPayload{}), owner.Only) // only owner (authority)

	return &Chaincode{r}
}

//========  Base methods ====================================
//
// Init initializes chain code - sets chaincode "owner"
func (cc *Chaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	// set owner of chain code with special permissions , based on tx creator certificate
	return owner.SetFromCreator(cc.router.Context(`init`, stub))
}

// Invoke - entry point for chain code invocations
func (cc *Chaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// delegate handling to router
	return cc.router.Handle(stub)
}

// Key for car entry in chaincode state
func Key(id string) []string {
	return []string{CarKeyPrefix, id}
}

// ======= Chaincode methods

// car get info chaincode method handler
func car(c router.Context) peer.Response {
	return c.Response().Create(
		c.State().Get( // get state entry
			Key(c.ArgString(`id`)), // by composite key using CarKeyPrefix and car.Id
			&Car{}))                // unmarshal from []byte to Car struct
}

// cars car list chaincode method handler
func cars(c router.Context) peer.Response {
	return c.Response().Create(
		c.State().List(
			CarKeyPrefix, // get list of state entries of type CarKeyPrefix
			&Car{}))      // unmarshal from []byte and append to []Car slice
}

// carRegister car register chaincode method handler
func carRegister(c router.Context) peer.Response {
	// arg name defined in router method definition
	p := c.Arg(`car`).(CarPayload)
	if exists, err := c.State().Exists(Key(p.Id)); exists || err != nil {
		return c.Response().Error(ErrCarAlreadyExists)
	}

	t, _ := c.Time() // tx time
	car := &Car{
		Id:        p.Id,
		Title:     p.Title,
		Owner:     p.Owner,
		UpdatedAt: t} // data for chaincode state

	return c.Response().Create(
		car, // peer.Response payload will json serialized car data
		c.State().Put( //put json serialized data to state
			Key(car.Id), // create composite key using CarKeyPrefix and car.Id
			car))
}
```

### Test for chaincode

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