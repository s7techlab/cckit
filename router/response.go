package router

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/s7techlab/cckit/response"
)

type Response interface {
	Error(err interface{}) peer.Response
	Success(data interface{}) peer.Response
	Create(data interface{}, err interface{}) peer.Response
}

type contextResponse struct {
	context Context
}

func (c contextResponse) Error(err interface{}) peer.Response {
	res := response.Error(err)
	c.context.Logger().Warning(`router.handle.error: `, c.context.Path(), `, err: `, string(res.Message))
	return res
}

func (c contextResponse) Success(data interface{}) peer.Response {
	res := response.Success(data)
	c.context.Logger().Debug(`router.handle.success: `, c.context.Path(), `, data: `, string(res.Payload))
	return res
}

func (c contextResponse) Create(data interface{}, err interface{}) peer.Response {
	result := response.Create(data, err)

	if result.Status == shim.ERROR {
		return c.Error(result.Message)
	} else {
		return c.Success(result.Payload)
	}
}
