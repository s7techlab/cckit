package erc20_test

import (
	"encoding/base64"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/s7techlab/cckit/examples/token/chaincode/erc20"
	"github.com/s7techlab/cckit/examples/token/service/account"
	"github.com/s7techlab/cckit/examples/token/service/balance"
	"github.com/s7techlab/cckit/examples/token/service/config_erc20"
	"github.com/s7techlab/cckit/identity"
	"github.com/s7techlab/cckit/identity/testdata"
	testcc "github.com/s7techlab/cckit/testing"
	expectcc "github.com/s7techlab/cckit/testing/expect"
)

func TestERC20(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ERC20 Test suite")
}

var (
	ownerIdentity = testdata.Certificates[0].MustIdentity(testdata.DefaultMSP)
	user1Identity = testdata.Certificates[1].MustIdentity(testdata.DefaultMSP)
	user2Identity = testdata.Certificates[2].MustIdentity(testdata.DefaultMSP)

	ownerAddress = base64.StdEncoding.EncodeToString(identity.MarshalPublicKey(ownerIdentity.Cert.PublicKey))
	user1Address = base64.StdEncoding.EncodeToString(identity.MarshalPublicKey(user1Identity.Cert.PublicKey))
	user2Address = base64.StdEncoding.EncodeToString(identity.MarshalPublicKey(user2Identity.Cert.PublicKey))

	cc *testcc.MockStub
)

var _ = Describe(`ERC`, func() {

	BeforeSuite(func() {
		chaincode, err := erc20.New()
		Expect(err).NotTo(HaveOccurred())
		cc = testcc.NewMockStub(`erc20`, chaincode)

		expectcc.ResponseOk(cc.From(ownerIdentity).Init())
	})

	It(`Allows to call init once more time `, func() {
		expectcc.ResponseOk(cc.From(ownerIdentity).Init())
	})

	Context(`token info`, func() {

		It(`Allows to get token name`, func() {
			name := expectcc.PayloadIs(
				cc.From(user1Identity).
					Query(config_erc20.ConfigERC20ServiceChaincode_GetName, nil),
				&config_erc20.NameResponse{}).(*config_erc20.NameResponse)

			Expect(name.Name).To(Equal(erc20.Token.Name))
		})
	})

	Context(`initial balance`, func() {

		It(`Allows to know invoker address `, func() {
			address := expectcc.PayloadIs(
				cc.From(user1Identity).
					Query(account.AccountServiceChaincode_GetInvokerAddress, nil),
				&account.AddressId{}).(*account.AddressId)

			Expect(address.Address).To(Equal(user1Address))

			address = expectcc.PayloadIs(
				cc.From(user2Identity).
					Query(account.AccountServiceChaincode_GetInvokerAddress, nil),
				&account.AddressId{}).(*account.AddressId)

			Expect(address.Address).To(Equal(user2Address))
		})

		It(`Allows to get owner balance`, func() {
			b := expectcc.PayloadIs(
				cc.From(user1Identity). // call by any user
							Query(balance.BalanceServiceChaincode_GetBalance,
						&balance.BalanceId{Address: ownerAddress, Token: []string{erc20.Token.Name}}),
				&balance.Balance{}).(*balance.Balance)

			Expect(b.Address).To(Equal(ownerAddress))
			Expect(b.Amount).To(Equal(uint64(erc20.Token.TotalSupply)))
		})

		It(`Allows to get zero balance`, func() {
			b := expectcc.PayloadIs(
				cc.From(user1Identity).
					Query(balance.BalanceServiceChaincode_GetBalance,
						&balance.BalanceId{Address: user1Address, Token: []string{erc20.Token.Name}}),
				&balance.Balance{}).(*balance.Balance)

			Expect(b.Amount).To(Equal(uint64(0)))
		})

	})

	Context(`transfer`, func() {
		transferAmount := 100

		It(`Disallow to transfer balance  by user with zero balance`, func() {
			expectcc.ResponseError(
				cc.From(user1Identity).
					Invoke(balance.BalanceServiceChaincode_Transfer,
						&balance.TransferRequest{
							RecipientAddress: user2Address,
							Token:            []string{erc20.Token.Name},
							Amount:           100,
						}), balance.ErrAmountInsuficcient)

		})

		It(`Allows to transfer balance by owner`, func() {
			r := expectcc.PayloadIs(
				cc.From(ownerIdentity).
					Invoke(balance.BalanceServiceChaincode_Transfer,
						&balance.TransferRequest{
							RecipientAddress: user1Address,
							Token:            []string{erc20.Token.Name},
							Amount:           100,
						}),
				&balance.TransferResponse{}).(*balance.TransferResponse)

			Expect(r.SenderAddress).To(Equal(ownerAddress))
			Expect(r.Amount).To(Equal(uint64(transferAmount)))
		})

		It(`Allows to get new non zero balance`, func() {
			b := expectcc.PayloadIs(
				cc.From(user1Identity).
					Query(balance.BalanceServiceChaincode_GetBalance,
						&balance.BalanceId{Address: user1Address, Token: []string{erc20.Token.Name}}),
				&balance.Balance{}).(*balance.Balance)

			Expect(b.Amount).To(Equal(uint64(transferAmount)))
		})

	})
})
