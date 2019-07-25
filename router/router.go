// Package router provides base router for using in chaincode Invoke function
package router

import (
	"fmt"
	"os"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/s7techlab/cckit/response"
)

type (
	MethodType string

	// InterfaceMap map of interfaces
	InterfaceMap map[string]interface{}

	// ContextHandlerFunc use stub context as input parameter
	ContextHandlerFunc func(Context) peer.Response

	// StubHandlerFunc acts as raw chaincode invoke method, accepts stub and returns peer.Response
	StubHandlerFunc func(shim.ChaincodeStubInterface) peer.Response

	// HandlerFunc returns result as interface and error, this is converted to peer.Response via response.Create
	HandlerFunc func(Context) (interface{}, error)

	// ContextMiddlewareFunc middleware for ContextHandlerFun
	ContextMiddlewareFunc func(nextOrPrev ContextHandlerFunc, pos ...int) ContextHandlerFunc

	// MiddlewareFunc middleware for HandlerFunc
	MiddlewareFunc func(HandlerFunc, ...int) HandlerFunc

	HandlerMeta struct {
		Hdl  HandlerFunc
		Type MethodType
	}

	// Group of chain code functions
	Group struct {
		logger *shim.ChaincodeLogger
		prefix string

		// mapping chaincode method  => handler
		stubHandlers    map[string]StubHandlerFunc
		contextHandlers map[string]ContextHandlerFunc
		handlers        map[string]*HandlerMeta

		contextMiddleware []ContextMiddlewareFunc
		middleware        []MiddlewareFunc

		preMiddleware   []ContextMiddlewareFunc
		afterMiddleware []MiddlewareFunc
	}

	Router interface {
		HandleInit(shim.ChaincodeStubInterface)
		Handle(shim.ChaincodeStubInterface)
		Query(path string, handler HandlerFunc, middleware ...MiddlewareFunc) Router
		Invoke(path string, handler HandlerFunc, middleware ...MiddlewareFunc) Router
	}
)

const (
	InitFunc                = `init`
	MethodInvoke MethodType = `invoke`
	MethodQuery  MethodType = `query`
)

func (g *Group) buildHandler() ContextHandlerFunc {
	return func(c Context) peer.Response {
		h := g.handleContext
		// build pre part
		for i := len(g.preMiddleware) - 1; i >= 0; i-- {
			h = g.preMiddleware[i](h, i)
		}

		return h(c)
	}
}

// HandleInit handle chaincode init method
func (g *Group) HandleInit(stub shim.ChaincodeStubInterface) peer.Response {
	// Pre context handling middleware
	h := g.buildHandler()

	// add "init" as first arg
	return h(g.Context(stub).ReplaceArgs(append([][]byte{[]byte(InitFunc)}, stub.GetArgs()...)))
}

// Handle used for using in CC Invoke function
// Must be called after adding new routes using Add function
func (g *Group) Handle(stub shim.ChaincodeStubInterface) peer.Response {
	args := stub.GetArgs()
	if len(args) == 0 {
		return response.Error(ErrEmptyArgs)
	}

	h := g.buildHandler()
	return h(g.Context(stub))
}

func (g *Group) handleContext(c Context) peer.Response {

	// handle standard stub handler (accepts StubInterface, returns peer.Response)
	if stubHandler, ok := g.stubHandlers[c.Path()]; ok {
		g.logger.Debug(`router stubHandler: `, c.Path())
		return stubHandler(c.Stub())

		// handle context handler (accepts Context, returns peer.Response)
	} else if contextHandler, ok := g.contextHandlers[c.Path()]; ok {
		g.logger.Debug(`router contextHandler: `, c.Path())
		h := func(c Context) peer.Response {
			h := contextHandler
			for i := len(g.contextMiddleware) - 1; i >= 0; i-- {
				h = g.contextMiddleware[i](h, i)
			}
			return h(c)
		}
		return h(c)
	} else if handlerMeta, ok := g.handlers[c.Path()]; ok {

		g.logger.Debug(`router handler: `, c.Path())
		h := func(c Context) (interface{}, error) {

			c.SetHandler(handlerMeta)
			h := handlerMeta.Hdl
			for i := len(g.middleware) - 1; i >= 0; i-- {
				h = g.middleware[i](h, i)
			}

			for i := 0; i <= len(g.afterMiddleware)-1; i++ {
				h = g.afterMiddleware[i](h, 0)
			}

			return h(c)
		}
		resp := response.Create(h(c))
		if resp.Status != shim.OK {
			g.logger.Errorf(`%s: %s: %s`, ErrHandlerError, c.Path(), resp.Message)
		}
		return resp
	}

	err := fmt.Errorf(`%s: %s`, ErrMethodNotFound, c.Path())
	g.logger.Error(err)
	return shim.Error(err.Error())
}

func (g *Group) Pre(middleware ...ContextMiddlewareFunc) *Group {
	g.preMiddleware = append(g.preMiddleware, middleware...)
	return g
}

func (g *Group) After(middleware ...MiddlewareFunc) *Group {
	g.afterMiddleware = append(g.afterMiddleware, middleware...)
	return g
}

// Use middleware function in chain code functions group
func (g *Group) Use(middleware ...MiddlewareFunc) *Group {
	g.middleware = append(g.middleware, middleware...)
	return g
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

// ContextHandler adds new context handler using presented path
func (g *Group) ContextHandler(path string, fn ContextHandlerFunc) *Group {
	g.contextHandlers[g.prefix+path] = fn
	return g
}

// Query defines handler and middleware for querying chaincode method (no state change, no send to orderer)
func (g *Group) Query(path string, handler HandlerFunc, middleware ...MiddlewareFunc) *Group {
	return g.addHandler(MethodQuery, path, handler, middleware...)
}

// Invoke defines handler and middleware for invoke chaincode method  (state change,  need to send to orderer)
func (g *Group) Invoke(path string, handler HandlerFunc, middleware ...MiddlewareFunc) *Group {
	return g.addHandler(MethodInvoke, path, handler, middleware...)
}

func (g *Group) addHandler(t MethodType, path string, handler HandlerFunc, middleware ...MiddlewareFunc) *Group {
	g.handlers[g.prefix+path] = &HandlerMeta{
		Type: t,
		Hdl: func(context Context) (interface{}, error) {
			h := handler
			for i := len(middleware) - 1; i >= 0; i-- {
				h = middleware[i](h, i)
			}
			return h(context)
		}}
	return g
}

func (g *Group) Init(handler HandlerFunc, middleware ...MiddlewareFunc) *Group {
	return g.Invoke(InitFunc, handler, middleware...)
}

// Context returns chain code invoke context  for provided path and stub
func (g *Group) Context(stub shim.ChaincodeStubInterface) Context {
	return NewContext(stub, g.logger)
}

// New group of chain code functions
func New(name string) *Group {
	g := new(Group)
	g.logger = NewLogger(name)
	g.stubHandlers = make(map[string]StubHandlerFunc)
	g.contextHandlers = make(map[string]ContextHandlerFunc)
	g.handlers = make(map[string]*HandlerMeta)

	return g
}

// NewContext creates new instance of router.Context
func NewContext(stub shim.ChaincodeStubInterface, logger *shim.ChaincodeLogger) *context {
	return &context{
		stub:   stub,
		logger: logger,
	}
}

// NewLogger creates new instance of shim.ChaincodeLogger
func NewLogger(name string) *shim.ChaincodeLogger {
	logger := shim.NewLogger(name)
	loggingLevel, err := shim.LogLevel(os.Getenv(`CORE_CHAINCODE_LOGGING_LEVEL`))
	if err == nil {
		logger.SetLevel(loggingLevel)
	}

	return logger
}
