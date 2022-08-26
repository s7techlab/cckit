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

// NewEncryptPaymentCC chaincode with required encrypting (encrypting key must be provided in transient map)
// WITHOUT mapping
func NewEncryptPaymentCC() *router.Chaincode {
	r := router.New(`encrypted`).
		Pre(encryption.ArgsDecrypt). //  encrypted args required - key must be provided in transient map
		Init(router.EmptyContextHandler)

	r.Use(m.MapStates(StateMappings)) // use state mapping

	debug.AddHandlers(r, `debug`)

	// same as NewEncryptOnDemandPaymentCC
	r.Group(`payment`).
		Invoke(`Create`, invokePaymentCreateManualEncryptWithoutMapping, p.String(`type`), p.String(`id`), p.Int(`amount`)).
		Query(`List`, queryPaymentsWithoutMapping, p.String(`type`)).
		Query(`Get`, queryPaymentWithoutMapping, p.String(`type`), p.String(`id`))

	return router.NewChaincode(r)
}

// Payment creation chaincode function handler - can be used in encryption and no encryption mode
func invokePaymentCreateManualEncryptWithoutMapping(c router.Context) (interface{}, error) {
	var (
		paymentType   = c.ParamString(`type`)
		paymentId     = c.ParamString(`id`)
		paymentAmount = c.ParamInt(`amount`)
		s             state.State
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

	// manually set key
	return returnVal, s.Put(
		[]string{`PaymentManual`, paymentType, paymentId},
		&schema.Payment{
			Type:   paymentType,
			Id:     paymentId,
			Amount: int32(paymentAmount)},
	)
}

func queryPaymentsWithoutMapping(c router.Context) (interface{}, error) {
	var (
		//paymentType = c.ParamString(`type`)
		s   state.State
		err error
	)
	// Explicit encrypted state wrapper with key from transient map if key provided
	if s, err = encryption.StateWithTransientKeyIfProvided(c); err != nil {
		return nil, err
	}
	// manually set key
	return s.List(`Payment`)
}

// Payment query chaincode function handler - can be used in encryption and no encryption mode
func queryPaymentWithoutMapping(c router.Context) (interface{}, error) {
	var (
		paymentType = c.ParamString(`type`)
		paymentId   = c.ParamString(`id`)
		s           state.State
		err         error
	)
	if s, err = encryption.StateWithTransientKeyIfProvided(c); err != nil {
		return nil, err
	}
	// manually set key
	return s.Get([]string{`Payment`, paymentType, paymentId})
}
