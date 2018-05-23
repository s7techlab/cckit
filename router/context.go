package router

import (
	"time"

	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type (
	// Context of chaincode invoke
	Context interface {
		Stub() shim.ChaincodeStubInterface
		Client() (cid.ClientIdentity, error)
		Response() Response
		Logger() *shim.ChaincodeLogger
		Path() string
		State() State
		Time() (time.Time, error)
		Args() InterfaceMap
		Arg(string) interface{}
		ArgString(string) string
		ArgBytes(string) []byte
		SetArg(string, interface{})
		Get(string) interface{}
		Set(string, interface{})
	}

	context struct {
		stub       shim.ChaincodeStubInterface
		logger     *shim.ChaincodeLogger
		path       string
		invokeArgs InterfaceMap
		store      InterfaceMap
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
	return c.path
}

func (c *context) State() State {
	return ContextState{c}
}

// Time
func (c *context) Time() (time.Time, error) {
	txTimestamp, err := c.stub.GetTxTimestamp()
	if err != nil {
		return time.Unix(0, 0), err
	}
	return time.Unix(txTimestamp.GetSeconds(), int64(txTimestamp.GetNanos())), nil
}

func (c *context) Args() InterfaceMap {
	return c.invokeArgs
}

func (c *context) SetArg(name string, value interface{}) {
	if c.invokeArgs == nil {
		c.invokeArgs = make(InterfaceMap)
	}
	c.invokeArgs[name] = value
}

func (c *context) Arg(name string) interface{} {
	return c.invokeArgs[name]
}

func (c *context) ArgString(name string) string {
	return c.Arg(name).(string)
}

func (c *context) ArgBytes(name string) []byte {
	return c.Arg(name).([]byte)
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
