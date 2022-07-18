package config_erc20

import (
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/s7techlab/cckit/router"
)

type Service struct {
}

func (s *Service) GetName(ctx router.Context, _ *emptypb.Empty) (*NameResponse, error) {
	config, err := s.GetConfig(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &NameResponse{Name: config.Name}, nil
}

func (s *Service) GetSymbol(ctx router.Context, _ *emptypb.Empty) (*SymbolResponse, error) {
	config, err := s.GetConfig(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &SymbolResponse{Symbol: config.Symbol}, nil
}

func (s *Service) GetDecimals(ctx router.Context, _ *emptypb.Empty) (*DecimalsResponse, error) {
	config, err := s.GetConfig(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &DecimalsResponse{Decimals: config.Decimals}, nil
}

func (s *Service) GetTotalSupply(ctx router.Context, _ *emptypb.Empty) (*TotalSupplyResponse, error) {
	config, err := s.GetConfig(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &TotalSupplyResponse{TotalSupply: config.TotalSupply}, nil
}
