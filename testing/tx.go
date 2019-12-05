package testing

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/s7techlab/cckit/router"
	"github.com/s7techlab/cckit/testing/expect"
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

func (p *TxHandler) Tx(tx func()) {
	uuid := p.MockStub.generateTxUID()

	p.MockStub.MockTransactionStart(uuid)
	tx()
	p.MockStub.MockTransactionEnd(uuid)
}

func (p *TxHandler) Invoke(invoke func() (interface{}, error)) *TxResult {
	uuid := p.MockStub.generateTxUID()

	p.MockStub.MockTransactionStart(uuid)
	res, err := invoke()
	p.MockStub.MockTransactionEnd(uuid)

	txRes := &TxResult{
		Result: res,
		Err:    err,
		Event:  p.MockStub.ChaincodeEvent,
	}

	return txRes
}

func (p *TxHandler) SvcExpect(res interface{}, err error) *expect.TxRes {
	return &expect.TxRes{
		Result: res,
		Err:    err,
		Event:  p.MockStub.ChaincodeEvent,
	}
}

// TxEvent returs last tx timestamp
func (p *TxHandler) TxTimestamp() *timestamp.Timestamp {
	return p.MockStub.TxTimestamp
}

// TxEvent returs last tx event
func (p *TxHandler) TxEvent() *peer.ChaincodeEvent {
	return p.MockStub.ChaincodeEvent
}

func (r *TxResult) Expect() *expect.TxRes {
	return &expect.TxRes{
		Result: r.Result,
		Err:    r.Err,
		Event:  r.Event,
	}
}
