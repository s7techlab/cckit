package encryption_test

import (
	"encoding/base64"
	"testing"

	"github.com/s7techlab/cckit/convert"

	"math/rand"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/s7techlab/cckit/extensions/encryption"
	testcc "github.com/s7techlab/cckit/testing"
	expectcc "github.com/s7techlab/cckit/testing/expect"
)

func TestEncryption(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Router suite")
}

var (
	encryptOnDemandPaymentCC *testcc.MockStub
	encryptPaymentCC         *testcc.MockStub

	pType  = `SALE`
	encKey []byte

	// fixtures
	pId1     = `id-1`
	pAmount1 = 111

	pId2     = `id-2`
	pAmount2 = 222

	pId3     = `id-3`
	pAmount3 = 333
)

var _ = Describe(`Router`, func() {

	//Create chaincode mock
	encryptOnDemandPaymentCC = testcc.NewMockStub(`paymentsEncOnDemand`, encryption.NewEncryptOnDemandPaymentCC())
	encryptPaymentCC = testcc.NewMockStub(`paymentsEnc`, encryption.NewEncryptPaymentCC())
	// Create encode key. In real case it can be calculated via ECDH
	encKey = make([]byte, 32)
	rand.Seed(time.Now().UnixNano())
	rand.Read(encKey)

	Describe("Encrypting in demand", func() {

		encryptedType, err := encryption.Encrypt(encKey, pType)
		if err != nil {
			panic(err)
		}
		encryptedPaymentId1, err := encryption.Encrypt(encKey, pId1)
		if err != nil {
			panic(err)
		}

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
			Expect(args[1]).To(Equal(encryptedType))
			// second argument is encoded payment id
			Expect(args[2]).To(Equal(encryptedPaymentId1))

			// invoke chaincode with encoded args and encKey via transientMap, receives encoded payment id
			ccPId := expectcc.PayloadBytes(
				encryptOnDemandPaymentCC.WithTransient(encryption.TransientMapWithKey(encKey)).
					InvokeBytes(args...), encryptedPaymentId1)

			decryptedPaymentId, err := encryption.Decrypt(encKey, ccPId)
			Expect(err).To(BeNil())

			Expect(string(decryptedPaymentId)).To(Equal(pId1))
		})

		It("Allow to get encrypted payment by type and id", func() {
			// encrypt all arguments
			args, err := encryption.EncryptArgs(encKey, `paymentGet`, pType, pId1)
			Expect(err).To(BeNil())
			Expect(len(args)).To(Equal(3))

			payment := encryption.Payment{
				Type:   pType,
				Id:     pId1,
				Amount: pAmount1,
			}

			bb, err := convert.ToBytes(payment)
			Expect(err).To(BeNil())

			encPayment, err := encryption.Encrypt(encKey, bb)
			Expect(err).To(BeNil())

			//Check that value is encrypted in chaincode state - use debugStateGet func
			expectcc.PayloadBytes(encryptOnDemandPaymentCC.Invoke(`debugStateGet`, []string{
				base64.StdEncoding.EncodeToString(encryptedType),
				base64.StdEncoding.EncodeToString(encryptedPaymentId1)}), encPayment)

			//returns unencrypted
			paymentFromCC := expectcc.PayloadIs(
				encryptOnDemandPaymentCC.WithTransient(encryption.TransientMapWithKey(encKey)).InvokeBytes(args...), &encryption.Payment{}).(encryption.Payment)

			Expect(paymentFromCC).To(Equal(payment))
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
})
