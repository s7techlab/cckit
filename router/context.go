package router

import (
	"time"

	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/s7techlab/cckit/convert"
	"github.com/s7techlab/cckit/state"
)

type (
	// Context of chaincode invoke
	Context interface {
		Stub() shim.ChaincodeStubInterface
		Client() (cid.ClientIdentity, error)
		Response() Response
		Logger() *shim.ChaincodeLogger
		Path() string
		State() state.State
		Time() (time.Time, error)

		ReplaceArgs(args [][]byte) Context // replace args, for usage in preMiddleware
		GetArgs() [][]byte

		// Deprecated: Use Params instead.
		Args() InterfaceMap

		// Deprecated: Use Arg instead.
		Arg(string) interface{}

		// Deprecated: Use ParamString instead.
		ArgString(string) string

		// Deprecated: Use ParamBytes instead.
		ArgBytes(string) []byte

		// Deprecated: Use ParamInt instead.
		ArgInt(string) int

		// Deprecated: Use SetParam instead.
		SetArg(string, interface{})

		// Params returns parameter values.
		Params() InterfaceMap

		// Param returns parameter value.
		Param(name string) interface{}

		// ParamString returns parameter value as string.
		ParamString(name string) string

		// ParamBytes returns parameter value as bytes.
		ParamBytes(name string) []byte

		// ParamInt returns parameter value as bytes.
		ParamInt(name string) int

		// SetParam sets parameter value.
		SetParam(name string, value interface{})

		// Get retrieves data from the context.
		Get(key string) interface{}
		// Set saves data in the context.
		Set(key string, value interface{})

		SetEvent(string, interface{}) error
	}

	context struct {
		stub   shim.ChaincodeStubInterface
		logger *shim.ChaincodeLogger
		state  *state.StateImpl
		path   string
		args   [][]byte
		params InterfaceMap
		store  InterfaceMap
	}
)

func (c *context) Stub() shim.ChaincodeStubInterface {
	return c.stub
}

func (c *context) Client() (cid.ClientIdentity, error) {
	return cid.New(c.Stub())
}

func (c *context) Response() Response {
	return ContextResponse{c}
}

func (c *context) Logger() *shim.ChaincodeLogger {
	return c.logger
}

func (c *context) Path() string {
	return string(c.GetArgs()[0])
}

func (c *context) State() state.State {
	if c.state == nil {
		c.state = state.New(c.stub)
	}
	return c.state
}

// Time
func (c *context) Time() (time.Time, error) {
	txTimestamp, err := c.stub.GetTxTimestamp()
	if err != nil {
		return time.Unix(0, 0), err
	}
	return time.Unix(txTimestamp.GetSeconds(), int64(txTimestamp.GetNanos())), nil
}

// ReplaceArgs replace args, for usage in preMiddleware
func (c *context) ReplaceArgs(args [][]byte) Context {
	c.args = args
	return c
}

func (c *context) GetArgs() [][]byte {
	if c.args != nil {
		return c.args
	}
	return c.stub.GetArgs()
}

// Deprecated: Use Params instead.
func (c *context) Args() InterfaceMap {
	return c.Params()
}

func (c *context) Params() InterfaceMap {
	return c.params
}

// Deprecated: Use SetParam instead.
func (c *context) SetArg(name string, value interface{}) {
	c.SetParam(name, value)
}

func (c *context) SetParam(name string, value interface{}) {
	if c.params == nil {
		c.params = make(InterfaceMap)
	}
	c.params[name] = value
}

// Deprecated: Use Param instead.
func (c *context) Arg(name string) interface{} {
	return c.Param(name)
}

func (c *context) Param(name string) interface{} {
	return c.params[name]
}

// Deprecated: Use ParamString instead.
func (c *context) ArgString(name string) string {
	return c.ParamString(name)
}

func (c *context) ParamString(name string) string {
	out, _ := c.Param(name).(string)
	return out
}

// Deprecated: Use ParamBytes instead.
func (c *context) ArgBytes(name string) []byte {
	return c.ParamBytes(name)
}

func (c *context) ParamBytes(name string) []byte {
	out, _ := c.Param(name).([]byte)
	return out
}

// Deprecated: Use ParamInt instead.
func (c *context) ArgInt(name string) int {
	return c.ParamInt(name)
}

func (c *context) ParamInt(name string) int {
	out, _ := c.Param(name).(int)
	return out
}

func (c *context) Set(key string, val interface{}) {
	if c.store == nil {
		c.store = make(InterfaceMap)
	}
	c.store[key] = val
}

func (c *context) Get(key string) interface{} {
	return c.store[key]
}

func (c *context) SetEvent(name string, payload interface{}) error {
	bb, err := convert.ToBytes(payload)
	if err != nil {
		return err
	}
	return c.stub.SetEvent(name, bb)
}
