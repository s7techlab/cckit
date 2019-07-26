package expect

import (
	g "github.com/onsi/gomega"
	"github.com/s7techlab/cckit/testing"
)

type (
	Stringer interface {
		String() string
	}

	TxResult struct {
		*testing.TxResult
	}
)

func SvcResponse(res *testing.TxResult) *TxResult {
	return &TxResult{TxResult: res}
}

func (r *TxResult) HasError(err error) *TxResult {
	if err == nil {
		g.Expect(r.Err).NotTo(g.HaveOccurred())
	} else {
		g.Expect(r.Err).To(g.Equal(err))
	}
	return r
}

func (r *TxResult) Is(expectedResult interface{}) *TxResult {
	r.HasError(nil)

	_, ok1 := r.Result.(Stringer)
	_, ok2 := expectedResult.(Stringer)
	if ok1 && ok2 {
		g.Expect(r.Result.(Stringer).String()).To(g.Equal(expectedResult.(Stringer).String()))
	} else {
		g.Expect(r.Result).To(g.BeEquivalentTo(expectedResult))
	}

	return r
}

func (r *TxResult) ProduceEvent(eventName string, eventPayload interface{}) {
	r.HasError(nil)
	EventIs(r.Event, eventName, eventPayload)
}
