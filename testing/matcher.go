package testing

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	g "github.com/onsi/gomega"
)

// ExpectResponseOk expects peer.Response has shim.OK status and message has okSubstr prefix
func ExpectResponseOk(response peer.Response, okSubstr ...string) {
	g.Expect(int(response.Status)).To(g.Equal(shim.OK), response.Message)

	if len(okSubstr) > 0 {
		g.Expect(response.Message).To(g.HavePrefix(okSubstr[0]), "ok message not match: "+response.Message)
	}

}

// ExpectResponseError expects peer.Response has shim.ERROR status and message has errorSubstr prefix
func ExpectResponseError(response peer.Response, errorSubstr ...string) {
	g.Expect(int(response.Status)).To(g.Equal(shim.ERROR), response.Message)

	if len(errorSubstr) > 0 {
		g.Expect(response.Message).To(g.HavePrefix(errorSubstr[0]), "error message not match: "+response.Message)
	}
}
