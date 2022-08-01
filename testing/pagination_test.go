package testing_test

import (
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/ledger/queryresult"
	"github.com/hyperledger/fabric-protos-go/peer"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	testcc "github.com/s7techlab/cckit/testing"
)

var _ = Describe("MockStateRangePagedIterator", func() {
	var (
		mockStub *testcc.MockStub
		iter     *testcc.MockStateRangeQueryPagedIterator
		state    []*queryresult.KV
	)

	var _ = BeforeEach(func() {
		mockStub = testcc.NewMockStub("test", nil)
		state = []*queryresult.KV{
			{Key: "aa", Value: []byte{10}},
			{Key: "ab", Value: []byte{11}},
			{Key: "ac", Value: []byte{12}},
			{Key: "ad", Value: []byte{13}},
			{Key: "ae", Value: []byte{14}},
			{Key: "af", Value: []byte{15}},
			{Key: "ag", Value: []byte{16}},
			{Key: "ba", Value: []byte{20}},
			{Key: "bb", Value: []byte{21}},
		}

		if err := populateState(mockStub, state); err != nil {
			Fail(fmt.Sprintf("Couldn't populate state: %s", err.Error()))
		}
	})

	Context("without bookmark", func() {
		It("should iterate over first 2 items", func() {
			iter = testcc.NewMockStatesRangeQueryPagedIterator(
				mockStub, "aa", "b", 2, "")
			Expect(iter.Len()).To(Equal(int32(2)))
			Expect(iter.NextBookmark()).To(Equal("ac"))

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
			Expect(iter.Len()).To(Equal(int32(3)))
			Expect(iter.NextBookmark()).To(Equal("ae"))

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
			Expect(iter.Len()).To(Equal(int32(2)))
			Expect(iter.NextBookmark()).To(Equal("ae"))

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
			Expect(iter.Len()).To(Equal(int32(0)))
			Expect(iter.NextBookmark()).To(Equal(""))
		})
	})

	Context("with empty state", func() {
		It("shouldn't contains elements", func() {
			emptyStub := testcc.NewMockStub("test", nil)
			iter = testcc.NewMockStatesRangeQueryPagedIterator(emptyStub, "", "", 10, "")

			Expect(iter.Len()).To(Equal(int32(0)))
			Expect(iter.HasNext()).To(Equal(false))
			Expect(iter.NextBookmark()).To(Equal(""))
		})
	})

	Context("with unbound range", func() {
		It("should contains upto pageSize elements", func() {
			iter = testcc.NewMockStatesRangeQueryPagedIterator(mockStub, "", "", 6, "")

			Expect(iter.Len()).To(Equal(int32(6)))
			Expect(iter.HasNext()).To(Equal(true))
			Expect(iter.NextBookmark()).To(Equal("ag"))
		})
	})

	Context("when iterate over last elements", func() {
		It("shouldn't has next page", func() {
			iter = testcc.NewMockStatesRangeQueryPagedIterator(mockStub, "", "", 6, "ae")
			Expect(iter.Len()).To(Equal(int32(5)))
			Expect(iter.HasNext()).To(Equal(true))
			Expect(iter.NextBookmark()).To(Equal(""))
		})
	})
})

var _ = Describe("MockStub", func() {
	Describe("GetStateByRangeWithPagination", func() {
		var (
			mockStub *testcc.MockStub
			state    []*queryresult.KV
		)

		var _ = BeforeEach(func() {
			mockStub = testcc.NewMockStub("test", nil)
			state = []*queryresult.KV{
				{Key: "aa", Value: []byte{10}},
				{Key: "ab", Value: []byte{11}},
				{Key: "ac", Value: []byte{12}},
				{Key: "ad", Value: []byte{13}},
				{Key: "ae", Value: []byte{14}},
				{Key: "af", Value: []byte{15}},
				{Key: "ag", Value: []byte{16}},
				{Key: "ba", Value: []byte{20}},
				{Key: "bb", Value: []byte{21}},
			}
			if err := populateState(mockStub, state); err != nil {
				Fail(fmt.Sprintf("Couldn't populate state: %s", err.Error()))
			}
		})

		It("should return first 5 elements in range", func() {
			iter, md, err := mockStub.GetStateByRangeWithPagination("aa", "ba", 5, "")

			Expect(err).NotTo(HaveOccurred())
			Expect(md.Bookmark).To(Equal("af"))
			Expect(md.FetchedRecordsCount).To(Equal(int32(5)))

			for _, expect := range state[0:5] {
				kv, err := iter.Next()
				Expect(err).NotTo(HaveOccurred())
				Expect(kv).To(Equal(expect))
			}

			Expect(iter.HasNext()).To(Equal(false))
		})

		It("should return 3 elements after bookmark(inclusive)", func() {
			iter, md, err := mockStub.GetStateByRangeWithPagination("aa", "ba", 3, "ad")

			Expect(err).NotTo(HaveOccurred())
			Expect(md.Bookmark).To(Equal("ag"))
			Expect(md.FetchedRecordsCount).To(Equal(int32(3)))

			for _, expect := range state[3:6] {
				kv, err := iter.Next()
				Expect(err).NotTo(HaveOccurred())
				Expect(kv).To(Equal(expect))
			}

			Expect(iter.HasNext()).To(Equal(false))
		})

		It("should return last 2 elements", func() {
			iter, md, err := mockStub.GetStateByRangeWithPagination("aa", "ba", 3, "ag")

			Expect(err).NotTo(HaveOccurred())
			Expect(md.Bookmark).To(Equal(""))
			Expect(md.FetchedRecordsCount).To(Equal(int32(1)))

			for _, expect := range state[6:7] {
				kv, err := iter.Next()
				Expect(err).NotTo(HaveOccurred())
				Expect(kv).To(Equal(expect))
			}

			Expect(iter.HasNext()).To(Equal(false))
		})

		It("should be empty when bookmark equal to endKey", func() {
			iter, md, err := mockStub.GetStateByRangeWithPagination("aa", "ba", 3, "ba")

			Expect(err).NotTo(HaveOccurred())
			Expect(md.Bookmark).To(Equal(""))
			Expect(md.FetchedRecordsCount).To(Equal(int32(0)))

			Expect(iter.HasNext()).To(Equal(false))
		})

		It("should return items in range when bookmark equal to startKey", func() {
			iter, md, err := mockStub.GetStateByRangeWithPagination("aa", "ba", 3, "aa")

			Expect(err).NotTo(HaveOccurred())
			Expect(md.Bookmark).To(Equal("ad"))
			Expect(md.FetchedRecordsCount).To(Equal(int32(3)))

			for _, expect := range state[0:3] {
				kv, err := iter.Next()
				Expect(err).NotTo(HaveOccurred())
				Expect(kv).To(Equal(expect))
			}

			Expect(iter.HasNext()).To(Equal(false))
		})

		It("should returns items in range when bookmark is less than startKey", func() {
			iter, md, err := mockStub.GetStateByRangeWithPagination("af", "ba", 3, "ab")

			Expect(err).NotTo(HaveOccurred())
			Expect(md.Bookmark).To(Equal(""))
			Expect(md.FetchedRecordsCount).To(Equal(int32(2)))

			for _, expect := range state[5:7] {
				kv, err := iter.Next()
				Expect(err).NotTo(HaveOccurred())
				Expect(kv).To(Equal(expect))
			}

			Expect(iter.HasNext()).To(Equal(false))
		})
	})

	Describe("GetStateByPartialCompositeKeyWithPagination", func() {
		var (
			mockStub *testcc.MockStub
			state    []*queryresult.KV
		)

		var _ = BeforeEach(func() {
			mockStub = testcc.NewMockStub("test", nil)
			k := keyComposer(mockStub)
			state = []*queryresult.KV{
				{Key: k("test/Foo", "a"), Value: []byte{10}},
				{Key: k("test/Foo", "b"), Value: []byte{11}},
				{Key: k("test/Foo", "c"), Value: []byte{12}},
				{Key: k("test/Foo", "d"), Value: []byte{13}},
				{Key: k("test/Foo", "e"), Value: []byte{14}},
				{Key: k("test/Foo", "f"), Value: []byte{15}},
				{Key: k("test/Zoo", "a"), Value: []byte{20}},
				{Key: k("test/Zoo", "b"), Value: []byte{21}},
				{Key: k("test/Zoo", "c"), Value: []byte{22}},
			}
			if err := populateState(mockStub, state); err != nil {
				Fail(fmt.Sprintf("Couldn't populate state: %s", err.Error()))
			}
		})

		It("should returns single element", func() {
			iter, md, err := mockStub.GetStateByPartialCompositeKeyWithPagination(
				"test/Foo", []string{"b"}, 10, "")
			Expect(err).NotTo(HaveOccurred())
			Expect(md.FetchedRecordsCount).To(Equal(int32(1)))
			Expect(md.Bookmark).To(Equal(""))
			Expect(iter.HasNext()).To(Equal(true))

			kv, err := iter.Next()
			Expect(err).NotTo(HaveOccurred())
			Expect(kv).To(Equal(state[1]))

			Expect(iter.HasNext()).To(Equal(false))
		})

		It("should returns bookmark to next page when matched elements more than page size", func() {
			iter, md, err := mockStub.GetStateByPartialCompositeKeyWithPagination("test/Foo", nil, 3, "")
			Expect(err).NotTo(HaveOccurred())
			Expect(md.Bookmark).To(Equal(state[3].Key))
			Expect(md.FetchedRecordsCount).To(Equal(int32(3)))

			for _, expect := range state[0:3] {
				Expect(iter.HasNext()).To(Equal(true))
				kv, err := iter.Next()
				Expect(err).NotTo(HaveOccurred())
				Expect(kv).To(Equal(expect))
			}

			Expect(iter.HasNext()).To(Equal(false))
		})

		It("should return elements after bookmark", func() {
			iter, md, err := mockStub.GetStateByPartialCompositeKeyWithPagination("test/Foo", nil, 3, state[4].Key)
			Expect(err).NotTo(HaveOccurred())
			Expect(md.Bookmark).To(Equal(""))
			Expect(md.FetchedRecordsCount).To(Equal(int32(2)))

			for _, expect := range state[4:6] {
				Expect(iter.HasNext()).To(Equal(true))
				kv, err := iter.Next()
				Expect(err).NotTo(HaveOccurred())
				Expect(kv).To(Equal(expect))
			}

			Expect(iter.HasNext()).To(Equal(false))
		})

		It("should returns empty result for incorrect bookmark", func() {
			iter, md, err := mockStub.GetStateByPartialCompositeKeyWithPagination("test/Foo", nil, 3, state[7].Key)

			Expect(err).NotTo(HaveOccurred())
			Expect(md.Bookmark).To(Equal(""))
			Expect(md.FetchedRecordsCount).To(Equal(int32(0)))
			Expect(iter.HasNext()).To(Equal(false))
		})
	})
})

// keyComposer returns wrapper upon mockStub.CreateCompositeKey function
func keyComposer(mockStub *testcc.MockStub) func(objectType string, attrs ...string) string {
	return func(objectType string, attrs ...string) string {
		key, err := mockStub.CreateCompositeKey(objectType, attrs)
		if err != nil {
			Fail(fmt.Sprintf("Couldn't compose key: %s", err.Error()))
		}
		return key
	}
}

// populateState populate mock stub state with given key - value pairs
func populateState(mockStub *testcc.MockStub, values []*queryresult.KV) error {
	mockStub.MockTransactionStart("init")
	for _, kv := range values {
		if err := mockStub.PutState(kv.Key, kv.Value); err != nil {
			return err
		}
	}
	// workaround
	mockStub.TxResult = peer.Response{
		Status:  shim.OK,
		Message: "",
		Payload: nil,
	}
	mockStub.MockTransactionEnd("init")

	return nil
}
