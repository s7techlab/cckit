package testing

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	. "github.com/onsi/gomega"
)

func ExpectResponseOk(response peer.Response, okSubstr ...string) {
	Expect(int(response.Status)).To(Equal(shim.OK), response.Message)

	if len(okSubstr) > 0 {
		Expect(response.Message).To(HavePrefix(okSubstr[0]), "ok message not match: "+response.Message)
	}

}

// ExpectResponseError  expects peer.Response.Status is shim.ERROR
func ExpectResponseError(response peer.Response, errorSubstr ...string) {
	Expect(int(response.Status)).To(Equal(shim.ERROR), response.Message)

	if len(errorSubstr) > 0 {
		Expect(response.Message).To(HavePrefix(errorSubstr[0]), "error message not match: "+response.Message)
	}
}
