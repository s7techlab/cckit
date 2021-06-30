package state_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/s7techlab/cckit/state/testdata"
	testcc "github.com/s7techlab/cckit/testing"
	expectcc "github.com/s7techlab/cckit/testing/expect"
)

const (
	StateCachedChaincode = `state_cached`
)

var _ = Describe(`State caching`, func() {

	//Create chaincode mocks
	stateCachedCC := testcc.NewMockStub(StateCachedChaincode, testdata.NewStateCachedCC())

	It("Read after write returns non empty entry", func() {
		resp := expectcc.PayloadIs(stateCachedCC.Invoke(testdata.TxStateCachedReadAfterWrite), &testdata.Value{})
		Expect(resp).To(Equal(testdata.KeyValue(testdata.Keys[0])))
	})

	It("Read after delete returns empty entry", func() {
		resp := stateCachedCC.Invoke(testdata.TxStateCachedReadAfterDelete)
		Expect(resp.Payload).To(Equal([]byte{}))
	})

	It("List after write returns list", func() {
		resp := expectcc.PayloadIs(
			stateCachedCC.Invoke(testdata.TxStateCachedListAfterWrite), &[]testdata.Value{}).([]testdata.Value)

		// all key exists
		Expect(resp).To(Equal([]testdata.Value{
			testdata.KeyValue(testdata.Keys[0]), testdata.KeyValue(testdata.Keys[1]), testdata.KeyValue(testdata.Keys[2])}))
	})

	It("List after delete returns list without deleted item", func() {
		resp := expectcc.PayloadIs(
			stateCachedCC.Invoke(testdata.TxStateCachedListAfterDelete), &[]testdata.Value{}).([]testdata.Value)

		// first key is deleted
		Expect(resp).To(Equal([]testdata.Value{
			testdata.KeyValue(testdata.Keys[1]), testdata.KeyValue(testdata.Keys[2])}))
	})
})
