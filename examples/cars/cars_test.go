package cars

import (
	"testing"

	examplecert "github.com/s7techlab/cckit/examples/cert"
	"github.com/s7techlab/cckit/extensions/owner"
	"github.com/s7techlab/cckit/state"
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
	ccWithoutAC := testcc.NewMockStub(`cars`, NewWithoutAccessControl())

	// load actor certificates
	actors, err := testcc.IdentitiesFromFiles(`SOME_MSP`, map[string]string{
		`authority`: `s7techlab.pem`,
		`someone`:   `victor-nosov.pem`}, examplecert.Content)
	if err != nil {
		panic(err)
	}

	BeforeSuite(func() {
		// init chaincode
		expectcc.ResponseOk(cc.From(actors[`authority`]).Init()) // init chaincode from authority
	})

	Describe("Car", func() {

		It("Allow authority to add information about car", func() {
			//invoke chaincode method from authority actor
			expectcc.ResponseOk(cc.From(actors[`authority`]).Invoke(`carRegister`, Payloads[0]))
		})

		It("Disallow non authority to add information about car", func() {
			//invoke chaincode method from non authority actor
			expectcc.ResponseError(
				cc.From(actors[`someone`]).Invoke(`carRegister`, Payloads[0]),
				owner.ErrOwnerOnly) // expect "only owner" error
		})

		It("Allow non authority to add information about car to chaincode without access control", func() {
			//invoke chaincode method from non authority actor
			expectcc.ResponseOk(
				ccWithoutAC.From(actors[`someone`]).Invoke(`carRegister`, Payloads[0]))
		})

		It("Disallow authority to add duplicate information about car", func() {
			expectcc.ResponseError(
				cc.From(actors[`authority`]).Invoke(`carRegister`, Payloads[0]),
				state.ErrKeyAlreadyExists) //expect car id already exists
		})

		It("Allow everyone to retrieve car information", func() {
			car := expectcc.PayloadIs(cc.Invoke(`carGet`, Payloads[0].Id),
				&Car{}).(Car)

			Expect(car.Title).To(Equal(Payloads[0].Title))
			Expect(car.Id).To(Equal(Payloads[0].Id))
		})

		It("Allow everyone to get car list", func() {
			//  &[]Car{} - declares target type for unmarshalling from []byte received from chaincode
			cars := expectcc.PayloadIs(cc.Invoke(`carList`), &[]Car{}).([]Car)

			Expect(len(cars)).To(Equal(1))
			Expect(cars[0].Id).To(Equal(Payloads[0].Id))
		})

		It("Allow authority to add more information about car", func() {
			// register second car
			expectcc.ResponseOk(cc.From(actors[`authority`]).Invoke(`carRegister`, Payloads[1]))
			cars := expectcc.PayloadIs(
				cc.From(actors[`authority`]).Invoke(`carList`),
				&[]Car{}).([]Car)

			Expect(len(cars)).To(Equal(2))
		})
	})
})
