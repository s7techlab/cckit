package config_erc20

import (
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/s7techlab/cckit/examples/token/service/config"
	"github.com/s7techlab/cckit/router"
)

type ERC20Service struct {
	Token config.TokenGetter
}

func (s *ERC20Service) GetName(ctx router.Context, e *emptypb.Empty) (*NameResponse, error) {
	token, err := s.Token.GetDefaultToken(ctx, e)
	if err != nil {
		return nil, err
	}

	return &NameResponse{Name: token.GetType().GetName()}, nil
}

func (s *ERC20Service) GetSymbol(ctx router.Context, e *emptypb.Empty) (*SymbolResponse, error) {
	token, err := s.Token.GetDefaultToken(ctx, e)
	if err != nil {
		return nil, err
	}

	return &SymbolResponse{Symbol: token.GetType().GetSymbol()}, nil
}

func (s *ERC20Service) GetDecimals(ctx router.Context, e *emptypb.Empty) (*DecimalsResponse, error) {
	token, err := s.Token.GetDefaultToken(ctx, e)
	if err != nil {
		return nil, err
	}

	return &DecimalsResponse{Decimals: token.GetType().GetDecimals()}, nil
}

func (s *ERC20Service) GetTotalSupply(ctx router.Context, e *emptypb.Empty) (*TotalSupplyResponse, error) {
	token, err := s.Token.GetDefaultToken(ctx, e)
	if err != nil {
		return nil, err
	}

	return &TotalSupplyResponse{TotalSupply: token.GetType().GetTotalSupply()}, nil
}
