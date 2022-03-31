package testdata

import (
	"github.com/golang/protobuf/proto"
	. "github.com/onsi/gomega"

	"github.com/s7techlab/cckit/examples/fabcar"
)

type (
	MakerSample struct {
		Create *fabcar.CreateMakerRequest
	}
)

func (ms *MakerSample) CreateClone() *fabcar.CreateMakerRequest {
	return proto.Clone(ms.Create).(*fabcar.CreateMakerRequest)
}

func (ms *MakerSample) ExpectEqual(maker *fabcar.Maker) {
	Expect(ms.Create.Name).To(Equal(maker.Name))
	Expect(ms.Create.Country).To(Equal(maker.Country))
	Expect(ms.Create.FoundationYear).To(Equal(maker.FoundationYear))
}

var (
	MakerNonexistent = MakerSample{
		Create: &fabcar.CreateMakerRequest{
			Name: "Nonexistent",
			Country: "Nonexistent",
			FoundationYear: 1884,
		},
	}

	MakerToyota = MakerSample{
		Create: &fabcar.CreateMakerRequest{
			Name: "Toyota",
			Country: "Japan",
			FoundationYear: 1937,
		},
	}

	MakerAudi = MakerSample{
		Create: &fabcar.CreateMakerRequest{
			Name: "Audi",
			Country: "German",
			FoundationYear: 1909,
		},
	}

	MakerPeugeot = MakerSample{
		Create: &fabcar.CreateMakerRequest{
			Name: "Peugeot",
			Country: "France",
			FoundationYear: 1886,
		},
	}

	MakerFord = MakerSample{
		Create: &fabcar.CreateMakerRequest{
			Name: "Ford",
			Country: "USA",
			FoundationYear: 1903,
		},
	}
)
