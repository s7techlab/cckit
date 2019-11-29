package testing

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/s7techlab/cckit/router"
)

type (
	TxHandler struct {
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

func NewTxHandler(name string) (*TxHandler, router.Context) {
	var (
		mockStub = NewMockStub(name, nil)
		logger   = router.NewLogger(name)
	)
	return &TxHandler{
			MockStub: mockStub,
			Logger:   logger,
		},
		router.NewContext(mockStub, logger)
}

func (p *TxHandler) Exec(
	txHdl func() (interface{}, error), middleware ...TxMiddleware) *TxResult {
	uuid := p.MockStub.generateTxUID()

	p.MockStub.MockTransactionStart(uuid)
	res, err := txHdl()
	p.MockStub.MockTransactionEnd(uuid)

	txRes := &TxResult{
		Result: res,
		Err:    err,
		Event:  p.MockStub.ChaincodeEvent,
	}
	for _, m := range middleware {
		m(txRes)
	}

	return txRes
}
