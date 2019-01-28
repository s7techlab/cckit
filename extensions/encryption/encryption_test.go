package encryption_test

import (
	"encoding/base64"
	"math/rand"
	"testing"
	"time"

	"github.com/s7techlab/cckit/identity"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	examplecert "github.com/s7techlab/cckit/examples/cert"
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

	actors identity.Actors

	pType  = `SALE`
	encKey []byte

	// fixtures
	pId1     = `id-1`
	pAmount1 = 111

	pId2     = `id-2`
	pAmount2 = 222

	pId3     = `id-3`
	pAmount3 = 333

	encryptedPType, encryptedPId1 []byte
	payment1                      encryption.Payment
	encPayment1                   []byte
	err                           error
)

var _ = Describe(`Router`, func() {

	// Create encode key. In real case it can be calculated via ECDH
	encKey = make([]byte, 32)
	rand.Seed(time.Now().UnixNano())
	rand.Read(encKey)

	//Create chaincode mock
	encryptOnDemandPaymentCC = testcc.NewMockStub(`paymentsEncOnDemand`, encryption.NewEncryptOnDemandPaymentCC())
	encryptPaymentCC = testcc.NewMockStub(`paymentsEnc`, encryption.NewEncryptPaymentCC())
	encryptPaymentCCWithEncStateContext = testcc.NewMockStub(`paymentsEnc`, encryption.NewEncryptedPaymentCCWithEncStateContext())
	encCCInvoker := encryption.NewMockStub(encryptPaymentCCWithEncStateContext, encKey)

	BeforeSuite(func() {

		actors, err = identity.ActorsFromPemFile(`SOME_MSP`, map[string]string{
			`owner`:   `s7techlab.pem`,
			`someone`: `victor-nosov.pem`}, examplecert.Content)
		Expect(err).To(BeNil())

		encryptedPType, err = encryption.Encrypt(encKey, pType)
		Expect(err).To(BeNil())

		encryptedPId1, err = encryption.Encrypt(encKey, pId1)
		Expect(err).To(BeNil())

		payment1 = encryption.Payment{
			Type:   pType,
			Id:     pId1,
			Amount: pAmount1,
		}

		encPayment1, err = encryption.Encrypt(encKey, payment1)
		Expect(err).To(BeNil())
	})

	Describe("Encrypting in demand", func() {

		It("Allow to init encrypt-on-demand payment chaincode", func() {
			expectcc.ResponseOk(encryptOnDemandPaymentCC.Init())
		})

		It("Disallow to create encrypted payment providing unencrypted arguments", func() {
			expectcc.ResponseError(
				encryptOnDemandPaymentCC.WithTransient(encryption.TransientMapWithKey(encKey)).
					Invoke(`paymentCreate`, pType, pId1, pAmount1), `args: decryption error`)
		})

		It("Allow to create encrypted payment", func() {
			// encode all arguments
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

			//Check that value is encrypted in chaincode state - use debugStateGet func
			// without providing key in tramsient map
			expectcc.PayloadBytes(encryptOnDemandPaymentCC.Invoke(`debugStateGet`, []string{
				base64.StdEncoding.EncodeToString(encryptedPType),
				base64.StdEncoding.EncodeToString(encryptedPId1)}), encPayment1)

			//returns unencrypted
			paymentFromCC := expectcc.PayloadIs(
				encryptOnDemandPaymentCC.WithTransient(encryption.TransientMapWithKey(encKey)).InvokeBytes(args...), &encryption.Payment{}).(encryption.Payment)

			Expect(paymentFromCC).To(Equal(payment1))
		})

		It("Allow to get encrypted payments by type as unencrypted values", func() {
			// MockInvoke sets key in transient map and encrypt input arguments
			payments := expectcc.PayloadIs(encryption.MockInvoke(encryptOnDemandPaymentCC, encKey, `paymentList`, pType), &[]encryption.Payment{}).([]encryption.Payment)

			Expect(len(payments)).To(Equal(1))
			// Returned value is not encrypted
			Expect(payments[0].Id).To(Equal(pId1))
		})

		It("Allow to create non encrypted payment encryptOnDemandPaymentCC", func() {
			expectcc.ResponseOk(encryptOnDemandPaymentCC.Invoke(`paymentCreate`, pType, pId2, pAmount2))
		})

		It("Allow to get non encrypted payment by type and id", func() {
			//returns unencrypted
			paymentFromCC := expectcc.PayloadIs(
				encryptOnDemandPaymentCC.Query(`paymentGet`, pType, pId2), &encryption.Payment{}).(encryption.Payment)

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
			expectcc.ResponseOk(encryptPaymentCCWithEncStateContext.Init())
		})

		It("Allow to create payment providing key in encryptPaymentCC ", func() {
			// encode all arguments
			expectcc.ResponseOk(encryption.MockInvoke(encryptPaymentCCWithEncStateContext, encKey, `paymentCreate`, pType, pId1, pAmount3))

			//Check that value is encrypted in chaincode state - use debugStateGet func
			expectcc.PayloadBytes(encryptOnDemandPaymentCC.Invoke(`debugStateGet`, []string{
				base64.StdEncoding.EncodeToString(encryptedPType),
				base64.StdEncoding.EncodeToString(encryptedPId1)}), encPayment1)
		})

		It("Allow to get encrypted payments by type as unencrypted values", func() {
			payments := expectcc.PayloadIs(encCCInvoker.Invoke(`paymentList`, pType), &[]encryption.Payment{}).([]encryption.Payment)

			Expect(len(payments)).To(Equal(1))
			// Returned value is not encrypted
			Expect(payments[0].Id).To(Equal(pId1))
		})

		It("Allow to get payment by type and id", func() {
			//returns unencrypted
			paymentFromCC := expectcc.PayloadIs(encCCInvoker.From(actors[`owner`]).Query(`paymentGet`, pType, pId1), &encryption.Payment{}).(encryption.Payment)
			Expect(string(paymentFromCC.Id)).To(Equal(pId1))
		})
	})
})
