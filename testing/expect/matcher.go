package expect

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	g "github.com/onsi/gomega"
	"github.com/s7techlab/cckit/convert"
)

// ResponseOk expects peer.Response has shim.OK status and message has okSubstr prefix
func ResponseOk(response peer.Response, okSubstr ...string) {
	g.Expect(int(response.Status)).To(g.Equal(shim.OK), response.Message)

	if len(okSubstr) > 0 {
		g.Expect(response.Message).To(g.HavePrefix(okSubstr[0]), "ok message not match: "+response.Message)
	}

}

// ResponseError expects peer.Response has shim.ERROR status and message has errorSubstr prefix
func ResponseError(response peer.Response, errorSubstr ...interface{}) {
	g.Expect(int(response.Status)).To(g.Equal(shim.ERROR), response.Message)

	if len(errorSubstr) > 0 {
		g.Expect(response.Message).To(g.HavePrefix(fmt.Sprintf(`%s`, errorSubstr[0])),
			"error message not match: "+response.Message)
	}
}

// PayloadIs expects peer.Response payload can be marshalled to target interface{} and returns converted value
func PayloadIs(response peer.Response, target interface{}) interface{} {
	ResponseOk(response)
	data, err := convert.FromBytes(response.Payload, target)
	description := ``
	if err != nil {
		description = err.Error()
	}
	g.Expect(err).To(g.BeNil(), description)
	return data
}

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
