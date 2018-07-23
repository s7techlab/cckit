package main

import (
	"testing"

	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	testcc "github.com/s7techlab/cckit/testing"
	expectcc "github.com/s7techlab/cckit/testing/expect"
)

func TestInsuranceApp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Insurance app suite")
}

var _ = Describe(`Insurance`, func() {

	//Create chaincode mock
	cc := testcc.NewMockStub(`insurance`, new(SmartContract))

	BeforeSuite(func() {
		// init chaincode
		expectcc.ResponseOk(cc.Init()) // init chaincode
	})

	Describe("Contract", func() {

		It("Allow everyone to create contract", func() {

			const MyUserName = `vitiko`
			// we get chaincode invoke response payload and expect
			contractCreateResponse := expectcc.PayloadIs(
				cc.Invoke(
					// invoke chaincode function
					`contract_create`,
					// цith ContractDTO payload, it will automatically will be marshalled to json
					&ContractDTO{
						UUID:             `xxx-aaa-bbb`,
						ContractTypeUUID: `xxx-ddd-ccc`,
						Username:         MyUserName,
						Password:         `Root123AsUsual`,
						FirstName:        `Victor`,
						LastName:         `Nosov`,
						Item: item{
							ID:          1,
							Brand:       `NoName`,
							Model:       `Model-XYZ`,
							Price:       123.45,
							Description: `Coolest thing ever`,
							SerialNo:    `ooo-999-222`,
						},
						StartDate: time.Now(),
						EndDate:   time.Now(),
					}),
				// We expect than payload response is marshalled ContractCreateResponse structure
				&ContractCreateResponse{}).(ContractCreateResponse)

			Expect(contractCreateResponse.Username).To(Equal(MyUserName))
		})

	})

})
