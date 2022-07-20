package config

import (
	"errors"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/s7techlab/cckit/router"
)

var (
	ErrTokenAlreadyExists = errors.New(`token already exists`)
)

type TokenGetter interface {
	GetToken(router.Context, *TokenId) (*Token, error)
	GetDefaultToken(router.Context, *emptypb.Empty) (*Token, error)
}

func TokenByTypeGroup(tokenType *TokenType, tokenGroup *TokenGroup) []string {
	return append([]string{tokenType.GetName()}, tokenGroup.GetName()...)
}

func CreateDefaultToken(
	ctx router.Context, configSvc ConfigServiceChaincode, createToken *CreateTokenTypeRequest) ([]string, error) {

	existsTokenType, _ := configSvc.GetTokenType(ctx, &TokenTypeId{Name: createToken.Name})
	if existsTokenType != nil {
		return nil, ErrTokenAlreadyExists
	}

	// init token on first Init call
	tokenType, err := configSvc.CreateTokenType(ctx, createToken)
	if err != nil {
		return nil, err
	}

	token := TokenByTypeGroup(tokenType, nil)

	if _, err = configSvc.SetConfig(ctx, &Config{DefaultToken: token}); err != nil {
		return nil, err
	}

	return token, nil
}
