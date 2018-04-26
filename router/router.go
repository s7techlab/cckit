// Package router provides base router for using in chaincode Invoke function
package router

import (
	"os"

	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/pkg/errors"
	"github.com/s7techlab/cckit/response"
)

var (
	// ErrMethodNotFound occurs when trying to invoke non existent chaincode method
	ErrMethodNotFound  = errors.New(`method not found`)
	errNoRoutes        = errors.New(`no routes presented`)
	errArgsNumMismatch = errors.New(`method args count mismatch`)
)

type (
	// InterfaceMap map of interfaces
	InterfaceMap map[string]interface{}

	// ContextHandlerFunc use stub context as input parameter
	ContextHandlerFunc func(Context) peer.Response

	// StubHandlerFunc acts as raw chaincode invoke method, accepts stub and returns peer.Response
	StubHandlerFunc func(shim.ChaincodeStubInterface) peer.Response

	// HandlerFunc returns result as interface and error, this is converted to peer.Response via response.Create
	HandlerFunc func(Context) (interface{}, error)

	// MiddlewareFunc middleware for chain code invoke
	ContextMiddlewareFunc func(ContextHandlerFunc, ...int) ContextHandlerFunc

	// MiddlewareFunc middleware for chain code invoke
	MiddlewareFunc func(HandlerFunc, ...int) HandlerFunc

	// Group of chain code functions
	Group struct {
		logger *shim.ChaincodeLogger
		prefix string

		stubHandlers    map[string]StubHandlerFunc
		contextHandlers map[string]ContextHandlerFunc
		handlers        map[string]HandlerFunc

		contextMiddleware []ContextMiddlewareFunc
		middleware        []MiddlewareFunc
	}
)

// Handle used for using in CC Invoke function
// Must be called after adding new routes using Add function
func (g *Group) Handle(stub shim.ChaincodeStubInterface) peer.Response {
	fnString, _ := stub.GetFunctionAndParameters()

	if stubHandler, ok := g.stubHandlers[fnString]; ok {
		g.logger.Debug(`router.stubHandler: `, fnString)
		return stubHandler(stub)
	} else if contextHandler, ok := g.contextHandlers[fnString]; ok {
		g.logger.Debug(`router.contextHandler: `, fnString)

		h := func(c Context) peer.Response {
			h := contextHandler
			for i := len(g.contextMiddleware) - 1; i >= 0; i-- {
				h = g.contextMiddleware[i](h, i)
			}
			return h(c)
		}

		return h(g.Context(fnString, stub))
	} else if handler, ok := g.handlers[fnString]; ok {
		g.logger.Debug(`router.handler: `, fnString)

		h := func(c Context) (interface{}, error) {
			h := handler
			for i := len(g.middleware) - 1; i >= 0; i-- {
				h = g.middleware[i](h, i)
			}
			return h(c)
		}
		return response.Create(h(g.Context(fnString, stub)))
	}

	err := errors.Wrap(fmt.Errorf(`%s: %s`, ErrMethodNotFound, fnString), `chaincode router`)
	g.logger.Error(err)
	return shim.Error(err.Error())
}

// Use middleware function in chain code functions group
func (g *Group) Use(middleware ...MiddlewareFunc) {
	g.middleware = append(g.middleware, middleware...)
}

// Group gets new group using presented path
// New group can be used as independent
func (g *Group) Group(path string) *Group {
	return &Group{
		logger:          g.logger,
		prefix:          g.prefix + path,
		stubHandlers:    g.stubHandlers,
		contextHandlers: g.contextHandlers,
		handlers:        g.handlers,
		middleware:      g.middleware,
	}
}

// StubHandler adds new stub handler using presented path
func (g *Group) StubHandler(path string, fn StubHandlerFunc) *Group {
	g.stubHandlers[g.prefix+path] = fn
	return g
}

// StubHandler adds new stub handler using presented path
func (g *Group) ContextHandler(path string, fn ContextHandlerFunc) *Group {
	g.contextHandlers[g.prefix+path] = fn
	return g
}

// Query alias for invoke
func (g *Group) Query(path string, handler HandlerFunc, middleware ...MiddlewareFunc) *Group {
	return g.Invoke(path, handler, middleware...)
}

// Invoke configure handler and middleware functions for chain code function name
func (g *Group) Invoke(path string, handler HandlerFunc, middleware ...MiddlewareFunc) *Group {
	g.handlers[g.prefix+path] = func(context Context) (interface{}, error) {
		h := handler
		for i := len(middleware) - 1; i >= 0; i-- {
			h = middleware[i](h, i)
		}
		return h(context)
	}
	return g
}

// Context returns chain code invoke context  for provided path and stub
func (g *Group) Context(path string, stub shim.ChaincodeStubInterface) Context {
	return &context{path: path, stub: stub, logger: g.logger}
}

// New group of chain code functions
func New(name string) *Group {

	logger := shim.NewLogger(name)
	loggingLevel, err := shim.LogLevel(os.Getenv(`CORE_CHAINCODE_LOGGING_LEVEL`))
	if err == nil {
		logger.SetLevel(loggingLevel)
	}

	g := new(Group)
	g.logger = logger
	g.stubHandlers = make(map[string]StubHandlerFunc)
	g.contextHandlers = make(map[string]ContextHandlerFunc)
	g.handlers = make(map[string]HandlerFunc)

	return g
}
