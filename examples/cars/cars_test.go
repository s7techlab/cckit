package cars_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/s7techlab/cckit/examples/cars"
	"github.com/s7techlab/cckit/extensions/owner"
	"github.com/s7techlab/cckit/identity/testdata"
	"github.com/s7techlab/cckit/state"
	testcc "github.com/s7techlab/cckit/testing"
	expectcc "github.com/s7techlab/cckit/testing/expect"
)

func TestCars(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cars Suite")
}

var (
	Authority = testdata.Certificates[0].MustIdentity(`SOME_MSP`)
	Someone   = testdata.Certificates[1].MustIdentity(`SOME_MSP`)
)

var _ = Describe(`Cars`, func() {

	//Create chaincode mock
	cc := testcc.NewMockStub(`cars`, cars.New())
	ccWithoutAC := testcc.NewMockStub(`cars`, cars.NewWithoutAccessControl())

	BeforeSuite(func() {
		// init chaincode
		expectcc.ResponseOk(cc.From(Authority).Init()) // init chaincode from authority
	})

	Describe("Car", func() {

		It("Allow authority to add information about car", func() {
			//invoke chaincode method from authority actor
			expectcc.ResponseOk(cc.From(Authority).Invoke(`carRegister`, cars.Payloads[0]))
		})

		It("Disallow non authority to add information about car", func() {
			//invoke chaincode method from non authority actor
			expectcc.ResponseError(
				cc.From(Someone).Invoke(`carRegister`, cars.Payloads[0]),
				owner.ErrOwnerOnly) // expect "only owner" error
		})

		It("Allow non authority to add information about car to chaincode without access control", func() {
			//invoke chaincode method from non authority actor
			expectcc.ResponseOk(
				ccWithoutAC.From(Authority).Invoke(`carRegister`, cars.Payloads[0]))
		})

		It("Disallow authority to add duplicate information about car", func() {
			expectcc.ResponseError(
				cc.From(Authority).Invoke(`carRegister`, cars.Payloads[0]),
				state.ErrKeyAlreadyExists) //expect car id already exists
		})

		It("Allow everyone to retrieve car information", func() {
			car := expectcc.PayloadIs(cc.Invoke(`carGet`, cars.Payloads[0].Id),
				&cars.Car{}).(cars.Car)

			Expect(car.Title).To(Equal(cars.Payloads[0].Title))
			Expect(car.Id).To(Equal(cars.Payloads[0].Id))
		})

		It("Allow everyone to get car list", func() {
			//  &[]Car{} - declares target type for unmarshalling from []byte received from chaincode
			cc := expectcc.PayloadIs(cc.Invoke(`carList`), &[]cars.Car{}).([]cars.Car)

			Expect(cc).To(HaveLen(1))
			Expect(cc[0].Id).To(Equal(cars.Payloads[0].Id))
		})

		It("Allow authority to add more information about car", func() {
			// register second car
			expectcc.ResponseOk(cc.From(Authority).Invoke(`carRegister`, cars.Payloads[1]))
			cc := expectcc.PayloadIs(
				cc.From(Authority).Invoke(`carList`),
				&[]cars.Car{}).([]cars.Car)

			Expect(cc).To(HaveLen(2))
		})
	})
})
