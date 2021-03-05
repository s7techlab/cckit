package expect

import (
	"github.com/hyperledger/fabric-protos-go/peer"
	g "github.com/onsi/gomega"
	"github.com/s7techlab/cckit/convert"
	"github.com/s7techlab/cckit/testing/gomega"
)

// EventIs expects ChaincodeEvent name is equal to expectName and event payload can be marshaled to expectPayload
func EventIs(event *peer.ChaincodeEvent, expectName string, expectPayload interface{}) interface{} {
	g.Expect(event.EventName).To(g.Equal(expectName), `event name not match`)

	return EventPayloadIs(event, expectPayload)
}

// EventStringerEqual expects ChaincodeEvent name is equal to expectName and
// event payload String() equal expectPayload String()
func EventStringerEqual(event *peer.ChaincodeEvent, expectName string, expectPayload interface{}) {
	payload := EventIs(event, expectName, expectPayload)

	g.Expect(payload).To(gomega.StringerEqual(expectPayload))
}

// EventPayloadIs expects peer.ChaincodeEvent payload can be marshaled to
// target interface{} and returns converted value
func EventPayloadIs(event *peer.ChaincodeEvent, target interface{}) interface{} {
	g.Expect(event).NotTo(g.BeNil())
	data, err := convert.FromBytes(event.Payload, target)
	description := ``
	if err != nil {
		description = err.Error()
	}
	g.Expect(err).To(g.BeNil(), description)
	return data
}
