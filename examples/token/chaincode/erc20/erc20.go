package erc20

import (
	"errors"

	"github.com/s7techlab/cckit/examples/token/service/account"
	"github.com/s7techlab/cckit/examples/token/service/balance"
	"github.com/s7techlab/cckit/examples/token/service/config"
	"github.com/s7techlab/cckit/examples/token/service/config_erc20"
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

	erc20ConfigSvc := &config_erc20.ERC20Service{Token: configSvc}

	r.Init(func(ctx router.Context) (interface{}, error) {
		// add token definition to state if not exists
		token, err := config.CreateDefaultToken(ctx, configSvc, Token)
		if err != nil {
			if errors.Is(err, config.ErrTokenAlreadyExists) {
				return nil, nil
			}
			return nil, err
		}

		ownerAddress, err := accountSvc.GetInvokerAddress(ctx, nil)
		if err != nil {
			return nil, err
		}

		// add  `TotalSupply` to chaincode first committer
		if err = balance.NewStore(ctx).Add(ownerAddress.Address, token, Token.TotalSupply); err != nil {
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
	if err := config_erc20.RegisterConfigERC20ServiceChaincode(r, erc20ConfigSvc); err != nil {
		return nil, err
	}
	//if err := RegisterAllowanceServiceChaincode(r, allowance.New(balanceSvc)); err != nil {
	//	return nil, err
	//}

	return router.NewChaincode(r), nil
}
