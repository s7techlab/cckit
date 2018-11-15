package expect

import (
	"fmt"

	"strconv"

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

// PayloadString expects payload content is string
func PayloadString(response peer.Response, expectedValue string) string {
	ResponseOk(response)
	str := string(response.Payload)
	g.Expect(str).To(g.Equal(expectedValue))
	return str
}

// PayloadBytes expects response is ok and compares response.Payload with expected value
func PayloadBytes(response peer.Response, expectedValue []byte) []byte {
	ResponseOk(response)
	g.Expect(response.Payload).To(g.Equal(expectedValue))
	return response.Payload
}

func PayloadInt(response peer.Response, expectedValue int) int {
	ResponseOk(response)
	d, err := strconv.Atoi(string((response.Payload)))
	g.Expect(err).To(g.BeNil())
	g.Expect(d).To(g.Equal(expectedValue))
	return d
}

// EventPayloadIs expects peer.ChaincodeEvent payload can be marshalled to target interface{} and returns converted value
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
