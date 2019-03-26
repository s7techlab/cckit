package encryption_test

import (
	"encoding/base64"
	"math/rand"
	"testing"
	"time"

	"github.com/s7techlab/cckit/extensions/encryption/testdata"

	"github.com/hyperledger/fabric/protos/peer"

	"github.com/s7techlab/cckit/state/mapping"

	"github.com/s7techlab/cckit/examples/payment/schema"

	examplecert "github.com/s7techlab/cckit/examples/cert"
	"github.com/s7techlab/cckit/examples/payment"
	"github.com/s7techlab/cckit/extensions/encryption"
	testcc "github.com/s7techlab/cckit/testing"
	expectcc "github.com/s7techlab/cckit/testing/expect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestEncryption(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Router suite")
}

var (
	encryptOnDemandPaymentCC            *testcc.MockStub
	encryptPaymentCC                    *testcc.MockStub
	encryptPaymentCCWithEncStateContext *testcc.MockStub
	externalCC                          *testcc.MockStub

	encCCInvoker *encryption.MockStub

	actors testcc.Identities

	pType  = `SALE`
	encKey []byte

	// fixtures
	pId1           = `id-1`
	pAmount1 int32 = 111

	pId2           = `id-2`
	pAmount2 int32 = 222

	pId3           = `id-3`
	pAmount3 int32 = 333

	encryptedPType, encryptedPId1 []byte
	payment1                      *schema.Payment
	encPayment1                   []byte
	paymentMapper                 mapping.StateMapper
	encPaymentNamespace           []byte
	err                           error
)

var _ = Describe(`Router`, func() {

	// Create encode key. In real case it can be calculated via ECDH
	encKey = make([]byte, 32)
	rand.Seed(time.Now().UnixNano())
	rand.Read(encKey)

	//Create chaincode mock
	encryptOnDemandPaymentCC = testcc.NewMockStub(`paymentsEncOnDemand`, payment.NewEncryptOnDemandPaymentCC())
	encryptPaymentCC = testcc.NewMockStub(`paymentsEnc`, payment.NewEncryptPaymentCC())

	encryptPaymentCCWithEncStateContext = testcc.NewMockStub(`paymentsEncWithContext`, payment.NewEncryptedPaymentCCWithEncStateContext())

	encCCInvoker = encryption.NewMockStub(encryptPaymentCCWithEncStateContext, encKey)
	encCCInvoker.DecryptInvokeResponse = true

	externalCC = testcc.NewMockStub(`external`,
		testdata.NewExternaldCC(`paymentsEncWithContext`, `payment-channel`))

	// external cc have access to encrypted payment chaincode
	externalCC.MockPeerChaincode(
		`paymentsEncWithContext/payment-channel`,
		encryptPaymentCCWithEncStateContext)

	BeforeSuite(func() {

		actors, err = testcc.IdentitiesFromFiles(`SOME_MSP`, map[string]string{
			`owner`:   `s7techlab.pem`,
			`someone`: `victor-nosov.pem`}, examplecert.Content)
		Expect(err).To(BeNil())

		encryptedPType, err = encryption.Encrypt(encKey, pType)
		Expect(err).To(BeNil())

		encryptedPId1, err = encryption.Encrypt(encKey, pId1)
		Expect(err).To(BeNil())

		payment1 = &schema.Payment{
			Type:   pType,
			Id:     pId1,
			Amount: pAmount1,
		}

		encPayment1, err = encryption.Encrypt(encKey, payment1)
		Expect(err).To(BeNil())

		paymentMapper, _ = payment.StateMappings.Get(&schema.Payment{})
		Expect(err).To(BeNil())

		// we know that Payment namespace contain only one part
		encPaymentNamespace, err = encryption.Encrypt(encKey, paymentMapper.Namespace()[0])
		Expect(err).To(BeNil())
	})

	Describe("Encrypting in demand with state mapping", func() {

		It("Allow to init encrypt-on-demand payment chaincode", func() {
			expectcc.ResponseOk(encryptOnDemandPaymentCC.Init())
		})

		It("Disallow to create encrypted payment providing unencrypted arguments", func() {
			expectcc.ResponseError(
				encryptOnDemandPaymentCC.WithTransient(encryption.TransientMapWithKey(encKey)).
					Invoke(`paymentCreate`, pType, pId1, pAmount1), `args: decryption error`)
		})

		It("Allow to create encrypted payment", func() {
			// encrypt all arguments
			args, err := encryption.EncryptArgs(encKey, `paymentCreate`, pType, pId1, pAmount1)
			Expect(err).To(BeNil())
			Expect(len(args)).To(Equal(4))

			Expect(err).To(BeNil())
			Expect(args[1]).To(Equal(encryptedPType))
			// second argument is encoded payment id
			Expect(args[2]).To(Equal(encryptedPId1))

			// invoke chaincode with encoded args and encKey via transientMap, receives encoded payment id
			ccPId := expectcc.PayloadBytes(
				encryptOnDemandPaymentCC.WithTransient(encryption.TransientMapWithKey(encKey)).
					InvokeBytes(args...), encryptedPId1)

			encEvent := <-encryptOnDemandPaymentCC.ChaincodeEventsChannel
			decryptedEvent := encryption.MustDecryptEvent(encKey, encEvent)

			Expect(decryptedEvent.Payload).To(BeEquivalentTo(
				testcc.MustProtoMarshal(&schema.PaymentEvent{Type: pType, Id: pId1, Amount: pAmount1})))

			Expect(decryptedEvent.EventName).To(Equal(`PaymentEvent`))

			decryptedPaymentId, err := encryption.Decrypt(encKey, ccPId)
			Expect(err).To(BeNil())
			Expect(string(decryptedPaymentId)).To(Equal(pId1))
		})

		It("Allow to get encrypted payment by type and id", func() {
			// encrypt all arguments
			args, err := encryption.EncryptArgs(encKey, `paymentGet`, pType, pId1)
			Expect(err).To(BeNil())
			Expect(len(args)).To(Equal(3))

			// Check that value is encrypted in chaincode state - use debugStateGet func
			// without providing key in transient map - so we need to provide encrypted key
			// and cause we dont't require key - state also returns unencrypted
			expectcc.PayloadBytes(encryptOnDemandPaymentCC.Invoke(`debugStateGet`, []string{
				base64.StdEncoding.EncodeToString(encPaymentNamespace),
				base64.StdEncoding.EncodeToString(encryptedPType),
				base64.StdEncoding.EncodeToString(encryptedPId1)}), encPayment1)

			//returns unencrypted
			paymentFromCC := expectcc.PayloadIs(
				encryptOnDemandPaymentCC.WithTransient(encryption.TransientMapWithKey(encKey)).InvokeBytes(args...), &schema.Payment{}).(*schema.Payment)

			Expect(paymentFromCC).To(Equal(payment1))
		})

		It("Allow to get encrypted payments by type as unencrypted values", func() {
			// MockInvoke sets key in transient map and encrypt input arguments
			payments := expectcc.PayloadIs(encryption.MockInvoke(encryptOnDemandPaymentCC, encKey, `paymentList`, pType), &schema.PaymentList{}).(*schema.PaymentList)

			Expect(payments.Items).To(HaveLen(1))
			// Returned value is not encrypted
			Expect(payments.Items[0].Id).To(Equal(pId1))
		})

		It("Allow to invoke with non encrypted data", func() {
			expectcc.ResponseOk(encryptOnDemandPaymentCC.Invoke(`paymentCreate`, pType, pId2, pAmount2))
		})

		It("Allow to get non encrypted payment by type and id", func() {
			//returns unencrypted
			paymentFromCC := expectcc.PayloadIs(
				encryptOnDemandPaymentCC.Query(`paymentGet`, pType, pId2), &schema.Payment{}).(*schema.Payment)

			Expect(string(paymentFromCC.Id)).To(Equal(pId2))
		})
	})

	Describe("Encrypting required", func() {

		It("Allow to init encrypt payment chaincodes", func() {
			expectcc.ResponseOk(encryptPaymentCC.Init())
		})

		It("Disallow to create payment without providing key in encryptPaymentCC ", func() {
			expectcc.ResponseError(encryptPaymentCC.Invoke(`paymentCreate`, pType, pId3, pAmount3),
				encryption.ErrKeyNotDefinedInTransientMap)
		})

		It("Allow to create payment providing key in encryptPaymentCC ", func() {
			// encode all arguments
			expectcc.ResponseOk(encryption.MockInvoke(encryptPaymentCC, encKey, `paymentCreate`, pType, pId3, pAmount3))
		})
	})

	Describe("Encrypted state context", func() {
		It("Allow to init encrypt payment chaincodes", func() {
			expectcc.ResponseOk(encryptPaymentCCWithEncStateContext.WithTransient(encryption.TransientMapWithKey(encKey)).Init())
		})
		//
		It("Allow to create payment providing key in encryptPaymentCC ", func(done Done) {
			events := encryptPaymentCCWithEncStateContext.EventSubscription()

			responsePayment := expectcc.PayloadIs(
				// encCCInvoker encrypts args before passing to cc invoke and pass key in transient map
				encCCInvoker.From(actors[`owner`]).Invoke(`paymentCreate`, pType, pId1, pAmount3),
				&schema.Payment{}).(*schema.Payment)

			// we use encryption.MockStub DecryptInvokeResponse feature
			Expect(responsePayment.Id).To(Equal(pId1))
			Expect(responsePayment.Type).To(Equal(pType))
			Expect(responsePayment.Amount).To(Equal(pAmount3))

			//event name and payload is encrypted with key
			Expect(<-events).To(BeEquivalentTo(encryption.MustEncryptEvent(encKey, &peer.ChaincodeEvent{
				EventName: `PaymentEvent`,
				Payload: testcc.MustProtoMarshal(&schema.PaymentEvent{
					Type:   pType,
					Id:     pId1,
					Amount: pAmount3,
				}),
			})))

			close(done)
		}, 0.2)

		It("Allow to get payment by type and id", func() {
			// encCCInvoker encrypts args before passing to cc invoke and pass key in transient map
			paymentFromCC := expectcc.PayloadIs(encCCInvoker.From(actors[`owner`]).Query(`paymentGet`, pType, pId1), &schema.Payment{}).(*schema.Payment)

			//returned payload is unencrypted
			Expect(paymentFromCC.Id).To(Equal(pId1))
			Expect(paymentFromCC.Type).To(Equal(pType))
			Expect(paymentFromCC.Amount).To(Equal(pAmount3))
		})

		It("Disallow to get payment providing key using debugStateGet", func() {
			// we didn't provide encrypting key,
			// chaincode use ArgsDecrypt middleware, requiring key in transient map
			// but we add exception for method debugStateGet
			// for all chaincode methods, except stateGet @see example/payments/cc_enc_context_with_mapping.go
			expectcc.ResponseOk(encCCInvoker.MockStub.Invoke(`debugStateGet`, []string{
				base64.StdEncoding.EncodeToString(encPaymentNamespace),
				base64.StdEncoding.EncodeToString(encryptedPType),
				base64.StdEncoding.EncodeToString(encryptedPId1)}))
		})

		It("Allow to get encrypted payments by type as unencrypted values", func() {
			paymentsFromCC := expectcc.PayloadIs(encCCInvoker.Query(`paymentList`, pType), &schema.PaymentList{}).(*schema.PaymentList)

			Expect(paymentsFromCC.Items).To(HaveLen(1))
			// Returned value is not encrypted
			Expect(paymentsFromCC.Items[0].Id).To(Equal(pId1))
			Expect(paymentsFromCC.Items[0].Type).To(Equal(pType))
			Expect(paymentsFromCC.Items[0].Amount).To(Equal(pAmount3))
		})

		It("Disallow to get payment by type and id without providing encrypting key in transient map", func() {
			expectcc.ResponseError(encCCInvoker.MockStub.From(actors[`owner`]).Query(`paymentGet`, pType, pId1),
				encryption.ErrKeyNotDefinedInTransientMap)
		})

		It("Allow to get payment by type and id", func() {
			//returns unencrypted
			paymentFromCC := expectcc.PayloadIs(encCCInvoker.From(actors[`owner`]).Query(`paymentGet`, pType, pId1),
				&schema.Payment{}).(*schema.Payment)
			Expect(paymentFromCC.Id).To(Equal(pId1))
			Expect(paymentFromCC.Type).To(Equal(pType))
			Expect(paymentFromCC.Amount).To(Equal(pAmount3))
		})

		It("Allow to get payment via external chaincode", func() {

			paymentFromExtCC := expectcc.PayloadIs(externalCC.WithTransient(encryption.
				TransientMapWithKey(encKey)).Query(`checkPayment`, pType, pId1), &schema.Payment{}).(*schema.Payment)
			Expect(paymentFromExtCC.Id).To(Equal(pId1))
			Expect(paymentFromExtCC.Type).To(Equal(pType))
			Expect(paymentFromExtCC.Amount).To(Equal(pAmount3))
		})

	})
})
