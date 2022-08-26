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

// NewEncryptOnDemandPaymentCC chaincode WITH schema mapping
// and encrypting data on demand (if encrypting key is provided in transient map)
func NewEncryptOnDemandPaymentCC() *router.Chaincode {
	r := router.New(`encrypted-on-demand`).
		Pre(encryption.ArgsDecryptIfKeyProvided). //  encrypted args optional - key can be provided in transient map
		Init(router.EmptyContextHandler)

	debug.AddHandlers(r, `debug`)

	r.Use(m.MapStates(StateMappings)) // use state mappings
	r.Use(m.MapEvents(EventMappings)) // use state mappings

	r.Group(`payment`).
		Invoke(`Create`, invokePaymentCreateManualEncryptWithMapping, p.String(`type`), p.String(`id`), p.Int(`amount`)).
		Query(`List`, queryPaymentsWithMapping, p.String(`type`)).
		Query(`Get`, queryPaymentWithMapping, p.String(`type`), p.String(`id`))

	return router.NewChaincode(r)
}

// Payment creation chaincode function handler - can be used in encryption and no encryption mode
func invokePaymentCreateManualEncryptWithMapping(c router.Context) (interface{}, error) {
	var (
		paymentType   = c.ParamString(`type`)
		paymentId     = c.ParamString(`id`)
		paymentAmount = c.ParamInt(`amount`)
		s             state.State
		e             state.Event
		err           error
		returnVal     []byte
	)

	// Explicit encrypted state wrapper with key from transient map if key provided
	if s, err = encryption.StateWithTransientKeyIfProvided(c); err != nil {
		return nil, err
	}

	returnVal = []byte(paymentId)
	// return encrypted val if key provided
	if key, _ := encryption.KeyFromTransient(c); key != nil {
		returnVal, err = encryption.Encrypt(key, paymentId)
	}

	if e, err = encryption.EventWithTransientKeyIfProvided(c); err != nil {
		return nil, err
	}

	if err = e.Set(&schema.PaymentEvent{Type: paymentType, Id: paymentId, Amount: int32(paymentAmount)}); err != nil {
		return nil, err
	}
	return returnVal, s.Put(&schema.Payment{Type: paymentType, Id: paymentId, Amount: int32(paymentAmount)})
}

func queryPaymentsWithMapping(c router.Context) (interface{}, error) {
	var (
		//paymentType = c.ParamString(`type`)
		s   state.State
		err error
	)
	// Explicit encrypted state wrapper with key from transient map if key provided
	if s, err = encryption.StateWithTransientKeyIfProvided(c); err != nil {
		return nil, err
	}
	return s.List(&schema.Payment{})
}

// Payment query chaincode function handler - can be used in encryption and no encryption mode
func queryPaymentWithMapping(c router.Context) (interface{}, error) {
	var (
		paymentType = c.ParamString(`type`)
		paymentId   = c.ParamString(`id`)
		s           state.State
		err         error
	)
	if s, err = encryption.StateWithTransientKeyIfProvided(c); err != nil {
		return nil, err
	}
	return s.Get(&schema.Payment{Type: paymentType, Id: paymentId})
}
