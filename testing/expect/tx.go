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
		g.Expect(fmt.Sprintf(`%s`, r.Err)).To(g.HavePrefix(fmt.Sprintf(`%s`, err)))
		//g.Expect(errors.Is(r.Err, err))
	}
	return r
}

func (r *TxRes) Is(expectedResult interface{}) *TxRes {
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

func (r *TxRes) ProduceEvent(eventName string, eventPayload interface{}) {
	r.HasError(nil)
	EventIs(r.Event, eventName, eventPayload)
}
