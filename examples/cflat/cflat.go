package cflat

import (
	"github.com/s7techlab/cckit/router"
	"github.com/s7techlab/cckit/state/mapping"
	"google.golang.org/protobuf/types/known/emptypb"
)

type FlatService struct{}

func (f *FlatService) flatMapper() *Mapper {
	return &Mapper{
		State: &mapping.EntryMapper{},
	}
}

func (f *FlatService) flatMapperByView(flatView *FlatView) *Mapper {
	return &Mapper{
		Flat: flatView.Flat,

		Residents: flatView.Residents.Residents,
		Rooms:     flatView.Rooms.Rooms,

		State: &mapping.EntryMapper{},
	}
}

func (f *FlatService) CreateFlat(ctx router.Context, req *CreateFlatRequest) (*FlatView, error) {
	if err := router.ValidateRequest(req); err != nil {
		return nil, err
	}

	mapper := f.flatMapper()

	if err := mapper.CreateFlat(ctx, req); err != nil {
		return nil, err
	}

	if err := mapper.State.Apply(State(ctx), Event(ctx)); err != nil {
		return nil, err
	}

	return mapper.View(), nil
}

func (f *FlatService) UpdateFlat(ctx router.Context, req *UpdateFlatRequest) (*FlatView, error) {
	if err := router.ValidateRequest(req); err != nil {
		return nil, err
	}

	flatView, err := f.GetFlatView(ctx, &FlatId{Id: req.Id})
	if err != nil {
		return nil, err
	}

	mapper := f.flatMapperByView(flatView)

	if err = mapper.SetFlat(ctx, req); err != nil {
		return nil, err
	}

	return mapper.View(), nil
}

func (f *FlatService) DeleteFlat(ctx router.Context, id *FlatId) (*FlatView, error) {
	if err := router.ValidateRequest(id); err != nil {
		return nil, err
	}

	flatView, err := f.GetFlatView(ctx, id)
	if err != nil {
		return nil, err
	}

	for _, resident := range flatView.Residents.Residents {
		if err = State(ctx).Delete(resident); err != nil {
			return nil, err
		}
	}

	for _, room := range flatView.Rooms.Rooms {
		if err = State(ctx).Delete(room); err != nil {
			return nil, err
		}
	}

	if err = State(ctx).Delete(flatView.Flat); err != nil {
		return nil, err
	}

	if err = Event(ctx).Set(&FlatDeleted{
		Id:                flatView.Flat.Id,
		Country:           flatView.Flat.Country,
		Region:            flatView.Flat.Region,
		City:              flatView.Flat.City,
		Street:            flatView.Flat.Street,
		HouseNum:          flatView.Flat.HouseNum,
		FlatNum:           flatView.Flat.FlatNum,
		Type:              flatView.Flat.Type,
		ResidentsQuantity: flatView.Flat.ResidentsQuantity,
		Area:              flatView.Flat.Area,
		Residents:         flatView.Residents,
		Rooms:             flatView.Rooms,
	}); err != nil {
		return nil, err
	}

	return flatView, nil
}

func (f *FlatService) GetFlat(ctx router.Context, id *FlatId) (*Flat, error) {
	if err := router.ValidateRequest(id); err != nil {
		return nil, err
	}

	flat, err := State(ctx).Get(id, &Flat{})
	if err != nil {
		return nil, err
	}
	return flat.(*Flat), nil
}

func (f *FlatService) GetFlatView(ctx router.Context, id *FlatId) (*FlatView, error) {
	var (
		flatView *FlatView
		err      error
	)

	flatView.Flat, err = f.GetFlat(ctx, id)
	if err != nil {
		return nil, err
	}

	flatView.Residents, err = f.ListFlatResidents(ctx, id)
	if err != nil {
		return nil, err
	}

	flatView.Residents, err = f.ListFlatResidents(ctx, id)
	if err != nil {
		return nil, err
	}

	flatView.Rooms, err = f.ListFlatRooms(ctx, id)
	if err != nil {
		return nil, err
	}

	return flatView, nil
}

func (f *FlatService) ListFlats(ctx router.Context, _ *emptypb.Empty) (*Flats, error) {
	panic("implement me")
}

func (f *FlatService) UpdateFlatResident(ctx router.Context, req *UpdateFlatResidentRequest) (*FlatResident, error) {
	panic("implement me")
}

func (f *FlatService) DeleteFlatResident(ctx router.Context, id *FlatResidentId) (*FlatResident, error) {
	panic("implement me")
}

func (f *FlatService) GetFlatResident(ctx router.Context, id *FlatResidentId) (*FlatResident, error) {
	panic("implement me")
}

func (f *FlatService) ListFlatResidents(ctx router.Context, id *FlatId) (*FlatResidents, error) {
	panic("implement me")
}

func (f *FlatService) UpdateFlatRoom(ctx router.Context, req *UpdateFlatRoomRequest) (*FlatRoom, error) {
	panic("implement me")
}

func (f *FlatService) DeleteFlatRoom(ctx router.Context, id *FlatRoomId) (*FlatRoom, error) {
	panic("implement me")
}

func (f *FlatService) GetFlatRoom(ctx router.Context, id *FlatRoomId) (*FlatRoom, error) {
	panic("implement me")
}

func (f *FlatService) ListFlatRooms(ctx router.Context, id *FlatId) (*FlatRooms, error) {
	panic("implement me")
}
