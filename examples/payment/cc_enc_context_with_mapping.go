package payment

import (
	"github.com/s7techlab/cckit/examples/payment/schema"
	"github.com/s7techlab/cckit/extensions/debug"
	"github.com/s7techlab/cckit/extensions/encryption"
	"github.com/s7techlab/cckit/router"
	p "github.com/s7techlab/cckit/router/param"
	"github.com/s7techlab/cckit/state"
	m "github.com/s7techlab/cckit/state/mapping"
)

// Chaincode WITH schema mapping
// and required encrypting
func NewEncryptedPaymentCCWithEncStateContext() *router.Chaincode {
	r := router.New(`encrypted-with-custom-context`).
		Pre(encryption.ArgsDecryptExcept(`debugStateGet`)). // encrypted args required, except method `stateGet`
		After(encryption.EncryptInvokeResponse()).
		Init(router.EmptyContextHandler)

	r.Use(m.MapStates(StateMappings)) // use state mappings
	r.Use(m.MapEvents(EventMappings)) // use event mappings

	// default Context replaced with EncryptedStateContext only if key is provided in transient map
	r.Use(encryption.EncStateContextIfKeyProvided)

	debug.AddHandlers(r, `debug`)

	r.Group(`payment`).
		// use multiple separate params
		// better way - to use single protobuf "payload" parameter
		Invoke(`Create`, invokePaymentCreateWithDefaultContext, p.String(`type`), p.String(`id`), p.Int(`amount`)).
		Query(`List`, queryPaymentsWithDefaultContext, p.String(`type`)).
		Query(`Get`, queryPaymentWithDefaultContext, p.String(`type`), p.String(`id`))

	return router.NewChaincode(r)
}

func invokePaymentCreateWithDefaultContext(c router.Context) (res interface{}, err error) {
	// params comes unencrypted - "before" middleware decrypts its using key from transient map
	var (
		paymentType   = c.ParamString(`type`)
		paymentId     = c.ParamString(`id`)
		paymentAmount = c.ParamInt(`amount`)
		payment       = &schema.Payment{Type: paymentType, Id: paymentId, Amount: int32(paymentAmount)}
		event         = &schema.PaymentEvent{Type: paymentType, Id: paymentId, Amount: int32(paymentAmount)}
	)

	// State use encryption setting from context, state key sets manually
	if err = c.State().Insert(payment); err != nil {
		return nil, err
	}

	if err = c.Event().Set(event); err != nil {
		return
	}

	// returned value will be placed in ledger - so if we don't want to show in in ledger - we must encrypt it
	return payment, nil
}

func queryPaymentsWithDefaultContext(c router.Context) (interface{}, error) {

	//paymentType := c.ParamString(`type`)
	//namespace, err := c.State().(m.MappedState).MappingNamespace(&schema.Payment{})
	//if err != nil {
	//	return nil, err
	//}
	//return c.State().List(namespace.Append(state.Key { paymentType }), &schema.Payment{})

	// some sugar to previous

	return c.State().(m.MappedState).ListWith(&schema.Payment{}, state.Key{c.ParamString(`type`)})
}

func queryPaymentWithDefaultContext(c router.Context) (interface{}, error) {
	var (
		paymentType = c.ParamString(`type`)
		paymentId   = c.ParamString(`id`)
	)

	return c.State().Get(&schema.Payment{Type: paymentType, Id: paymentId})
}
