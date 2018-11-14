package encryption_test

import (
	"testing"

	"math/rand"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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

func NewEncryptedCC() *router.Chaincode {
	r := router.New(`encrypted`).
		Pre(encryption.ArgsDecryptIfKeyProvided).
		Init(router.EmptyContextHandler).
		Invoke(`putEncryptedToState`, putEncryptedToState, p.String(`key`), p.String(`value`))

	debug.AddHandlers(r, `debug`)
	return router.NewChaincode(r)
}

func putEncryptedToState(c router.Context) (interface{}, error) {
	tm, _ := c.Stub().GetTransient()
	es := encryption.State(c, tm[encryption.TransientMapKey])

	v := c.ArgString(`value`)
	encValue, _ := encryption.Encrypt(tm[encryption.TransientMapKey], []byte(v))

	return encValue, es.Put(c.ArgString(`key`), v)
}

var _ = Describe(`Router`, func() {

	//Create chaincode mock
	encryptedCC := testcc.NewMockStub(`routerEncrypted`, NewEncryptedCC())

	encKey := make([]byte, 32)
	rand.Seed(time.Now().UnixNano())
	rand.Read(encKey)

	//fmt.Println(encKey)

	Describe("Encrypted", func() {

		It("Allow to put encoded data to state", func() {

			key1 := `key1`
			value1 := `value1`

			args, err := encryption.EncryptArgs(encKey, `putEncryptedToState`, key1, value1)
			if err != nil {
				panic(err)
			}

			expectcc.ResponseOk(encryptedCC.Init())
			ccRes := encryptedCC.WithTransient(encryption.TransientMapWithKey(encKey)).InvokeBytes(args...).Payload
			decryptedRes, err := encryption.Decrypt(encKey, ccRes)
			if err != nil {
				panic(err)
			}

			Expect(string(decryptedRes)).To(Equal(value1))

			getRes := encryptedCC.Invoke(`debugStateGet`, []string{key1}).Payload
			decryptedRes, err = encryption.Decrypt(encKey, getRes)
			if err != nil {
				panic(err)
			}

			Expect(string(decryptedRes)).To(Equal(value1))
		})

	})
})
