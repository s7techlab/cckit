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

func (p *TxHandler) From(creator ...interface{}) *TxHandler {
	p.MockStub.From(creator...)
	return p
}

func (p *TxHandler) Init(txHdl func() (interface{}, error)) *TxResult {
	return p.Invoke(txHdl)
}

func (p *TxHandler) Invoke(txHdl func() (interface{}, error)) *TxResult {
	uuid := p.MockStub.generateTxUID()

	p.MockStub.MockTransactionStart(uuid)
	res, err := txHdl()
	p.MockStub.MockTransactionEnd(uuid)

	txRes := &TxResult{
		Result: res,
		Err:    err,
		Event:  p.MockStub.ChaincodeEvent,
	}

	return txRes
}
