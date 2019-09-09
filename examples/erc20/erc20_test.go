package erc20_test

import (
	"testing"

	"github.com/s7techlab/cckit/examples/erc20"
	identitytestdata "github.com/s7techlab/cckit/identity/testdata"
	testcc "github.com/s7techlab/cckit/testing"
	expectcc "github.com/s7techlab/cckit/testing/expect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCars(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cars Suite")
}

var (
	// load actor certificates
	TokenOwner     = identitytestdata.Certificates[0].MustIdentity(`SOME_MSP`)
	AccountHolder1 = identitytestdata.Certificates[1].MustIdentity(`SOME_MSP`)
)
var _ = Describe(`ERC-20`, func() {

	const TokenSymbol = `HLF`
	const TokenName = `HLFCoin`
	const TotalSupply = 10000
	const Decimals = 3

	//Create chaincode mock
	erc20fs := testcc.NewMockStub(`erc20`, erc20.NewErc20FixedSupply())

	BeforeSuite(func() {
		// init token haincode
		expectcc.ResponseOk(erc20fs.From(TokenOwner).Init(TokenSymbol, TokenName, TotalSupply, Decimals))
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
					`balanceOf`, TokenOwner.MspID, TokenOwner.GetID()), TotalSupply)
		})

		It("Allow everyone to get holder's token balance", func() {
			expectcc.PayloadInt(
				erc20fs.Query(
					`balanceOf`, AccountHolder1.MspID, AccountHolder1.GetID()), 0)
		})
	})

	Describe("ERC-20 transfer", func() {

		It("Disallow to transfer token to same account", func() {
			expectcc.ResponseError(
				erc20fs.From(TokenOwner).Invoke(
					`transfer`, TokenOwner.MspID, TokenOwner.GetID(), 100),
				erc20.ErrForbiddenToTransferToSameAccount)
		})

		It("Disallow token holder with zero balance to transfer tokens", func() {
			expectcc.ResponseError(
				erc20fs.From(AccountHolder1).Invoke(
					`transfer`, TokenOwner.MspID, TokenOwner.GetID(), 100),
				erc20.ErrNotEnoughFunds)
		})

		It("Allow token holder with non zero balance to transfer tokens", func() {
			expectcc.PayloadInt(
				erc20fs.From(TokenOwner).Invoke(
					`transfer`, AccountHolder1.MspID, AccountHolder1.GetID(), 100),
				TotalSupply-100)

			expectcc.PayloadInt(
				erc20fs.Query(
					`balanceOf`, TokenOwner.MspID, TokenOwner.GetID()), TotalSupply-100)

			expectcc.PayloadInt(
				erc20fs.Query(
					`balanceOf`, AccountHolder1.MspID, AccountHolder1.GetID()), 100)
		})

	})

	Describe("ERC-20 transfer allowance", func() {

		It("Allow everyone to check token transfer allowance - zero initially", func() {
			expectcc.PayloadInt(
				erc20fs.Query(
					`allowance`,
					TokenOwner.MspID, TokenOwner.GetID(),
					AccountHolder1.MspID, AccountHolder1.GetID()), 0)
		})

		It("Allow token owner to set transfer allowance", func() {
			expectcc.ResponseOk(
				erc20fs.From(TokenOwner).Invoke(
					`approve`, AccountHolder1.MspID, AccountHolder1.GetID(), 10))
		})

		It("Allow everyone to check token transfer allowance", func() {
			expectcc.PayloadInt(
				erc20fs.Query(
					`allowance`,
					TokenOwner.MspID, TokenOwner.GetID(),
					AccountHolder1.MspID, AccountHolder1.GetID()), 10)
		})
	})
})
