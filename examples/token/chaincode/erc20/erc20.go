package erc20

import (
	"fmt"

	"github.com/s7techlab/cckit/examples/token/service/account"
	"github.com/s7techlab/cckit/examples/token/service/balance"
	"github.com/s7techlab/cckit/examples/token/service/config"
	"github.com/s7techlab/cckit/router"
)

var (
	Token = &config.CreateTokenTypeRequest{
		Name:        `SomeToken`,
		Symbol:      `@`,
		Decimals:    2,
		TotalSupply: 10000000,
	}
)

func New() (*router.Chaincode, error) {
	r := router.New(`erc20`)

	// accountSvc resolves address as base58( invoker.Cert.PublicKey )
	accountSvc := account.NewLocalService()
	configSvc := config.NewStateService()
	balanceSvc := balance.New(accountSvc, configSvc)

	r.Init(func(ctx router.Context) (interface{}, error) {
		// add token definition to state if not exists
		tokenType, _ := configSvc.GetTokenType(ctx, &config.TokenTypeId{Name: Token.Name})
		if tokenType != nil {
			return nil, nil
		}

		// init token on first Init call
		_, err := configSvc.CreateTokenType(ctx, Token)
		if err != nil {
			return nil, err
		}

		ownerAddress, err := accountSvc.GetInvokerAddress(ctx, nil)
		if err != nil {
			return nil, err
		}

		fmt.Println(`---`, ownerAddress.Address)

		// add  `TotalSupply` to identity
		if err = balance.NewStore(ctx).Add(ownerAddress.Address, []string{Token.Name}, Token.TotalSupply); err != nil {
			return nil, err
		}

		return nil, nil
	})

	if err := balance.RegisterBalanceServiceChaincode(r, balanceSvc); err != nil {
		return nil, err
	}
	if err := account.RegisterAccountServiceChaincode(r, accountSvc); err != nil {
		return nil, err
	}

	//if err := RegisterAllowanceServiceChaincode(r, allowance.New(balanceSvc)); err != nil {
	//	return nil, err
	//}

	return router.NewChaincode(r), nil
}
