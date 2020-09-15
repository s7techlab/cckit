package expect

import (
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	g "github.com/onsi/gomega"
	"github.com/s7techlab/cckit/convert"
)

// ResponseOk expects peer.Response has shim.OK status and message has okMatcher matcher
func ResponseOk(response peer.Response, okMatcher ...g.OmegaMatcher) peer.Response {
	g.Expect(int(response.Status)).To(g.Equal(shim.OK), response.Message)

	if len(okMatcher) > 0 {
		g.Expect(response.Message).To(okMatcher[0], "ok message not match: "+response.Message)
	}
	return response
}

// ResponseError expects peer.Response has shim.ERROR status and message has errMatcher matcher
func ResponseError(response peer.Response, errMatcher ...g.OmegaMatcher) peer.Response {
	g.Expect(int(response.Status)).To(g.Equal(shim.ERROR), response.Message)

	if len(errMatcher) > 0 {
		g.Expect(response.Message).To(errMatcher[0],
			"error message not match: "+response.Message)
	}

	return response
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
