package encryption_test

import (
	"encoding/base64"
	"testing"

	"math/rand"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/s7techlab/cckit/convert"
	"github.com/s7techlab/cckit/extensions/debug"
	"github.com/s7techlab/cckit/extensions/encryption"
	"github.com/s7techlab/cckit/router"
	p "github.com/s7techlab/cckit/router/param"
	testcc "github.com/s7techlab/cckit/testing"
	expectcc "github.com/s7techlab/cckit/testing/expect"
)

func TestEncryption(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Router suite")
}

type Payment struct {
	Type   string
	Id     string
	Amount int
}

func (p Payment) Key() ([]string, error) {
	return []string{p.Type, p.Id}, nil
}

func NewEncryptedPaymentCC() *router.Chaincode {
	r := router.New(`encrypted`).
		Pre(encryption.ArgsDecryptIfKeyProvided).
		Init(router.EmptyContextHandler)

	debug.AddHandlers(r, `debug`)

	r.Group(`payment`).
		Invoke(`Create`, invokePaymentCreate, p.String(`type`), p.String(`id`), p.Int(`amount`)).
		Query(`List`, queryPayments, p.String(`type`)).
		Query(`Get`, queryPayment, p.String(`type`), p.String(`id`))

	return router.NewChaincode(r)
}

func invokePaymentCreate(c router.Context) (interface{}, error) {
	paymentType := c.ArgString(`type`)
	paymentId := c.ArgString(`id`)
	paymentAmount := c.ArgInt(`amount`)

	// use encryption key from transient map
	key, err := encryption.KeyFromTransient(c)
	if err != nil {
		return nil, err
	}

	es, err := encryption.State(c, key)
	if err != nil {
		return nil, err
	}

	encId, _ := encryption.Encrypt(key, paymentId)

	// use encoded state to put encoded amount with encoded key
	return encId, es.Put(&Payment{Type: paymentType, Id: paymentId, Amount: paymentAmount})
}

func queryPayments(c router.Context) (interface{}, error) {
	paymentType := c.ArgString(`type`)
	es, err := encryption.StateWithTransientKey(c)
	if err != nil {
		return nil, err
	}
	return es.List(paymentType, &Payment{})
}

func queryPayment(c router.Context) (interface{}, error) {
	paymentType := c.ArgString(`type`)
	paymentId := c.ArgString(`id`)

	es, err := encryption.StateWithTransientKey(c)
	if err != nil {
		return nil, err
	}

	return es.Get([]string{paymentType, paymentId}, &Payment{})
}

var _ = Describe(`Router`, func() {

	//Create chaincode mock
	encryptedCC := testcc.NewMockStub(`payments`, NewEncryptedPaymentCC())

	// Create encode key. In real case it can be calculated via ECDH
	encKey := make([]byte, 32)
	rand.Seed(time.Now().UnixNano())
	rand.Read(encKey)

	Describe("Encrypted", func() {

		pType := `SALE`
		pId := `some-id`
		pAmount := 222

		encryptedType, err := encryption.Encrypt(encKey, pType)
		if err != nil {
			panic(err)
		}

		encryptedPaymentId, err := encryption.Encrypt(encKey, pId)
		if err != nil {
			panic(err)
		}

		It("Allow to init payment chaincode", func() {
			expectcc.ResponseOk(encryptedCC.Init())
		})

		It("Allow to create encoded payment", func() {

			// encode all arguments
			args, err := encryption.EncryptArgs(encKey, `paymentCreate`, pType, pId, pAmount)
			Expect(err).To(BeNil())
			Expect(len(args)).To(Equal(4))

			Expect(err).To(BeNil())

			Expect(args[1]).To(Equal(encryptedType))
			// second argument is encoded payment id
			Expect(args[2]).To(Equal(encryptedPaymentId))

			// invoke chaincode with encoded args and encKey via transientMap, recieves encoded payment id
			ccPId := expectcc.PayloadBytes(
				encryptedCC.WithTransient(encryption.TransientMapWithKey(encKey)).InvokeBytes(args...), encryptedPaymentId)

			decryptedPaymentId, err := encryption.Decrypt(encKey, ccPId)
			Expect(err).To(BeNil())

			Expect(string(decryptedPaymentId)).To(Equal(pId))
		})

		It("Allow to get encoded paymentby type and ids", func() {

			// encode all arguments
			args, err := encryption.EncryptArgs(encKey, `paymentGet`, pType, pId)
			Expect(err).To(BeNil())
			Expect(len(args)).To(Equal(3))

			payment := Payment{
				Type:   pType,
				Id:     pId,
				Amount: pAmount,
			}

			bb, err := convert.ToBytes(payment)
			Expect(err).To(BeNil())

			encPayment, err := encryption.Encrypt(encKey, bb)
			Expect(err).To(BeNil())

			//Check that value is encrypted in chaincode state
			expectcc.PayloadBytes(encryptedCC.Invoke(`debugStateGet`, []string{
				base64.StdEncoding.EncodeToString(encryptedType),
				base64.StdEncoding.EncodeToString(encryptedPaymentId)}), encPayment)

			//returns unencrypted
			paymentFromCC := expectcc.PayloadIs(
				encryptedCC.WithTransient(encryption.TransientMapWithKey(encKey)).InvokeBytes(args...), &Payment{}).(Payment)

			Expect(paymentFromCC).To(Equal(payment))
		})

		It("Allow to get encoded payments by type", func() {
			args, err := encryption.EncryptArgs(encKey, `paymentList`, pType)
			Expect(err).To(BeNil())

			payments := expectcc.PayloadIs(encryptedCC.WithTransient(encryption.TransientMapWithKey(encKey)).InvokeBytes(args...), &[]Payment{}).([]Payment)

			Expect(len(payments)).To(Equal(1))
			Expect(payments[0].Id).To(Equal(pId))
		})

	})
})
