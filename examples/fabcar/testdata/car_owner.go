package testdata

import (
	. "github.com/onsi/gomega"

	"github.com/s7techlab/cckit/examples/fabcar"
)

func ExpectOwnersViewContain(setOwners []*fabcar.SetCarOwner, getOwners []*fabcar.CarOwner) {
	Expect(setOwners).To(HaveLen(len(getOwners)))

	length := len(setOwners)
	for i := 0; i < length; i++ {
		for j := 0; j < length; j++ {
			if setOwners[i].FirstName == getOwners[j].FirstName && setOwners[i].SecondName == getOwners[j].SecondName {
				Expect(setOwners[i].FirstName).To(Equal(getOwners[j].FirstName))
				Expect(setOwners[i].SecondName).To(Equal(getOwners[j].SecondName))
				Expect(setOwners[i].VehiclePassport).To(Equal(getOwners[j].VehiclePassport))
			}
		}
	}
}
