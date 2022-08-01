// Package cars simple CRUD chaincode for store information about cars
package cars

import (
	"github.com/s7techlab/cckit/router"
	p "github.com/s7techlab/cckit/router/param"
)

func NewWithoutAccessControl() *router.Chaincode {
	r := router.New(`cars_without_access_control`) // also initialized logger with "cars" prefix

	r.Group(`car`).
		Query(`List`, queryCars).                                             // chain code method name is carList
		Query(`Get`, queryCar, p.String(`id`)).                               // chain code method name is carGet, method has 1 string argument "id"
		Invoke(`Register`, invokeCarRegister, p.Struct(`car`, &CarPayload{})) // allow access to everyone
	return router.NewChaincode(r)
}
