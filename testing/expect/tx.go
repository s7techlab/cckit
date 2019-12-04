package expect

import (
	"fmt"

	g "github.com/onsi/gomega"
	"github.com/s7techlab/cckit/testing"
)

type (
	Stringer interface {
		String() string
	}

	txResult struct {
		*testing.TxResult
	}
)

func TxResult(res *testing.TxResult) *txResult {
	return &txResult{TxResult: res}
}

func (r *txResult) HasError(err interface{}) *txResult {
	if err == nil {
		g.Expect(r.Err).NotTo(g.HaveOccurred())
	} else {
		g.Expect(fmt.Sprintf(`%s`, r.Err)).To(g.HavePrefix(fmt.Sprintf(`%s`, err)))
		//g.Expect(errors.Is(r.Err, err))
	}
	return r
}

func (r *txResult) Is(expectedResult interface{}) *txResult {
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

func (r *txResult) ProduceEvent(eventName string, eventPayload interface{}) {
	r.HasError(nil)
	EventIs(r.Event, eventName, eventPayload)
}
