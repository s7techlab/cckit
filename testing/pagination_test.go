package testing_test

import (
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	testcc "github.com/s7techlab/cckit/testing"
)

var (
	mockStub *testcc.MockStub
	iter     *testcc.MockStateRangeQueryPagedIterator
)

var _ = BeforeEach(func() {
	mockStub = testcc.NewMockStub("test", nil)
	mockStub.MockTransactionStart("init")
	mockStub.PutState("aa", []byte{10})
	mockStub.PutState("ab", []byte{11})
	mockStub.PutState("ac", []byte{12})
	mockStub.PutState("ad", []byte{13})
	mockStub.PutState("ae", []byte{14})
	mockStub.PutState("af", []byte{15})
	mockStub.PutState("ag", []byte{16})
	mockStub.PutState("ba", []byte{20})
	mockStub.PutState("bb", []byte{21})
	// workaround
	mockStub.TxResult = peer.Response{
		Status:  shim.OK,
		Message: "",
		Payload: nil,
	}

	mockStub.MockTransactionEnd("init")
})

var _ = Describe("MockStateRangePagedIterator", func() {
	Context("without bookmark", func() {
		It("should iterate over first 2 items", func() {
			iter = testcc.NewMockStatesRangeQueryPagedIterator(
				mockStub, "aa", "b", 2, "")
			Expect(iter.HasNext()).To(Equal(true))

			kv, err := iter.Next()
			Expect(err).NotTo(HaveOccurred())
			Expect(kv.Key).To(Equal("aa"))
			Expect(kv.Value).To(Equal([]byte{10}))
			Expect(iter.HasNext()).To(Equal(true))

			kv, err = iter.Next()
			Expect(err).NotTo(HaveOccurred())
			Expect(kv.Key).To(Equal("ab"))
			Expect(kv.Value).To(Equal([]byte{11}))
			Expect(iter.HasNext()).To(Equal(false))
		})
	})

	Context("with bookmark", func() {
		It("should iterate over 2 items from bookmark (inclusive)", func() {
			iter = testcc.NewMockStatesRangeQueryPagedIterator(
				mockStub, "aa", "b", 3, "ab")

			Expect(iter.HasNext()).To(Equal(true))
			kv, err := iter.Next()
			Expect(err).NotTo(HaveOccurred())
			Expect(kv.Key).To(Equal("ab"))
			Expect(kv.Value).To(Equal([]byte{11}))

			Expect(iter.HasNext()).To(Equal(true))
			kv, err = iter.Next()
			Expect(err).NotTo(HaveOccurred())
			Expect(kv.Key).To(Equal("ac"))
			Expect(kv.Value).To(Equal([]byte{12}))

			Expect(iter.HasNext()).To(Equal(true))
			kv, err = iter.Next()
			Expect(err).NotTo(HaveOccurred())
			Expect(kv.Key).To(Equal("ad"))
			Expect(kv.Value).To(Equal([]byte{13}))

			Expect(iter.HasNext()).To(Equal(false))
		})
	})

	Context("with bookmark less than startKey", func() {
		It("should iterate over 2 items from startKey", func() {
			iter = testcc.NewMockStatesRangeQueryPagedIterator(
				mockStub, "ac", "bb", 2, "ab")

			Expect(iter.HasNext()).To(Equal(true))
			kv, err := iter.Next()
			Expect(err).NotTo(HaveOccurred())
			Expect(kv.Key).To(Equal("ac"))
			Expect(kv.Value).To(Equal([]byte{12}))

			Expect(iter.HasNext()).To(Equal(true))
			kv, err = iter.Next()
			Expect(err).NotTo(HaveOccurred())
			Expect(kv.Key).To(Equal("ad"))
			Expect(kv.Value).To(Equal([]byte{13}))

			Expect(iter.HasNext()).To(Equal(false))
		})
	})

	Context("with bookmark greater than endKey", func() {
		It("shouldn't contains elements", func() {
			iter = testcc.NewMockStatesRangeQueryPagedIterator(
				mockStub, "ac", "ae", 2, "ba")

			Expect(iter.HasNext()).To(Equal(false))
		})
	})
})
