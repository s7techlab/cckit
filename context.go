package cckit

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"time"
)

type (
	Context interface {
		Stub() shim.ChaincodeStubInterface
		Client() (cid.ClientIdentity, error)
		Response() Response
		Logger() (*shim.ChaincodeLogger)
		Path() string
		State() State
		Time() (time.Time, error)
		Args() Map
		Arg(string) interface{}
		ArgString(string) string
		SetArg(string, interface{})
		Get(string) interface{}
		Set(string, interface{})
	}

	context struct {
		stub       shim.ChaincodeStubInterface
		logger     *shim.ChaincodeLogger
		path       string
		invokeArgs Map
		store      Map
	}
)

func (c *context) Stub() shim.ChaincodeStubInterface {
	return c.stub
}

func (c *context) Client() (cid.ClientIdentity, error) {
	return cid.New(c.Stub())
}

func (c *context) Response() Response {
	return contextResponse{c}
}

func (c *context) Logger() (*shim.ChaincodeLogger) {
	return c.logger
}

func (c *context) Path() (string) {
	return c.path
}

func (c *context) State() State {
	return &stateOp{c}
}

func (c *context) Time() (time.Time, error) {

	txTimestamp, err := c.stub.GetTxTimestamp()
	if err != nil {
		return time.Unix(0, 0), err
	}

	return time.Unix(txTimestamp.GetSeconds(), int64(txTimestamp.GetNanos())), nil
}

func (c *context) Args() Map {
	return c.invokeArgs
}

func (c *context) SetArg(name string, value interface{}) {
	if c.invokeArgs == nil {
		c.invokeArgs = make(Map)
	}
	c.invokeArgs[name] = value
}

func (c *context) Arg(name string) interface{} {
	return c.invokeArgs[name]
}

func (c *context) ArgString(name string) string {
	return c.Arg(name).(string)
}

func (c *context) Set(key string, val interface{}) {
	if c.store == nil {
		c.store = make(Map)
	}
	c.store[key] = val
}

func (c *context) Get(key string) interface{} {
	return c.store[key]
}


