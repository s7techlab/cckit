// Marbles chaincode example
package marbles

import (
	"github.com/s7techlab/cckit/examples/marbles/schema"
	"github.com/s7techlab/cckit/extensions/owner"
	"github.com/s7techlab/cckit/router"
	p "github.com/s7techlab/cckit/router/param"
	"github.com/s7techlab/cckit/state/middleware"
)

// New chaincode marbles
func New() *router.Chaincode {
	r := router.New(`marbles`).Use(middleware.Mapper(Mappings)) // use cc state mapping

	r.Init(owner.InvokeSetFromCreator)
	r.Group(`marble`).
		Query(`Init`, invokeMarbleInit, p.Proto(`marble`, &schema.Marble{})). // chain code method name is "marbleInit"
		Query(`Delete`, invokeMarbleDelete, p.String(`id`)).
		Invoke(`SetOwner`, invokeMarbleSetOwner, p.String(`id`))
	return router.NewChaincode(r)
}

// ======= Chain code methods

// marbleInit - create a new marble, store into chaincode state
func invokeMarbleInit(c router.Context) (interface{}, error) {
	marble := c.Param(`marble`).(*schema.Marble)
	return marble, c.State().Put(marble)
}

func invokeMarbleDelete(c router.Context) (interface{}, error) {

	return nil, nil
}

func invokeMarbleSetOwner(c router.Context) (interface{}, error) {

	return nil, nil
}

func marbleSetOwner(c router.Context) (interface{}, error) {

	return nil, nil
}
