// Simple CRUD chaincode for store information about cars
package main

import (
	"errors"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/s7techlab/cckit/response"
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
		Query(`List`, cars).               // chain code method name is carList
		Query(`Get`, car, p.String(`id`)). // chain code method name is carGet
		Invoke(`Register`, carRegister, p.Struct(`car`, &CarPayload{}))

	return &Chaincode{r}
}

//========  Base methods ====================================
//
// Init initializes chain code
func (cc *Chaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return response.Success(nil)
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
func car(c router.Context) peer.Response {
	return c.Response().Create(
		c.State().Get(
			Key(c.ArgString(`id`)), // get state entry by composite key  using CarKeyPrefix and car.Id
			&Car{}))                // unmarshal from []byte to Car struct
}

func cars(c router.Context) peer.Response {
	return c.Response().Create(
		c.State().List(
			CarKeyPrefix, // get list of state entries of type CarKeyPrefix
			&Car{}))      // unmarshal from []byte and append to []Car slice
}

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
