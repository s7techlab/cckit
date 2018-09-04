// Simple CRUD chaincode for store information about cars
package cars

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/s7techlab/cckit/router"
	p "github.com/s7techlab/cckit/router/param"
)

type ChaincodeWithoutAccessControl struct {
	router *router.Group
}

func NewWithoutAccessControl() *ChaincodeWithoutAccessControl {
	r := router.New(`cars_without_access_control`) // also initialized logger with "cars" prefix

	r.Group(`car`).
		Query(`List`, cars).               // chain code method name is carList
		Query(`Get`, car, p.String(`id`)). // chain code method name is carGet, method has 1 string argument "id"
		Invoke(`Register`, carRegister, p.Struct(`car`, &CarPayload{}))

	return &ChaincodeWithoutAccessControl{r}
}

//========  Base methods ====================================
//
// Init initializes chain code - sets chaincode "owner"
func (cc *ChaincodeWithoutAccessControl) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

// Invoke - entry point for chain code invocations
func (cc *ChaincodeWithoutAccessControl) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// delegate handling to router
	return cc.router.Handle(stub)
}
