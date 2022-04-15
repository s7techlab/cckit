package pinger_test

import (
	"github.com/golang/protobuf/ptypes/empty"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/s7techlab/cckit/extensions/pinger"
	"github.com/s7techlab/cckit/identity/testdata"
	testcc "github.com/s7techlab/cckit/testing"
)

var _ = Describe(`Chaincode pinger`, func() {

	var (
		Someone = testdata.Certificates[0].MustIdentity(`SOME_MSP`)

		pingerSvc = pinger.NewService()
		cc, ctx   = testcc.NewTxHandler(`Chaincode owner`)
	)

	It(`Ping`, func() {
		cc.From(Someone).Tx(func() {
			pingInfo, err := pingerSvc.Ping(ctx, &empty.Empty{})
			Expect(err).NotTo(HaveOccurred())

			Expect(pingInfo.InvokerId).To(Equal(Someone.GetID()))
			Expect(pingInfo.InvokerCert).To(Equal(Someone.GetPEM()))
			Expect(pingInfo.EndorsingServerTime).To(Not(BeNil()))
			Expect(pingInfo.TxTime).To(Not(BeNil()))
		})
	})
})
