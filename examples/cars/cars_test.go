package main

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	testcc "github.com/s7techlab/cckit/testing"
	expectcc "github.com/s7techlab/cckit/testing/expect"
)

func TestCars(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cars Suite")
}

var _ = Describe(`Cars`, func() {

	//Create chaincode mock
	cc := testcc.NewMockStub(`cars`, New())
	car1 := &Car{
		Id:    `A777MP77`,
		Title: `BMW`,
		Owner: `victor-nosov`,
	}

	car2 := &Car{
		Id:    `O888OO77`,
		Title: `TOYOTA`,
		Owner: `alexander`,
	}

	BeforeSuite(func() {
		expectcc.ResponseOk(cc.Init())
	})

	Describe("Car", func() {

		It("Allow everyone to add information about car", func() {
			//check that invoke method
			expectcc.ResponseOk(cc.Invoke(`carRegister`, car1))
		})

		It("Disallow everyone to add duplicate information about car", func() {
			expectcc.ResponseError(
				cc.Invoke(`carRegister`, car1), ErrCarAlreadyExists) //expect  this error
		})

		It("Allow everyone to retrieve car information", func() {
			car := expectcc.PayloadIs(cc.Invoke(`carGet`, car1.Id),
				&Car{}).(Car)

			Expect(car.Title).To(Equal(car1.Title))
			Expect(car.Id).To(Equal(car1.Id))
		})

		It("Allow everyone to get car list", func() {
			//  &[]Car{} - declares target type for unmarshalling from []byte received from chaincode
			cars := expectcc.PayloadIs(cc.Invoke(`carList`), &[]Car{}).([]Car)

			Expect(len(cars)).To(Equal(1))
			Expect(cars[0].Id).To(Equal(car1.Id))
		})

		It("Allow everyone to add more information about car", func() {
			//check that invoke method
			expectcc.ResponseOk(cc.Invoke(`carRegister`, car2))
			cars := expectcc.PayloadIs(cc.Invoke(`carList`), &[]Car{}).([]Car)

			Expect(len(cars)).To(Equal(2))
		})
	})
})
