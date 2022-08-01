package router

import (
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
	"go.uber.org/zap"

	"github.com/s7techlab/cckit/response"
)

// Response chaincode interface
type Response interface {
	Error(err interface{}) peer.Response
	Success(data interface{}) peer.Response
	Create(data interface{}, err interface{}) peer.Response
}

// ContextResponse implementation
type ContextResponse struct {
	context Context
}

// Error response
func (c ContextResponse) Error(err interface{}) peer.Response {
	res := response.Error(err)
	c.context.Logger().Error(`router handler error`, zap.String(`path`, c.context.Path()), zap.String(`message`, res.Message))
	return res
}

// Success response
func (c ContextResponse) Success(data interface{}) peer.Response {
	res := response.Success(data)
	c.context.Logger().Debug(`route handle success`, zap.String(`path`, c.context.Path()), zap.ByteString(`data`, res.Payload))
	return res
}

// Create  returns error response if err != nil
func (c ContextResponse) Create(data interface{}, err interface{}) peer.Response {
	result := response.Create(data, err)

	if result.Status == shim.ERROR {
		return c.Error(result.Message)
	}
	return c.Success(result.Payload)
}
