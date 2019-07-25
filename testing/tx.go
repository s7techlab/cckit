package testing

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/s7techlab/cckit/router"
)

type (
	CCService struct {
		MockStub *MockStub
		Logger   *shim.ChaincodeLogger
	}

	TxResult struct {
		Result interface{}
		Err    error
		Event  *peer.ChaincodeEvent
	}

	TxMiddleware func(*TxResult)
)

func NewCCService(name string) *CCService {
	return &CCService{
		MockStub: NewMockStub(name, nil),
		Logger:   router.NewLogger(name),
	}
}

func (p *CCService) Context() router.Context {
	return router.NewContext(p.MockStub, p.Logger)
}

func (p *CCService) Exec(
	txHdl func(ctx router.Context) (interface{}, error), middleware ...TxMiddleware) *TxResult {
	uuid := p.MockStub.generateTxUID()

	p.MockStub.MockTransactionStart(uuid)
	res, err := txHdl(p.Context())
	p.MockStub.MockTransactionEnd(uuid)

	txRes := &TxResult{Result: res, Err: err, Event: p.MockStub.ChaincodeEvent}
	for _, m := range middleware {
		m(txRes)
	}

	return txRes
}
