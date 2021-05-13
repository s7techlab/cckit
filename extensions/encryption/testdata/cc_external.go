package testdata

import (
	"github.com/s7techlab/cckit/examples/payment"
	"github.com/s7techlab/cckit/examples/payment/schema"

	"github.com/s7techlab/cckit/extensions/encryption"
	"github.com/s7techlab/cckit/router"
	"github.com/s7techlab/cckit/router/param"
)

// Test interaction with external encrypted chaincode (payments)
// see example/payment
func NewExternaldCC(encCCName, channelName string) *router.Chaincode {
	r := router.New(`external`)

	r.
		// returns payment state entry from external encrypted cc
		Query(`checkPayment`, func(c router.Context) (interface{}, error) {
			var (
				paymentType = c.ParamString(`paymentType`)
				paymentId   = c.ParamString(`paymentId`)
				encKey, err = encryption.KeyFromTransient(c)
			)
			if err != nil {
				return nil, err
			}

			// create state key using payments mapping
			paymentKey, err := payment.StateMappings.PrimaryKey(&schema.Payment{Type: paymentType, Id: paymentId})
			if err != nil {
				return nil, err
			}

			// we need to encrypt state key, not all args (method name `debugStateGet` must remain unencrypted )
			encPaymentKey, err := encryption.KeyEncryptor(encKey)(paymentKey)
			if err != nil {
				return nil, err
			}

			// payment state entry returns decrypted
			return encryption.InvokeChaincode(c.Stub(), encKey,
				encCCName, []interface{}{`debugStateGet`, []string(encPaymentKey)}, channelName, &schema.Payment{})
		}, param.String(`paymentType`), param.String(`paymentId`))

	return router.NewChaincode(r)
}
