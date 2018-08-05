package main

import (
	"fmt"
	"testing"

	"ibm_app"

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
	cc := testcc.NewMockStub(`insurance`, new(ibm_app.SmartContract))

	Describe("Chaincode initialization ", func() {
		It("Allow to provide contract types attributes  during chaincode creation", func() {
			expectcc.ResponseOk(cc.Init(`init`, &ContractTypesDTO{ContractType1}))
		})
	})

	Describe("Contract type ", func() {

		It("Allow to retrieve all contract type, added during chaincode init", func() {
			contractTypes := expectcc.PayloadIs(
				cc.Invoke(`contract_type_ls`), &ContractTypesDTO{}).(ContractTypesDTO)

			Expect(len(contractTypes)).To(Equal(1))
			Expect(contractTypes[0].ShopType).To(Equal(ContractType1.ShopType))
		})
	})

	Describe("Contract", func() {

		It("Allow everyone to create contract", func() {

			// get chaincode invoke response payload and expect returned payload is serialized instance of some structure
			contractCreateResponse := expectcc.PayloadIs(
				cc.Invoke(
					// invoke chaincode function
					`contract_create`,
					// with ContractDTO payload, it will automatically will be marshalled to json
					&Contract1,
				),
				// We expect than payload response is marshalled ContractCreateResponse structure
				&ContractCreateResponse{}).(ContractCreateResponse)

			Expect(contractCreateResponse.Username).To(Equal(Contract1.Username))
		})

		It("ERROR: During contract creation we also created user - ERROR in smart contract", func() {

			fmt.Printf(`%+v`, cc.Invoke(`user_get_info`, &GetUserDTO{
				Username: Contract1.Username,
			}))
		})

	})

})
