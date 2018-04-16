// Package router provides base router for using in chaincode Invoke function
package router

import (
	"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"os"
	"sort"
)

var (
	errMethodNotFound  = errors.New(`method not found`)
	errNoRoutes        = errors.New(`no routes presented`)
	errArgsNumMismatch = errors.New(`method args count mismatch`)
)

type (
	Map map[string]interface{}

	HandlerFunc    func(Context) peer.Response
	MiddlewareFunc func(HandlerFunc, ...int) HandlerFunc

	PathHandler struct {
		Handler HandlerFunc
	}

	Group struct {
		logger     *shim.ChaincodeLogger
		prefix     string
		middleware []MiddlewareFunc
		methods    map[string]func(stub shim.ChaincodeStubInterface) peer.Response
		handlers   map[string]PathHandler
	}
)

// Handle used for using in CC Invoke function
// Must be called after adding new routes using Add function
func (g *Group) Handle(stub shim.ChaincodeStubInterface) peer.Response {
	fnString, _ := stub.GetFunctionAndParameters()

	if fn, ok := g.methods[fnString]; ok {
		return fn(stub)
	}
	if pathHandler, ok := g.handlers[fnString]; ok {

		g.logger.Debug(`router.invoke: `, fnString)
		h := func(c Context) peer.Response {
			h := pathHandler.Handler
			for i := len(g.middleware) - 1; i >= 0; i-- {
				h = g.middleware[i](h, i)
			}
			return h(c)
		}
		return h(g.Context(fnString, stub))
	}

	g.logger.Error(`router.methodnotfound: `, fnString)
	return shim.Error(errMethodNotFound.Error())
}

func (g *Group) Use(middleware ...MiddlewareFunc) {
	g.middleware = append(g.middleware, middleware...)
}

// Group gets new group using presented path
// New group can be used as independent
func (g *Group) Group(path string) *Group {
	return &Group{
		logger:     g.logger,
		prefix:     g.prefix + path,
		methods:    g.methods,
		handlers:   g.handlers,
		middleware: g.middleware,
	}
}

// Add adds new handler using presented path
// Sets methods handler container
func (g *Group) Add(path string, fn func(stub shim.ChaincodeStubInterface) peer.Response) *Group {
	g.methods[g.prefix+path] = fn
	return g
}

// Query alias for invoke
func (g *Group) Query(path string, handler HandlerFunc, middleware ...MiddlewareFunc) *Group {
	return g.Invoke(path, handler, middleware...)
}

func (g *Group) Invoke(path string, handler HandlerFunc, middleware ...MiddlewareFunc) *Group {
	g.handlers[g.prefix+path] = PathHandler{
		func(context Context) peer.Response {
			h := handler
			for i := len(middleware) - 1; i >= 0; i-- {
				h = middleware[i](h, i)
			}
			return h(context)
		}}
	return g
}

// Routes
// Get ordered string view of routes
func (g *Group) Routes() ([]string, error) {
	rLen := len(g.methods)
	if rLen == 0 {
		return nil, errNoRoutes
	}
	r := make([]string, len(g.methods))
	i := 0
	for k, _ := range g.methods {
		r[i] = k
		i++
	}
	sort.Strings(r)
	return r, nil
}

func (g *Group) Context(path string, stub shim.ChaincodeStubInterface) Context {
	return &context{path: path, stub: stub, logger: g.logger}
}

func New(name string) *Group {

	logger := shim.NewLogger(name)
	loggingLevel, err := shim.LogLevel(os.Getenv(`CORE_CHAINCODE_LOGGING_LEVEL`))
	if err == nil {
		logger.SetLevel(loggingLevel)
	}

	g := new(Group)
	g.logger = logger
	g.methods = make(map[string]func(stub shim.ChaincodeStubInterface) peer.Response)
	g.handlers = make(map[string]PathHandler)

	return g
}
