package gateway

import (
	"github.com/hyperledger/fabric-chaincode-go/shim"

	"github.com/s7techlab/cckit/state"
)

// ChaincodeClient for querying external chaincode from chaincodes
type ChaincodeClient interface {
	Query(stub shim.ChaincodeStubInterface, fn string, args []interface{}, target interface{}) (interface{}, error)
}

type ClientOpt func(*chaincodeClient)

type chaincodeClient struct {
	Channel   string
	Chaincode string
}

func NewChaincodeClient(channelName, chaincodeName string, opts ...OptFunc) *chaincodeClient {
	c := &chaincodeClient{
		Channel:   channelName,
		Chaincode: chaincodeName,
	}

	return c
}

func (c *chaincodeClient) Query(stub shim.ChaincodeStubInterface, fn string, args []interface{}, target interface{}) (interface{}, error) {

	// if target chaincode is encrypted we only can invoke `stateGet` function
	// for example < encrypted(`orgGet`),  encrypted( &schema.OrganizationId { Id : `123` }) > will be <  `stateGet`, []string { key } >
	// we know target ( &schema.Organization ), know input parameter type ( &schema.OrganizationId )
	// if target has primary key with this type or uniq key with this type - we create state,Ke
	//if 1==2 {
	//	fn = `stateGet`
	//
	//}

	return state.InvokeChaincode(stub, c.Chaincode, append([]interface{}{fn}, args...), c.Channel, target)
}
