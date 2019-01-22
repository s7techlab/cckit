package encryption

import (
	"github.com/s7techlab/cckit/state"

	"github.com/s7techlab/cckit/extensions/debug"
	"github.com/s7techlab/cckit/router"
	p "github.com/s7techlab/cckit/router/param"
)

type Payment struct {
	Type   string
	Id     string
	Amount int
}

func (p Payment) Key() ([]string, error) {
	return []string{p.Type, p.Id}, nil
}

// Chaincode with required encrypting (encrypting key must be provided in transient map)
func NewEncryptPaymentCC() *router.Chaincode {
	r := router.New(`encrypted`).
		Pre(ArgsDecrypt).
		Init(router.EmptyContextHandler)

	debug.AddHandlers(r, `debug`)

	r.Group(`payment`).
		Invoke(`Create`, invokePaymentCreate, p.String(`type`), p.String(`id`), p.Int(`amount`)).
		Query(`List`, queryPayments, p.String(`type`)).
		Query(`Get`, queryPayment, p.String(`type`), p.String(`id`))

	return router.NewChaincode(r)
}

// Chaincode with encrypting data on demand (if encrypting key is provided in transient map)
func NewEncryptOnDemandPaymentCC() *router.Chaincode {
	r := router.New(`encrypted-on-demand`).
		Pre(ArgsDecryptIfKeyProvided).
		Init(router.EmptyContextHandler)

	debug.AddHandlers(r, `debug`)

	r.Group(`payment`).
		Invoke(`Create`, invokePaymentCreate, p.String(`type`), p.String(`id`), p.Int(`amount`)).
		Query(`List`, queryPayments, p.String(`type`)).
		Query(`Get`, queryPayment, p.String(`type`), p.String(`id`))

	return router.NewChaincode(r)
}

func NewEncryptedPaymentCCWithEncStateContext() *router.Chaincode {
	r := router.New(`encrypted-with-custom-context`).
		Pre(ArgsDecrypt).
		Use(EncStateContext). // default Context replaced with EncryptedStateContext
		Init(router.EmptyContextHandler)

	debug.AddHandlers(r, `debug`)

	r.Group(`payment`).
		Invoke(`Create`, invokePaymentCreateWithDefaultContext, p.String(`type`), p.String(`id`), p.Int(`amount`)).
		Query(`List`, queryPaymentsWithDefaultContext, p.String(`type`)).
		Query(`Get`, queryPaymentWithDefaultContext, p.String(`type`), p.String(`id`))

	return router.NewChaincode(r)
}

// Payment creation chaincode function handler - can be used in encryption and no encryption mode
func invokePaymentCreate(c router.Context) (interface{}, error) {
	var (
		paymentType   = c.ArgString(`type`)
		paymentId     = c.ArgString(`id`)
		paymentAmount = c.ArgInt(`amount`)
		s             state.State
		err           error
		returnVal     []byte
	)

	// Explicit encrypted state wrapper with key from transient map if key provided
	if s, err = StateWithTransientKeyIfProvided(c); err != nil {
		return nil, err
	}

	returnVal = []byte(paymentId)
	// return encrypted val if key provided
	if key, _ := KeyFromTransient(c); key != nil {
		returnVal, err = Encrypt(key, paymentId)
	}

	return returnVal, s.Put(&Payment{Type: paymentType, Id: paymentId, Amount: paymentAmount})
}

func invokePaymentCreateWithDefaultContext(c router.Context) (interface{}, error) {
	var (
		paymentType   = c.ArgString(`type`)
		paymentId     = c.ArgString(`id`)
		paymentAmount = c.ArgInt(`amount`)
		returnVal     = []byte(paymentId) // unencrypted
	)
	// State use encryption setting from context
	return returnVal, c.State().Put(&Payment{Type: paymentType, Id: paymentId, Amount: paymentAmount})
}

func queryPayments(c router.Context) (interface{}, error) {
	var (
		paymentType = c.ArgString(`type`)
		s           state.State
		err         error
	)
	// Explicit encrypted state wrapper with key from transient map if key provided
	if s, err = StateWithTransientKeyIfProvided(c); err != nil {
		return nil, err
	}
	return s.List(paymentType, &Payment{})
}

func queryPaymentsWithDefaultContext(c router.Context) (interface{}, error) {
	var (
		paymentType = c.ArgString(`type`)
	)
	// State use encryption setting from context
	return c.State().List(paymentType, &Payment{})
}

// Payment query chaincode function handler - can be used in encryption and no encryption mode
func queryPayment(c router.Context) (interface{}, error) {
	var (
		paymentType = c.ArgString(`type`)
		paymentId   = c.ArgString(`id`)
		s           state.State
		err         error
	)
	if s, err = StateWithTransientKeyIfProvided(c); err != nil {
		return nil, err
	}
	return s.Get([]string{paymentType, paymentId}, &Payment{})
}

func queryPaymentWithDefaultContext(c router.Context) (interface{}, error) {
	var (
		paymentType = c.ArgString(`type`)
		paymentId   = c.ArgString(`id`)
	)

	return c.State().Get([]string{paymentType, paymentId}, &Payment{})
}
