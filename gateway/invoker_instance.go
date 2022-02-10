package gateway

import (
	"context"
	"fmt"
	"reflect"

	"github.com/hyperledger/fabric-protos-go/peer"

	"github.com/s7techlab/cckit/convert"
)

// ChaincodeInvoker used in generated service gateway code
type (
	ChaincodeInstanceInvoker interface {
		Query(ctx context.Context, fn string, args []interface{}, target interface{}) (interface{}, error)
		Invoke(ctx context.Context, fn string, args []interface{}, target interface{}) (interface{}, error)
	}

	ChaincodeInstanceServiceInvoker struct {
		ChaincodeInstance ChaincodeInstanceServiceServer
	}
)

func NewChaincodeInstanceServiceInvoker(ccInstance ChaincodeInstanceServiceServer) *ChaincodeInstanceServiceInvoker {
	c := &ChaincodeInstanceServiceInvoker{
		ChaincodeInstance: ccInstance,
	}

	return c
}

func (c *ChaincodeInstanceServiceInvoker) Query(
	ctx context.Context, fn string, args []interface{}, target interface{}) (interface{}, error) {

	ccInput, err := ccInput(ctx, fn, args)
	if err != nil {
		return nil, err
	}

	res, err := c.ChaincodeInstance.Query(ctx, &ChaincodeInstanceQueryRequest{
		Input: ccInput,
	})
	if err != nil {
		return nil, err
	}

	return ссOutput(res, target)
}

func (c *ChaincodeInstanceServiceInvoker) Invoke(
	ctx context.Context, fn string, args []interface{}, target interface{}) (interface{}, error) {

	ccInput, err := ccInput(ctx, fn, args)
	if err != nil {
		return nil, err
	}

	res, err := c.ChaincodeInstance.Invoke(ctx, &ChaincodeInstanceInvokeRequest{
		Input: ccInput,
	})
	if err != nil {
		return nil, err
	}

	return ссOutput(res, target)
}

func InvokerArgs(fn string, args []interface{}) ([][]byte, error) {
	argsBytes, err := convert.ArgsToBytes(args...)
	if err != nil {
		return nil, fmt.Errorf(`invoker args: %w`, err)
	}

	return append([][]byte{[]byte(fn)}, argsBytes...), nil
}

func ccInput(ctx context.Context, fn string, args []interface{}) (*ChaincodeInput, error) {
	argsBytes, err := InvokerArgs(fn, args)
	if err != nil {
		return nil, err
	}
	ccInput := &ChaincodeInput{
		Args: argsBytes,
	}

	if ccInput.Transient, err = TransientFromContext(ctx); err != nil {
		return nil, err
	}

	return ccInput, nil
}

func ссOutput(response *peer.Response, target interface{}) (res interface{}, err error) {
	output, err := convert.FromBytes(response.Payload, target)
	if err != nil {
		return nil, fmt.Errorf(`convert output to=%s: %w`, reflect.TypeOf(target), err)
	}

	return output, nil
}
