package fabcar_test

import (
	"testing"

	"github.com/golang/protobuf/ptypes/empty"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/s7techlab/cckit/examples/fabcar"
	"github.com/s7techlab/cckit/examples/fabcar/testdata"
	identitytestdata "github.com/s7techlab/cckit/identity/testdata"
	testcc "github.com/s7techlab/cckit/testing"
)

func TestFabCarService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Flight service suite")
}

func ExpectErrorContain(str string) func(interface{}, error) {
	return func(_ interface{}, err error) {
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring(str))
	}
}

var _ = Describe(`FabCar service`, func() {

	var (
		fabCarCC = &fabcar.FabCarService{}

		// same ctx for all related services - like all services in one chaincode
		cc, ctx = testcc.NewTxHandler(fabcar.ChaincodeName)

		car1IdString = testdata.Car1.IdStrings()
	)

	It("Allow to init related", func() {
		cc.From(identitytestdata.Certificates[0].MustIdentity(`SOME_MSP`)).Init(fabcar.ChaincodeInitFunc())
	})

	Context("Maker", func() {
		It("disallow to create maker", func() {
			cc.Tx(func() {
				ExpectErrorContain(`invalid field FoundationYear: value '1884' must be greater than '1885'`)(
					fabCarCC.CreateMaker(ctx, testdata.MakerNonexistent.Create))
			})
		})

		It("allow to create makers", func() {
			cc.Tx(func() {
				reqCreate := testdata.MakerNonexistent.CreateClone()
				reqCreate.FoundationYear = 1886

				maker, err := fabCarCC.CreateMaker(ctx, reqCreate)
				Expect(err).NotTo(HaveOccurred())
				Expect(maker.Name).To(Equal(reqCreate.Name))
				Expect(maker.Country).To(Equal(reqCreate.Country))
				Expect(maker.FoundationYear).To(Equal(reqCreate.FoundationYear))

				req := testdata.MakerToyota
				maker, err = fabCarCC.CreateMaker(ctx, req.Create)
				Expect(err).NotTo(HaveOccurred())
				req.ExpectEqual(maker)

				req = testdata.MakerAudi
				maker, err = fabCarCC.CreateMaker(ctx, req.Create)
				Expect(err).NotTo(HaveOccurred())
				req.ExpectEqual(maker)

				req = testdata.MakerPeugeot
				maker, err = fabCarCC.CreateMaker(ctx, req.Create)
				Expect(err).NotTo(HaveOccurred())
				req.ExpectEqual(maker)

				req = testdata.MakerFord
				maker, err = fabCarCC.CreateMaker(ctx, req.Create)
				Expect(err).NotTo(HaveOccurred())
				req.ExpectEqual(maker)
			})
		})

		It("allow to get maker", func() {
			cc.Tx(func() {
				maker, err := fabCarCC.GetMaker(ctx, &fabcar.MakerName{Name: testdata.MakerPeugeot.Create.Name})
				Expect(err).NotTo(HaveOccurred())
				testdata.MakerPeugeot.ExpectEqual(maker)
			})
		})

		It("allow to get list makers", func() {
			cc.Tx(func() {
				makers, err := fabCarCC.ListMakers(ctx, &empty.Empty{})
				Expect(err).NotTo(HaveOccurred())
				Expect(makers.Items).To(HaveLen(5))
			})
		})

		It("allow to delete maker", func() {
			cc.Tx(func() {
				_, err := fabCarCC.DeleteMaker(ctx, &fabcar.MakerName{Name: testdata.MakerNonexistent.Create.Name})
				Expect(err).NotTo(HaveOccurred())
			})

			cc.Tx(func() {
				makers, err := fabCarCC.ListMakers(ctx, &empty.Empty{})
				Expect(err).NotTo(HaveOccurred())
				Expect(makers.Items).To(HaveLen(4))
			})
		})
	})

	Context(`Car`, func() {
		It("disallow to create car", func() {
			cc.Tx(func() {
				carReq := testdata.Car1.CreateClone()

				carReq.Make = testdata.MakerNonexistent.Create.Name
				ExpectErrorContain(`maker is not created: state entry not found: Maker | Nonexistent`)(fabCarCC.CreateCar(ctx, carReq))
			})
		})

		It("allow to create cars", func() {
			cc.Tx(func() {
				req := testdata.Car1
				carView, err := fabCarCC.CreateCar(ctx, req.Create)
				Expect(err).NotTo(HaveOccurred())
				req.ExpectCreateEqualCarView(carView)

				req = testdata.Car2
				carView, err = fabCarCC.CreateCar(ctx, req.Create)
				Expect(err).NotTo(HaveOccurred())
				req.ExpectCreateEqualCarView(carView)

				reqCreate := testdata.Car1.CreateClone()
				reqCreate.Make = testdata.MakerAudi.Create.Name
				_, err = fabCarCC.CreateCar(ctx, reqCreate)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		It("allow to get car", func() {
			cc.Tx(func() {
				car, err := fabCarCC.GetCar(ctx, &fabcar.CarId{Id: testdata.Car2.IdStrings()})
				Expect(err).NotTo(HaveOccurred())
				testdata.Car2.ExpectCreateEqualCar(car)
			})

			cc.Tx(func() {
				carOwners, err := fabCarCC.ListCarOwners(ctx, &fabcar.CarId{Id: car1IdString})
				Expect(err).NotTo(HaveOccurred())
				Expect(carOwners.Items).To(HaveLen(1))
			})

			cc.Tx(func() {
				carDetails, err := fabCarCC.ListCarDetails(ctx, &fabcar.CarId{Id: car1IdString})
				Expect(err).NotTo(HaveOccurred())
				Expect(carDetails.Items).To(HaveLen(2))
			})
		})

		It("allow to get car view", func() {
			cc.Tx(func() {
				carView, err := fabCarCC.GetCarView(ctx, &fabcar.CarId{Id: testdata.Car2.IdStrings()})
				Expect(err).NotTo(HaveOccurred())
				testdata.Car2.ExpectCreateEqualCarView(carView)
			})
		})

		It("allow to get list cars", func() {
			cc.Tx(func() {
				cars, err := fabCarCC.ListCars(ctx, &empty.Empty{})
				Expect(err).NotTo(HaveOccurred())
				Expect(cars.Items).To(HaveLen(3))
			})
		})

		It("allow to get car owner", func() {
			cc.Tx(func() {
				carOwner, err := fabCarCC.GetCarOwner(ctx, &fabcar.CarOwnerId{
					CarId:      car1IdString,
					FirstName:  testdata.Car1.Create.Owners[0].FirstName,
					SecondName: testdata.Car1.Create.Owners[0].SecondName,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(carOwner.CarId).To(Equal(car1IdString))
				Expect(carOwner.FirstName).To(Equal(testdata.Car1.Create.Owners[0].FirstName))
				Expect(carOwner.SecondName).To(Equal(testdata.Car1.Create.Owners[0].SecondName))
			})
		})

		It("allow to get list car owners", func() {
			cc.Tx(func() {
				carOwners, err := fabCarCC.ListCarOwners(ctx, &fabcar.CarId{Id: car1IdString})
				Expect(err).NotTo(HaveOccurred())
				Expect(carOwners.Items).To(HaveLen(1))
			})
		})

		It("allow to get car detail", func() {
			cc.Tx(func() {
				carDetail, err := fabCarCC.GetCarDetail(ctx, &fabcar.CarDetailId{
					CarId: car1IdString,
					Type:  testdata.Car1.Create.Details[0].Type,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(carDetail.CarId).To(Equal(car1IdString))
				Expect(carDetail.Type).To(Equal(testdata.Car1.Create.Details[0].Type))
				Expect(carDetail.Make).To(Equal(testdata.Car1.Create.Details[0].Make))
			})
		})

		It("allow to get list car details", func() {
			cc.Tx(func() {
				carDetails, err := fabCarCC.ListCarDetails(ctx, &fabcar.CarId{Id: car1IdString})
				Expect(err).NotTo(HaveOccurred())
				Expect(carDetails.Items).To(HaveLen(2))
			})
		})

		It("allow to delete car", func() {
			cc.Tx(func() {
				req := testdata.Car1.Clone()
				req.Create.Make = testdata.MakerAudi.Create.Name

				carView, err := fabCarCC.DeleteCar(ctx, &fabcar.CarId{Id: req.IdStrings()})
				Expect(err).NotTo(HaveOccurred())
				req.ExpectCreateEqualCarView(carView)
			})
		})

		It("allow to update car", func() {
			cc.Tx(func() {
				carView, err := fabCarCC.UpdateCar(ctx, testdata.Car1.Updates[0])
				Expect(err).NotTo(HaveOccurred())
				Expect(carView.Car.Id).To(Equal(car1IdString))
				Expect(carView.Car.Colour).To(Equal(testdata.Car1.Updates[0].Color))
				Expect(carView.Car.Number).To(Equal(testdata.Car1.Updates[0].Number))

				Expect(carView.Owners.Items[0].CarId).To(Equal(car1IdString))
				Expect(carView.Owners.Items[0].FirstName).To(Equal(testdata.Car1.Updates[0].Owners[0].FirstName))
				Expect(carView.Owners.Items[0].SecondName).To(Equal(testdata.Car1.Updates[0].Owners[0].SecondName))
				Expect(carView.Owners.Items[0].VehiclePassport).To(Equal(testdata.Car1.Updates[0].Owners[0].VehiclePassport))

				Expect(carView.Owners.Items[1].CarId).To(Equal(car1IdString))
				Expect(carView.Owners.Items[1].FirstName).To(Equal(testdata.Car1.Updates[0].Owners[1].FirstName))
				Expect(carView.Owners.Items[1].SecondName).To(Equal(testdata.Car1.Updates[0].Owners[1].SecondName))
				Expect(carView.Owners.Items[1].VehiclePassport).To(Equal(testdata.Car1.Updates[0].Owners[1].VehiclePassport))

				Expect(carView.Details.Items[0].CarId).To(Equal(car1IdString))
				Expect(carView.Details.Items[0].Type).To(Equal(testdata.Car1.Create.Details[1].Type))
				Expect(carView.Details.Items[0].Make).To(Equal(testdata.Car1.Create.Details[1].Make))

				Expect(carView.Details.Items[1].CarId).To(Equal(car1IdString))
				Expect(carView.Details.Items[1].Type).To(Equal(testdata.Car1.Updates[0].Details[0].Type))
				Expect(carView.Details.Items[1].Make).To(Equal(testdata.Car1.Updates[0].Details[0].Make))
			})
		})

		It("allow to update car owners", func() {
			cc.Tx(func() {
				carOwners, err := fabCarCC.UpdateCarOwners(ctx, testdata.Car1.UpdateOwners[0])
				Expect(err).NotTo(HaveOccurred())
				Expect(carOwners.Items[0].CarId).To(Equal(car1IdString))
				Expect(carOwners.Items[0].FirstName).To(Equal(testdata.Car1.Updates[0].Owners[1].FirstName))
				Expect(carOwners.Items[0].SecondName).To(Equal(testdata.Car1.Updates[0].Owners[1].SecondName))
				Expect(carOwners.Items[0].VehiclePassport).To(Equal(testdata.Car1.Updates[0].Owners[1].VehiclePassport))

				Expect(carOwners.Items[1].CarId).To(Equal(car1IdString))
				Expect(carOwners.Items[1].FirstName).To(Equal(testdata.Car1.UpdateOwners[0].Owners[0].FirstName))
				Expect(carOwners.Items[1].SecondName).To(Equal(testdata.Car1.UpdateOwners[0].Owners[0].SecondName))
				Expect(carOwners.Items[1].VehiclePassport).To(Equal(testdata.Car1.UpdateOwners[0].Owners[0].VehiclePassport))

				Expect(carOwners.Items[2].CarId).To(Equal(car1IdString))
				Expect(carOwners.Items[2].FirstName).To(Equal(testdata.Car1.UpdateOwners[0].Owners[1].FirstName))
				Expect(carOwners.Items[2].SecondName).To(Equal(testdata.Car1.UpdateOwners[0].Owners[1].SecondName))
				Expect(carOwners.Items[2].VehiclePassport).To(Equal(testdata.Car1.UpdateOwners[0].Owners[1].VehiclePassport))
			})

			cc.Tx(func() {
				carOwners, err := fabCarCC.ListCarOwners(ctx, &fabcar.CarId{Id: car1IdString})
				Expect(err).NotTo(HaveOccurred())
				Expect(carOwners.Items).To(HaveLen(3))
			})
		})

		It("allow to delete car owner", func() {
			cc.Tx(func() {
				carOwner, err := fabCarCC.DeleteCarOwner(ctx, &fabcar.CarOwnerId{
					CarId:      car1IdString,
					FirstName:  testdata.Car1.Updates[0].Owners[1].FirstName,
					SecondName: testdata.Car1.Updates[0].Owners[1].SecondName,
				})

				Expect(err).NotTo(HaveOccurred())
				Expect(carOwner.FirstName).To(Equal(testdata.Car1.Updates[0].Owners[1].FirstName))
				Expect(carOwner.SecondName).To(Equal(testdata.Car1.Updates[0].Owners[1].SecondName))
				Expect(carOwner.VehiclePassport).To(Equal(testdata.Car1.Updates[0].Owners[1].VehiclePassport))
			})

			cc.Tx(func() {
				carOwners, err := fabCarCC.ListCarOwners(ctx, &fabcar.CarId{Id: car1IdString})
				Expect(err).NotTo(HaveOccurred())
				Expect(carOwners.Items).To(HaveLen(2))
			})
		})

		It("allow to update car details", func() {
			cc.Tx(func() {
				carDetails, err := fabCarCC.UpdateCarDetails(ctx, testdata.Car1.UpdateDetails[0])
				Expect(err).NotTo(HaveOccurred())
				Expect(carDetails.Items[0].CarId).To(Equal(car1IdString))
				Expect(carDetails.Items[0].Type).To(Equal(testdata.Car1.UpdateDetails[0].Details[1].Type))
				Expect(carDetails.Items[0].Make).To(Equal(testdata.Car1.UpdateDetails[0].Details[1].Make))

				Expect(carDetails.Items[1].CarId).To(Equal(car1IdString))
				Expect(carDetails.Items[1].Type).To(Equal(testdata.Car1.UpdateDetails[0].Details[0].Type))
				Expect(carDetails.Items[1].Make).To(Equal(testdata.Car1.UpdateDetails[0].Details[0].Make))
			})
		})
	})
})
