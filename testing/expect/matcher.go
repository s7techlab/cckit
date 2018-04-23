package expect

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	g "github.com/onsi/gomega"
	"github.com/s7techlab/cckit/convert"
)

// ExpectResponseOk expects peer.Response has shim.OK status and message has okSubstr prefix
func ResponseOk(response peer.Response, okSubstr ...string) {
	g.Expect(int(response.Status)).To(g.Equal(shim.OK), response.Message)

	if len(okSubstr) > 0 {
		g.Expect(response.Message).To(g.HavePrefix(okSubstr[0]), "ok message not match: "+response.Message)
	}

}

// ExpectResponseError expects peer.Response has shim.ERROR status and message has errorSubstr prefix
func ResponseError(response peer.Response, errorSubstr ...string) {
	g.Expect(int(response.Status)).To(g.Equal(shim.ERROR), response.Message)

	if len(errorSubstr) > 0 {
		g.Expect(response.Message).To(g.HavePrefix(errorSubstr[0]), "error message not match: "+response.Message)
	}
}

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
