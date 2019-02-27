package convert_test

import (
	"testing"

	"github.com/s7techlab/cckit/convert"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestState(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "State suite")
}

var _ = Describe(`Convert`, func() {

	It(`Bool`, func() {
		bTrue, err := convert.ToBytes(true)
		Expect(err).NotTo(HaveOccurred())
		Expect(bTrue).To(Equal([]byte(`true`)))

		bFalse, err := convert.ToBytes(false)
		Expect(err).NotTo(HaveOccurred())
		Expect(bFalse).To(Equal([]byte(`false`)))

		eTrue, err := convert.FromBytes(bTrue, convert.TypeBool)
		Expect(err).NotTo(HaveOccurred())
		Expect(eTrue.(bool)).To(Equal(true))

		eFalse, err := convert.FromBytes(bFalse, convert.TypeBool)
		Expect(err).NotTo(HaveOccurred())
		Expect(eFalse.(bool)).To(Equal(false))
	})

	It(`String`, func() {
		const MyStr = `my-string`
		bStr, err := convert.ToBytes(MyStr)
		Expect(err).NotTo(HaveOccurred())
		Expect(bStr).To(Equal([]byte(MyStr)))

		eStr, err := convert.FromBytes(bStr, convert.TypeString)
		Expect(err).NotTo(HaveOccurred())
		Expect(eStr.(string)).To(Equal(MyStr))
	})

	It(`Nil`, func() {
		bNil, err := convert.ToBytes(nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(bNil).To(Equal([]byte{}))
	})

})
