package cars

import (
	"errors"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/s7techlab/cckit/router"
	p "github.com/s7techlab/cckit/router/param"
)

// NewProxy created chaincode, related to cars chaincode
func NewProxy(carsChannel, carsChaincode string) *router.Chaincode {
	r := router.New(`cars_proxy`) // also initialized logger with "cars_related" prefix

	r.Init(invokeInit)

	r.Query(`carGet`, queryCarProxy, p.String(`id`))

	return router.NewChaincode(r)
}

//
func queryCarProxy(c router.Context) (interface{}, error) {
	var (
		id = c.ParamString(`id`)
	)

	// query external chaincode
	response := c.Stub().InvokeChaincode(`cars`, [][]byte{[]byte(`carGet`), []byte(id)}, `my_channel`)
	if response.Status == shim.ERROR {
		return nil, errors.New(response.Message)
	}
	return response.Payload, nil
}
