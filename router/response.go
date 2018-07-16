package router

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
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
	c.context.Logger().Errorf(`%s: %s:%s`, ErrHandlerError, c.context.Path(), res.Message)
	return res
}

// Success response
func (c ContextResponse) Success(data interface{}) peer.Response {
	res := response.Success(data)
	c.context.Logger().Debug(`route handle success: `, c.context.Path(), `, data: `, string(res.Payload))
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
