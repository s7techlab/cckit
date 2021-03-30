package testing

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/s7techlab/cckit/router"
	"github.com/s7techlab/cckit/testing/expect"
	"go.uber.org/zap"
)

type (
	TxHandler struct {
		MockStub *MockStub
		Logger   *zap.Logger
		Context  router.Context
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
		ctx      = router.NewContext(mockStub, logger)
	)
	return &TxHandler{
			MockStub: mockStub,
			Logger:   logger,
			Context:  ctx,
		},

		ctx
}

func (p *TxHandler) From(creator ...interface{}) *TxHandler {
	p.MockStub.From(creator...)
	return p
}

func (p *TxHandler) Init(txHdl func(ctx router.Context) (interface{}, error)) *TxResult {
	return p.Invoke(txHdl)
}

// Invoke emulates chaincode invocation and returns transaction result
func (p *TxHandler) Invoke(invoke func(ctx router.Context) (interface{}, error)) *TxResult {
	uuid := p.MockStub.generateTxUID()

	p.MockStub.MockTransactionStart(uuid)
	res, err := invoke(p.Context)
	p.MockStub.MockTransactionEnd(uuid)

	txRes := &TxResult{
		Result: res,
		Err:    err,
		Event:  p.MockStub.ChaincodeEvent,
	}

	return txRes
}

// Tx emulates chaincode invocation
func (p *TxHandler) Tx(tx func()) {
	uuid := p.MockStub.generateTxUID()

	p.MockStub.MockTransactionStart(uuid)
	tx()
	p.MockStub.MockTransactionEnd(uuid)
}

// TxFund returns tx closure
func (p *TxHandler) TxFunc(tx func()) func() {
	return func() {
		p.Tx(tx)
	}
}

// Expect returns assertion helper
func (p *TxHandler) Expect(res interface{}, err error) *expect.TxRes {
	return &expect.TxRes{
		Result: res,
		Err:    err,
		Event:  p.MockStub.ChaincodeEvent,
	}
}

// TxTimestamp returns last tx timestamp
func (p *TxHandler) TxTimestamp() *timestamp.Timestamp {
	return p.MockStub.TxTimestamp
}

// TxEvent returns last tx event
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
