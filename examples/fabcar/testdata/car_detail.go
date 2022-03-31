package testdata

import (
	. "github.com/onsi/gomega"

	"github.com/s7techlab/cckit/examples/fabcar"
)

func ExpectDetailsViewContain(setDetails []*fabcar.SetCarDetail, getDetails []*fabcar.CarDetail) {
	Expect(setDetails).To(HaveLen(len(getDetails)))

	length := len(setDetails)
	for i := 0; i < length; i++ {
		for j := 0; j < length; j++ {
			if setDetails[i].Type == getDetails[j].Type {
				Expect(setDetails[i].Type).To(Equal(getDetails[j].Type))
				Expect(setDetails[i].Make).To(Equal(getDetails[j].Make))
			}
		}
	}
}
