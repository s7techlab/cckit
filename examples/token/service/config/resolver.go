package config

import (
	"github.com/s7techlab/cckit/router"
)

type TokenResolver interface {
	GetToken(ctx router.Context, id *TokenId) (*Token, error)
}
