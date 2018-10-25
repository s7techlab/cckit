package erc20

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	examplecert "github.com/s7techlab/cckit/examples/cert"
	"github.com/s7techlab/cckit/identity"
	testcc "github.com/s7techlab/cckit/testing"
	expectcc "github.com/s7techlab/cckit/testing/expect"
)

func TestCars(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cars Suite")
}

var _ = Describe(`ERC-20`, func() {

	const TokenSymbol = `HLF`
	const TokenName = `HLFCoin`
	const TotalSupply = 10000
	const Decimals = 3

	//Create chaincode mock
	erc20fs := testcc.NewMockStub(`erc20`, NewErc20FixedSupply())

	// load actor certificates
	actors, err := identity.ActorsFromPemFile(`SOME_MSP`, map[string]string{
		`token_owner`:     `s7techlab.pem`,
		`account_holder1`: `victor-nosov.pem`,
		//`accoubt_holder2`: `victor-nosov.pem`
	}, examplecert.Content)
	if err != nil {
		panic(err)
	}

	BeforeSuite(func() {
		// init token haincode
		expectcc.ResponseOk(erc20fs.From(actors[`token_owner`]).Init(TokenSymbol, TokenName, TotalSupply, Decimals))
	})

	Describe("ERC-20 creation", func() {

		It("Allow everyone to get token symbol", func() {
			expectcc.PayloadString(erc20fs.Query(`symbol`), TokenSymbol)
		})

		It("Allow everyone to get token name", func() {
			expectcc.PayloadString(erc20fs.Query(`name`), TokenName)
		})

		It("Allow everyone to get token total supply", func() {
			expectcc.PayloadInt(erc20fs.Query(`totalSupply`), TotalSupply)
		})

		It("Allow everyone to get owner's token balance", func() {
			expectcc.PayloadInt(
				erc20fs.Query(
					`balanceOf`, actors[`token_owner`].GetMSPID(), actors[`token_owner`].GetID()), TotalSupply)
		})

		It("Allow everyone to get holder's token balance", func() {
			expectcc.PayloadInt(
				erc20fs.Query(
					`balanceOf`, actors[`account_holder1`].GetMSPID(), actors[`account_holder1`].GetID()), 0)
		})
	})

	Describe("ERC-20 transfer", func() {

		It("Disallow to transfer token to same account", func() {
			expectcc.ResponseError(
				erc20fs.From(actors[`token_owner`]).Invoke(
					`transfer`, actors[`token_owner`].GetMSPID(), actors[`token_owner`].GetID(), 100),
				ErrForbiddenToTransferToSameAccount)
		})

		It("Disallow token holder with zero balance to transfer tokens", func() {
			expectcc.ResponseError(
				erc20fs.From(actors[`account_holder1`]).Invoke(
					`transfer`, actors[`token_owner`].GetMSPID(), actors[`token_owner`].GetID(), 100),
				ErrNotEnoughFunds)
		})

		It("Allow token holder with non zero balance to transfer tokens", func() {
			expectcc.PayloadInt(
				erc20fs.From(actors[`token_owner`]).Invoke(
					`transfer`, actors[`account_holder1`].GetMSPID(), actors[`account_holder1`].GetID(), 100),
				TotalSupply-100)

			expectcc.PayloadInt(
				erc20fs.Query(
					`balanceOf`, actors[`token_owner`].GetMSPID(), actors[`token_owner`].GetID()), TotalSupply-100)

			expectcc.PayloadInt(
				erc20fs.Query(
					`balanceOf`, actors[`account_holder1`].GetMSPID(), actors[`account_holder1`].GetID()), 100)
		})

	})

})
