package cflat

import (
	"github.com/s7techlab/cckit/extensions/owner"
	"github.com/s7techlab/cckit/router"
)

const ChaincodeName = `cflat`

func New() (*router.Chaincode, error) {

	r := router.New(ChaincodeName)

	if err := RegisterCFlatServiceChaincode(r, &CFlatService{}); err != nil {
		return nil, err
	}

	return router.NewChaincode(r), nil
}

func MustNew() *router.Chaincode {
	cc, err := New()
	if err != nil {
		panic(err)
	}

	return cc
}

func ChaincodeInitFunc(ownerSvc *owner.ChaincodeOwnerService) func(router.Context) (interface{}, error) {
	return func(ctx router.Context) (interface{}, error) {
		return ownerSvc.RegisterTxCreator(ctx)
	}
}
