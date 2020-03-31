package expect

import (
	"fmt"

	"github.com/hyperledger/fabric/protos/peer"
	g "github.com/onsi/gomega"
)

type (
	Stringer interface {
		String() string
	}

	TxRes struct {
		Result interface{}
		Err    error
		Event  *peer.ChaincodeEvent
	}
)

func (r *TxRes) HasError(err interface{}) *TxRes {
	if err == nil {
		g.Expect(r.Err).NotTo(g.HaveOccurred())
	} else {
		g.Expect(fmt.Sprintf(`%s`, r.Err)).To(g.ContainSubstring(fmt.Sprintf(`%s`, err)))
	}
	return r
}

func (r *TxRes) HasNoError() *TxRes {
	return r.HasError(nil)
}

func (r *TxRes) Is(expectedResult interface{}) *TxRes {
	r.HasNoError()

	_, ok1 := r.Result.(Stringer)
	_, ok2 := expectedResult.(Stringer)
	if ok1 && ok2 {
		g.Expect(r.Result.(Stringer).String()).To(g.Equal(expectedResult.(Stringer).String()))
	} else {
		g.Expect(r.Result).To(g.BeEquivalentTo(expectedResult))
	}

	return r
}

// ProduceEvent expects that tx produces event with particular payload
func (r *TxRes) ProduceEvent(eventName string, eventPayload interface{}) *TxRes {
	r.HasNoError()
	EventIs(r.Event, eventName, eventPayload)
	return r
}
