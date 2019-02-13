package payment

import (
	"github.com/s7techlab/cckit/examples/payment/schema"
	"github.com/s7techlab/cckit/extensions/debug"
	"github.com/s7techlab/cckit/extensions/encryption"
	"github.com/s7techlab/cckit/router"
	p "github.com/s7techlab/cckit/router/param"
	m "github.com/s7techlab/cckit/state/mapping"
)

// Chaincode WITH schema mapping
// and required encrypting
func NewEncryptedPaymentCCWithEncStateContext() *router.Chaincode {
	r := router.New(`encrypted-with-custom-context`).
		Pre(encryption.ArgsDecrypt).
		Init(router.EmptyContextHandler)

	r.Use(m.MapStates(StateMappings))
	r.Use(encryption.EncStateContext)
	// use state mappings
	// default Context replaced with EncryptedStateContext

	debug.AddHandlers(r, `debug`)

	r.Group(`payment`).
		Invoke(`Create`, invokePaymentCreateWithDefaultContext, p.String(`type`), p.String(`id`), p.Int(`amount`)).
		Query(`List`, queryPaymentsWithDefaultContext, p.String(`type`)).
		Query(`Get`, queryPaymentWithDefaultContext, p.String(`type`), p.String(`id`))

	return router.NewChaincode(r)
}

func invokePaymentCreateWithDefaultContext(c router.Context) (interface{}, error) {
	var (
		paymentType   = c.ParamString(`type`)
		paymentId     = c.ParamString(`id`)
		paymentAmount = c.ParamInt(`amount`)
		returnVal     = []byte(paymentId) // unencrypted
	)
	// State use encryption setting from context
	// and state key set manually
	return returnVal, c.State().Insert(&schema.Payment{Type: paymentType, Id: paymentId, Amount: int32(paymentAmount)})
}

func queryPaymentsWithDefaultContext(c router.Context) (interface{}, error) {

	paymentType := c.ParamString(`type`)
	namespace, err := c.State().(m.MappedState).MappingNamespace(&schema.Payment{})
	if err != nil {
		return nil, err
	}

	// State use encryption setting from context
	return c.State().List(namespace.Add(paymentType), &schema.Payment{})
}

func queryPaymentWithDefaultContext(c router.Context) (interface{}, error) {
	var (
		paymentType = c.ParamString(`type`)
		paymentId   = c.ParamString(`id`)
	)

	return c.State().Get(&schema.Payment{Type: paymentType, Id: paymentId})
}
