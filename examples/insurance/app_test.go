package insurance

import (
	"testing"

	"github.com/s7techlab/cckit/examples/insurance/app"

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
	cc := testcc.NewMockStub(`insurance`, new(app.SmartContract))

	Describe("Chaincode initialization ", func() {
		It("Allow to provide contract types attributes  during chaincode creation [init]", func() {
			expectcc.ResponseOk(cc.Init(`init`, &ContractTypesDTO{ContractType1}))
		})
	})

	Describe("Contract type ", func() {

		It("Allow to retrieve all contract type, added during chaincode init [contract_type_ls]", func() {
			contractTypes := expectcc.PayloadIs(
				cc.Invoke(`contract_type_ls`), &ContractTypesDTO{}).(ContractTypesDTO)

			Expect(len(contractTypes)).To(Equal(1))
			Expect(contractTypes[0].ShopType).To(Equal(ContractType1.ShopType))
		})

		It("Allow to create new contract type [contract_type_create]", func() {
			expectcc.ResponseOk(cc.Invoke(`contract_type_create`, &ContractType2))

			// get contract type list
			contractTypes := expectcc.PayloadIs(
				cc.Invoke(`contract_type_ls`), &ContractTypesDTO{}).(ContractTypesDTO)

			//expect now we have 2 contract type
			Expect(len(contractTypes)).To(Equal(2))
		})

		It("Allow to set active contract type [contract_type_set_active]", func() {
			Expect(ContractType2.Active).To(BeFalse())

			// active ContractType2
			expectcc.ResponseOk(cc.Invoke(`contract_type_set_active`, &ContractTypeActiveDTO{
				UUID: ContractType2.UUID, Active: true}))
		})

		It("Allow to retrieve filtered by shop type contract types [contract_type_ls]", func() {
			contractTypes := expectcc.PayloadIs(
				cc.Invoke(`contract_type_ls`, &ShopTypeDTO{ContractType2.ShopType}),
				&ContractTypesDTO{}).(ContractTypesDTO)

			Expect(len(contractTypes)).To(Equal(1))
			Expect(contractTypes[0].UUID).To(Equal(ContractType2.UUID))

			// Contract type 2 activated on previous step
			Expect(contractTypes[0].Active).To(BeTrue())
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

		It("Allow every one to get user info", func() {

			// orininally was error https://github.com/IBM/build-blockchain-insurance-app/pull/44
			user := expectcc.PayloadIs(
				cc.Invoke(`user_get_info`, &GetUserDTO{
					Username: Contract1.Username,
				}), &ResponseUserDTO{}).(ResponseUserDTO)

			Expect(user.LastName).To(Equal(Contract1.LastName))
		})

	})

})
