package encryption

import (
	"fmt"

	"testing"

	"math/rand"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
		Pre(ArgsDecrypt).
		Init(router.EmptyContextHandler).
		Invoke(`putEncryptedToState`, putEncryptedToState, p.String(`key`), p.String(`value`))
	return router.NewChaincode(r)
}

func putEncryptedToState(c router.Context) (interface{}, error) {

	fmt.Println(c.ArgString(`key`))
	fmt.Println(c.ArgString(`value`))

	return nil, nil
}

var _ = Describe(`Router`, func() {

	//Create chaincode mock
	encryptedCC := testcc.NewMockStub(`routerEncrypted`, NewEncryptedCC())

	encKey := make([]byte, 32)
	rand.Seed(time.Now().UnixNano())
	rand.Read(encKey)

	//fmt.Println(encKey)

	Describe("Encrypted", func() {

		It("Allow anyone to invoke ping method", func() {

			args, err := EncryptArgs(encKey, `putEncryptedToState`, `key1`, `value1`)
			if err != nil {
				panic(err)
			}

			expectcc.ResponseOk(encryptedCC.Init())

			//fmt.Println(string(args[0]))
			fmt.Println(encryptedCC.WithTransient(TransientMapWithKey(encKey)).InvokeBytes(args...).Message)

			////invoke chaincode method from authority actor
			//pingInfo := expectcc.PayloadIs(cc.From(invokerIdentity).Invoke(FuncPing), &PingInfo{}).(PingInfo)

			//Expect(pingInfoEvent.InvokerCert).To(Equal(invokerIdentity.GetPEM()))
		})

	})
})
