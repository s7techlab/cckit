package gateway

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"

	"github.com/s7techlab/cckit/convert"
)

// ChaincodeInvoker used in generated service gateway code
type (
	ChaincodeQuerier interface {
		Query(ctx context.Context, fn string, args []interface{}, target interface{}) (interface{}, error)
	}
	ChaincodeInvoker interface {
		ChaincodeQuerier
		Invoke(ctx context.Context, fn string, args []interface{}, target interface{}) (interface{}, error)
	}

	ChaincodeInstanceServiceInvoker struct {
		ChaincodeInstance ChaincodeInstanceServiceServer
	}

	CrossChaincodeServiceInvoker struct {
		Locator *ChaincodeLocator
		Stub    shim.ChaincodeStubInterface
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

	return c.prepareOutput(res, target)
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

	return c.prepareOutput(res, target)
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

func InvokerArgs(fn string, args []interface{}) ([][]byte, error) {
	argsBytes, err := convert.ArgsToBytes(args...)
	if err != nil {
		return nil, fmt.Errorf(`invoker args: %w`, err)
	}

	return append([][]byte{[]byte(fn)}, argsBytes...), nil
}
func (c *ChaincodeInstanceServiceInvoker) prepareOutput(response *peer.Response, target interface{}) (res interface{}, err error) {
	output, err := convert.FromBytes(response.Payload, target)
	if err != nil {
		return nil, fmt.Errorf(`convert output to=%s: %w`, reflect.TypeOf(target), err)
	}

	return output, nil
}

func (c *CrossChaincodeServiceInvoker) Query(
	ctx context.Context, fn string, args []interface{}, target interface{}) (interface{}, error) {

	argsBytes, err := InvokerArgs(fn, args)
	if err != nil {
		return nil, err
	}

	response := c.Stub.InvokeChaincode(c.Locator.Chaincode, argsBytes, c.Locator.Channel)
	if response.Status != shim.OK {
		return nil, fmt.Errorf(`cross chaincode=%s, channel=%s invoke: %w`,
			c.Locator.Chaincode, c.Locator.Channel, errors.New(response.Message))
	}

	return convert.FromBytes(response.Payload, target)
}
