package fabcar

import (
	"strconv"

	"github.com/s7techlab/cckit/router"
	"github.com/s7techlab/cckit/state/mapping"
)

type (
	Mapper struct {
		// main entry
		Car *Car

		// secondary entities
		Owners  []*CarOwner
		Details []*CarDetail

		State *mapping.EntryMapper
	}
)

func (m *Mapper) View() *CarView {
	return &CarView{
		Car:     m.Car,
		Owners:  &CarOwners{Items: m.Owners},
		Details: &CarDetails{Items: m.Details},
	}
}

func CreateCarID(car *Car) []string {
	return []string{car.Make, car.Model, strconv.Itoa(int(car.Number))}
}

func (m *Mapper) CreateCar(ctx router.Context, req *CreateCarRequest) error {
	updatedAt, _ := ctx.Stub().GetTxTimestamp()

	m.Car = &Car{
		Make:           req.Make,
		Model:          req.Model,
		Colour:         req.Colour,
		Number:         req.Number,
		OwnersQuantity: uint64(len(req.Owners)),
		UpdatedAt:      updatedAt,
	}

	m.Car.Id = CreateCarID(m.Car)

	if err := m.SetCarOwners(ctx, req.Owners); err != nil {
		return err
	}

	if err := m.SetCarDetails(ctx, req.Details); err != nil {
		return err
	}

	m.State.Commands.Insert(m.Car)
	m.State.Event = mapping.EventFromPayload(&CarCreated{
		Id:             m.Car.Id,
		Make:           m.Car.Make,
		Model:          m.Car.Model,
		Colour:         m.Car.Colour,
		Number:         m.Car.Number,
		OwnersQuantity: m.Car.OwnersQuantity,
	})

	return nil
}

func (m *Mapper) SetCar(ctx router.Context, req *UpdateCarRequest) error {
	updatedAt, _ := ctx.Stub().GetTxTimestamp()
	m.Car.UpdatedAt = updatedAt
	m.Car.Colour = req.Color
	m.Car.Number = req.Number

	if len(req.Owners) > 0 {
		if err := m.SetCarOwners(ctx, req.Owners); err != nil {
			return err
		}

		m.Car.OwnersQuantity = uint64(len(m.Owners))
	}

	if len(req.Details) > 0 {
		if err := m.SetCarDetails(ctx, req.Details); err != nil {
			return err
		}
	}

	m.State.Commands.Put(m.Car)
	m.State.Event = &mapping.Event{Payload: &CarUpdated{
		Id:             m.Car.Id,
		Colour:         m.Car.Colour,
		Number:         m.Car.Number,
		OwnersQuantity: m.Car.OwnersQuantity,
	}}
	return nil
}

func (m *Mapper) SetCarOwners(ctx router.Context, reqs []*SetCarOwner) error {
	updatedAt, _ := ctx.Stub().GetTxTimestamp()

	for _, req := range reqs {
		carOwner := &CarOwner{
			CarId:           m.Car.Id,
			FirstName:       req.FirstName,
			SecondName:      req.SecondName,
			VehiclePassport: req.VehiclePassport,
			UpdatedAt:       updatedAt,
		}

		var exists bool
		for i, owner := range m.Owners {
			if owner.FirstName == carOwner.FirstName && owner.SecondName == carOwner.SecondName {
				m.Owners[i] = carOwner
				exists = true
			}
		}

		if !exists {
			m.Owners = append(m.Owners, carOwner)
		}

		m.State.Commands.Put(carOwner)
	}

	m.State.Event = &mapping.Event{Payload: &CarOwnersUpdated{
		Owners: &CarOwners{Items: m.Owners},
	}}

	return nil
}

func (m *Mapper) SetCarDetails(ctx router.Context, reqs []*SetCarDetail) error {
	updatedAt, _ := ctx.Stub().GetTxTimestamp()

	for _, req := range reqs {
		carDetail := &CarDetail{
			CarId:     m.Car.Id,
			Type:      req.Type,
			Make:      req.Make,
			UpdatedAt: updatedAt,
		}

		var exists bool
		for i, room := range m.Details {
			if room.Type == carDetail.Type {
				m.Details[i] = carDetail
				exists = true
			}
		}

		if !exists {
			m.Details = append(m.Details, carDetail)
		}

		m.State.Commands.Put(carDetail)
	}

	m.State.Event = &mapping.Event{Payload: &CarDetailsUpdated{
		Details: &CarDetails{Items: m.Details},
	}}

	return nil
}
