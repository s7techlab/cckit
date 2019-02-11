package encryption_test

import (
	"encoding/base64"
	"math/rand"
	"testing"
	"time"

	"github.com/s7techlab/cckit/state/mapping"

	"github.com/s7techlab/cckit/examples/payment/schema"

	"github.com/s7techlab/cckit/identity"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	examplecert "github.com/s7techlab/cckit/examples/cert"
	"github.com/s7techlab/cckit/examples/payment"
	"github.com/s7techlab/cckit/extensions/encryption"
	testcc "github.com/s7techlab/cckit/testing"
	expectcc "github.com/s7techlab/cckit/testing/expect"
)

func TestEncryption(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Router suite")
}

var (
	encryptOnDemandPaymentCC            *testcc.MockStub
	encryptPaymentCC                    *testcc.MockStub
	encryptPaymentCCWithEncStateContext *testcc.MockStub

	encCCInvoker *encryption.MockStub

	actors identity.Actors

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

	BeforeSuite(func() {

		actors, err = identity.ActorsFromPemFile(`SOME_MSP`, map[string]string{
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

		// we know than Payment namespace contain only one part
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
			// without providing key in transient map

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
			payments := expectcc.PayloadIs(encryption.MockInvoke(encryptOnDemandPaymentCC, encKey, `paymentList`, pType), &[]schema.Payment{}).([]schema.Payment)

			Expect(len(payments)).To(Equal(1))
			// Returned value is not encrypted
			Expect(payments[0].Id).To(Equal(pId1))
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
		It("Allow to create payment providing key in encryptPaymentCC ", func() {
			// encode all arguments
			expectcc.ResponseOk(encryption.MockInvoke(encryptPaymentCCWithEncStateContext, encKey, `paymentCreate`, pType, pId1, pAmount3))
		})

		It("Allow to get payment by type and id", func() {
			//returns unencrypted
			paymentFromCC := expectcc.PayloadIs(encCCInvoker.From(actors[`owner`]).Query(`paymentGet`, pType, pId1), &schema.Payment{}).(*schema.Payment)
			Expect(string(paymentFromCC.Id)).To(Equal(pId1))
		})

		It("Allow to get encrypted payments by type as unencrypted values", func() {
			payments := expectcc.PayloadIs(encCCInvoker.Invoke(`paymentList`, pType), &[]schema.Payment{}).([]schema.Payment)

			Expect(len(payments)).To(Equal(1))
			// Returned value is not encrypted
			Expect(payments[0].Id).To(Equal(pId1))
		})

	})
})
