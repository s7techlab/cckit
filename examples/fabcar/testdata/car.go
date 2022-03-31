package testdata

import (
	"strconv"

	"github.com/golang/protobuf/proto"
	. "github.com/onsi/gomega"

	"github.com/s7techlab/cckit/examples/fabcar"
)

type (
	CarSample struct {
		Create        *fabcar.CreateCarRequest
		Updates       []*fabcar.UpdateCarRequest
		UpdateOwners  []*fabcar.UpdateCarOwnersRequest
		UpdateDetails []*fabcar.UpdateCarDetailsRequest
	}
)

func (cs CarSample) IdStrings() []string {
	return fabcar.CreateCarID(&fabcar.Car{Make: cs.Create.Make, Model: cs.Create.Model, Number: cs.Create.Number})
}

func (cs CarSample) ExpectCreateEqualCarView(cv *fabcar.CarView) {
	Expect(cv.Car.Id).To(Equal(cs.IdStrings()))

	Expect(cv.Car.Make).To(Equal(cs.Create.Make))
	Expect(cv.Car.Model).To(Equal(cs.Create.Model))
	Expect(cv.Car.Colour).To(Equal(cs.Create.Colour))
	Expect(cv.Car.Number).To(Equal(cs.Create.Number))
	Expect(cv.Car.OwnersQuantity).To(Equal(uint64(len(cs.Create.Owners))))

	ExpectOwnersViewContain(cs.Create.Owners, cv.Owners.Items)

	ExpectDetailsViewContain(cs.Create.Details, cv.Details.Items)
}

func (cs CarSample) ExpectCreateEqualCar(car *fabcar.Car) {
	Expect(car.Id).To(Equal(cs.IdStrings()))

	Expect(car.Make).To(Equal(cs.Create.Make))
	Expect(car.Model).To(Equal(cs.Create.Model))
	Expect(car.Colour).To(Equal(cs.Create.Colour))
	Expect(car.Number).To(Equal(cs.Create.Number))
	Expect(car.OwnersQuantity).To(Equal(uint64(len(cs.Create.Owners))))
}

func (cs CarSample) CreateClone() *fabcar.CreateCarRequest {
	return proto.Clone(cs.Create).(*fabcar.CreateCarRequest)
}

func (cs CarSample) Clone() CarSample {
	return CarSample{
		Create: cs.CreateClone(),
	}
}

var (
	Car1Create = &fabcar.CreateCarRequest{
		Make:   MakerToyota.Create.Name,
		Model:  "Prius",
		Colour: "blue",
		Number: 85322,
		Owners: []*fabcar.SetCarOwner{
			{
				FirstName:       "Tomoko",
				SecondName:      "Uemura",
				VehiclePassport: "Xsdkk4300FSa",
			},
		},
		Details: []*fabcar.SetCarDetail{
			{
				Type: fabcar.DetailType_WHEELS,
				Make: "Michelin",
			},
			{
				Type: fabcar.DetailType_BATTERY,
				Make: "BYD",
			},
		},
	}

	Car1ID = []string{Car1Create.Make, Car1Create.Model, strconv.Itoa(int(Car1Create.Number))}

	Car1 = CarSample{
		Create: Car1Create,
		Updates: []*fabcar.UpdateCarRequest{{
			Id:     Car1ID,
			Color:  "black",
			Number: 333211124,
			Owners: []*fabcar.SetCarOwner{{
				FirstName:       Car1Create.Owners[0].FirstName,
				SecondName:      Car1Create.Owners[0].SecondName,
				VehiclePassport: "Cok1239Dlk13p",
			},
				{
					FirstName:       "Michel",
					SecondName:      "Tailor",
					VehiclePassport: "daj12OkDas092cC",
				}},
			Details: []*fabcar.SetCarDetail{{
				Type: fabcar.DetailType_WHEELS,
				Make: "Continental",
			}},
		}},
		UpdateOwners: []*fabcar.UpdateCarOwnersRequest{{
			CarId: Car1ID,
			Owners: []*fabcar.SetCarOwner{{
				FirstName:       Car1Create.Owners[0].FirstName,
				SecondName:      Car1Create.Owners[0].SecondName,
				VehiclePassport: "23Ck7sAqo0Y7td",
			},
				{
					FirstName:       "Valeria",
					SecondName:      "Peach",
					VehiclePassport: "312jjkdASd98J87d",
				}},
		}},
		UpdateDetails: []*fabcar.UpdateCarDetailsRequest{{
			CarId: Car1ID,
			Details: []*fabcar.SetCarDetail{{
				Type: fabcar.DetailType_WHEELS,
				Make: "Yokohama",
			},
				{
					Type: fabcar.DetailType_BATTERY,
					Make: "Contemporary Amperex Technology",
				}},
		}},
	}

	Car2Create = &fabcar.CreateCarRequest{
		Make:   MakerFord.Create.Name,
		Model:  "Mustang",
		Colour: "red",
		Number: 85322,
		Owners: []*fabcar.SetCarOwner{{
			FirstName:       "Brad",
			SecondName:      "McDonald",
			VehiclePassport: "Iuuu7o9722C",
		},
			{
				FirstName:       "Adriana",
				SecondName:      "Grande",
				VehiclePassport: "9972jjaq812k",
			}},
		Details: []*fabcar.SetCarDetail{{
			Type: fabcar.DetailType_WHEELS,
			Make: "Pirelli",
		},
			{
				Type: fabcar.DetailType_BATTERY,
				Make: "Panasonic",
			}},
	}

	Car2 = CarSample{
		Create: Car2Create,
	}
)
