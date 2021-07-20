package testing

import (
	"sync"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
	"go.uber.org/zap"

	"github.com/s7techlab/cckit/response"
	"github.com/s7techlab/cckit/router"
	"github.com/s7techlab/cckit/testing/expect"
)

type (
	TxHandler struct {
		MockStub *MockStub
		m        sync.Mutex
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
	}, ctx
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
	p.MockStub.TxResult = response.Create(res, err)
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
	p.m.Lock()
	defer p.m.Unlock()

	uuid := p.MockStub.generateTxUID()
	p.MockStub.MockTransactionStart(uuid)

	// expect that invoke will be with shim.OK status, need for dump state changes
	// if during tx func got error - func setTxResult must be called
	p.MockStub.TxResult = peer.Response{
		Status:  shim.OK,
		Message: "",
		Payload: nil,
	}
	tx()
	p.MockStub.MockTransactionEnd(uuid)
}

// SetTxResult can be used for set txResult error during Tx
func (p *TxHandler) SetTxResult(err error) {
	if p.MockStub.TxID == `` {
		panic(`can be called only during Tx() evaluation`)
	}
	if err != nil {
		p.MockStub.TxResult.Status = shim.ERROR
		p.MockStub.TxResult.Message = err.Error()
	}
}

// TxFunc returns tx closure, can be used directly as ginkgo func
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
