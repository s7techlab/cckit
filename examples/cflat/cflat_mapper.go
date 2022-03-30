package cflat

import (
	"strconv"

	"github.com/s7techlab/cckit/router"
	"github.com/s7techlab/cckit/state/mapping"
)

type (
	Mapper struct {
		// main entry
		Flat *Flat

		// secondary entities
		Residents []*FlatResident
		Rooms     []*FlatRoom

		State *mapping.EntryMapper
	}
)

func (m *Mapper) View() *FlatView {
	return &FlatView{
		Flat:      m.Flat,
		Residents: &FlatResidents{Residents: m.Residents},
		Rooms:     &FlatRooms{Rooms: m.Rooms},
	}
}

func createFlatID(flat *Flat) []string {
	return []string{flat.Country, flat.Region, flat.City, flat.Street, strconv.Itoa(int(flat.HouseNum)), strconv.Itoa(int(flat.FlatNum))}
}

func (m *Mapper) CreateFlat(ctx router.Context, req *CreateFlatRequest) error {
	updatedAt, _ := ctx.Stub().GetTxTimestamp()

	m.Flat = &Flat{
		Id:                createFlatID(m.Flat),
		Country:           req.Country,
		Region:            req.Region,
		City:              req.City,
		Street:            req.Street,
		HouseNum:          req.HouseNum,
		FlatNum:           req.FlatNum,
		Type:              req.Type,
		ResidentsQuantity: uint64(len(req.Residents)),
		UpdatedAt:         updatedAt,
	}

	var area uint64
	for _, room := range req.Rooms {
		area += room.Area
	}
	m.Flat.Area = area

	if err := m.SetFlatResidents(ctx, req.Residents); err != nil {
		return err
	}

	if err := m.SetFlatRooms(ctx, req.Rooms); err != nil {
		return err
	}

	m.State.Commands.Insert(m.Flat)
	m.State.Event = mapping.EventFromPayload(&FlatCreated{
		Id:                m.Flat.Id,
		Country:           m.Flat.Country,
		Region:            m.Flat.Region,
		City:              m.Flat.Region,
		Street:            m.Flat.Street,
		HouseNum:          m.Flat.HouseNum,
		FlatNum:           m.Flat.FlatNum,
		Type:              m.Flat.Type,
		ResidentsQuantity: m.Flat.ResidentsQuantity,
		Area:              m.Flat.Area,
	})

	return nil
}

func (m *Mapper) SetFlat(ctx router.Context, req *UpdateFlatRequest) error {
	updatedAt, _ := ctx.Stub().GetTxTimestamp()
	m.Flat.UpdatedAt = updatedAt
	m.Flat.Type = req.Type

	if len(req.Residents) > 0 {
		if err := m.SetFlatResidents(ctx, req.Residents); err != nil {
			return err
		}

		m.Flat.ResidentsQuantity = uint64(len(m.Residents))
	}

	if len(req.Residents) > 0 {
		if err := m.SetFlatRooms(ctx, req.Rooms); err != nil {
			return err
		}

		var area uint64
		for _, room := range m.Rooms {
			area += room.Area
		}

		m.Flat.Area = area
	}

	m.State.Commands.Put(m.Flat)
	m.State.Event = &mapping.Event{Payload: &FlatUpdated{
		Id:                m.Flat.Id,
		Type:              m.Flat.Type,
		ResidentsQuantity: m.Flat.ResidentsQuantity,
		Area:              m.Flat.Area,
	}}
	return nil
}

func (m *Mapper) SetFlatResidents(ctx router.Context, reqs []*SetFlatResident) error {
	updatedAt, _ := ctx.Stub().GetTxTimestamp()

	for _, req := range reqs {
		flatResident := &FlatResident{
			FlatId:     m.Flat.Id,
			FirstName:  req.FirstName,
			SecondName: req.SecondName,
			UpdatedAt:  updatedAt,
		}

		var exists bool
		for i, resident := range m.Residents {
			if resident.FirstName == flatResident.FirstName && resident.SecondName == flatResident.SecondName {
				m.Residents[i] = flatResident
				exists = true
			}
		}

		if !exists {
			m.Residents = append(m.Residents, flatResident)
		}

		m.State.Commands.Put(flatResident)
	}

	m.State.Event = &mapping.Event{Payload: &FlatResidentsUpdated{
		FlatId:    m.Flat.Id,
		Residents: &FlatResidents{Residents: m.Residents},
	}}

	return nil
}

func (m *Mapper) SetFlatRooms(ctx router.Context, reqs []*SetFlatRoom) error {
	updatedAt, _ := ctx.Stub().GetTxTimestamp()

	for _, req := range reqs {
		flatRoom := &FlatRoom{
			FlatId:    m.Flat.Id,
			Type:      req.Type,
			Area:      req.Area,
			UpdatedAt: updatedAt,
		}

		var exists bool
		for i, room := range m.Rooms {
			if room.Type == flatRoom.Type {
				m.Rooms[i] = flatRoom
				exists = true
			}
		}

		if !exists {
			m.Rooms = append(m.Rooms, flatRoom)
		}

		m.State.Commands.Put(flatRoom)
	}

	m.State.Event = &mapping.Event{Payload: &FlatRoomsUpdated{
		FlatId: m.Flat.Id,
		Rooms:  &FlatRooms{Rooms: m.Rooms},
	}}

	return nil
}
