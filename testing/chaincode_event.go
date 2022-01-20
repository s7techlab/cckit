package testing

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/hyperledger/fabric-protos-go/peer"
)

type (
	ChaincodeEvent struct {
		event       *peer.ChaincodeEvent
		block       uint64
		txTimestamp *timestamp.Timestamp
	}
)

func (eb *ChaincodeEvent) Event() *peer.ChaincodeEvent {
	return eb.event
}

func (eb *ChaincodeEvent) Block() uint64 {
	return eb.block
}

func (eb *ChaincodeEvent) TxTimestamp() *timestamp.Timestamp {
	return eb.txTimestamp
}
