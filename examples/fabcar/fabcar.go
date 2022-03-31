package fabcar

import (
	"fmt"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/s7techlab/cckit/router"
	"github.com/s7techlab/cckit/state/mapping"
)

type FabCarService struct{}

func (f *FabCarService) CreateMaker(ctx router.Context, req *CreateMakerRequest) (*Maker, error) {
	if err := router.ValidateRequest(req); err != nil {
		return nil, err
	}

	maker := &Maker{
		Name:           req.Name,
		Country:        req.Country,
		FoundationYear: req.FoundationYear,
	}

	mapper := f.carMapper()
	mapper.State.Commands.Insert(maker)
	mapper.State.Event = mapping.EventFromPayload(&MakerCreated{
		Name:           maker.Name,
		Country:        maker.Country,
		FoundationYear: maker.FoundationYear,
	})

	if err := mapper.State.Apply(State(ctx), Event(ctx)); err != nil {
		return nil, err
	}

	return maker, nil
}

func (f *FabCarService) DeleteMaker(ctx router.Context, name *MakerName) (*Maker, error) {
	maker, err := f.GetMaker(ctx, name)
	if err != nil {
		return nil, err
	}

	if err = State(ctx).Delete(maker); err != nil {
		return nil, err
	}

	if err = Event(ctx).Set(&MakerDeleted{
		Name:           maker.Name,
		Country:        maker.Country,
		FoundationYear: maker.FoundationYear,
	}); err != nil {
		return nil, err
	}

	return maker, nil
}

func (f *FabCarService) GetMaker(ctx router.Context, name *MakerName) (*Maker, error) {
	if err := router.ValidateRequest(name); err != nil {
		return nil, err
	}

	maker, err := State(ctx).Get(name, &Maker{})
	if err != nil {
		return nil, err
	}

	return maker.(*Maker), nil
}

func (f *FabCarService) ListMakers(ctx router.Context, _ *emptypb.Empty) (*Makers, error) {
	res, err := State(ctx).List(&Maker{})
	if err != nil {
		return nil, err
	}

	return res.(*Makers), nil
}

func (f *FabCarService) carMapper() *Mapper {
	return &Mapper{
		State: &mapping.EntryMapper{},
	}
}

func (f *FabCarService) carMapperByView(carView *CarView) *Mapper {
	return &Mapper{
		Car: carView.Car,

		Owners:  carView.Owners.Items,
		Details: carView.Details.Items,

		State: &mapping.EntryMapper{},
	}
}

func (f *FabCarService) CreateCar(ctx router.Context, req *CreateCarRequest) (*CarView, error) {
	if err := router.ValidateRequest(req); err != nil {
		return nil, err
	}

	_, err := f.GetMaker(ctx, &MakerName{Name: req.Make})
	if err != nil {
		return nil, fmt.Errorf("maker is not created: %w", err)
	}

	mapper := f.carMapper()
	if err = mapper.CreateCar(ctx, req); err != nil {
		return nil, err
	}

	if err = mapper.State.Apply(State(ctx), Event(ctx)); err != nil {
		return nil, err
	}

	return mapper.View(), nil
}

func (f *FabCarService) UpdateCar(ctx router.Context, req *UpdateCarRequest) (*CarView, error) {
	if err := router.ValidateRequest(req); err != nil {
		return nil, err
	}

	carView, err := f.GetCarView(ctx, &CarId{Id: req.Id})
	if err != nil {
		return nil, err
	}

	mapper := f.carMapperByView(carView)
	if err = mapper.SetCar(ctx, req); err != nil {
		return nil, err
	}

	if err = mapper.State.Apply(State(ctx), Event(ctx)); err != nil {
		return nil, err
	}

	return mapper.View(), nil
}

func (f *FabCarService) DeleteCar(ctx router.Context, id *CarId) (*CarView, error) {
	if err := router.ValidateRequest(id); err != nil {
		return nil, err
	}

	carView, err := f.GetCarView(ctx, id)
	if err != nil {
		return nil, err
	}

	for _, owner := range carView.Owners.Items {
		if err = State(ctx).Delete(owner); err != nil {
			return nil, err
		}
	}

	for _, details := range carView.Details.Items {
		if err = State(ctx).Delete(details); err != nil {
			return nil, err
		}
	}

	if err = State(ctx).Delete(carView.Car); err != nil {
		return nil, err
	}

	if err = Event(ctx).Set(&CarDeleted{
		Id:             carView.Car.Id,
		Make:           carView.Car.Make,
		Model:          carView.Car.Model,
		Colour:         carView.Car.Colour,
		Number:         carView.Car.Number,
		OwnersQuantity: carView.Car.OwnersQuantity,
		Owners:         carView.Owners,
		Details:        carView.Details,
	}); err != nil {
		return nil, err
	}

	return carView, nil
}

func (f *FabCarService) GetCar(ctx router.Context, id *CarId) (*Car, error) {
	if err := router.ValidateRequest(id); err != nil {
		return nil, err
	}

	car, err := State(ctx).Get(id, &Car{})
	if err != nil {
		return nil, err
	}

	return car.(*Car), nil
}

func (f *FabCarService) GetCarView(ctx router.Context, id *CarId) (*CarView, error) {
	var (
		carView = &CarView{}
		err     error
	)

	carView.Car, err = f.GetCar(ctx, id)
	if err != nil {
		return nil, err
	}

	carView.Owners, err = f.ListCarOwners(ctx, id)
	if err != nil {
		return nil, err
	}

	carView.Details, err = f.ListCarDetails(ctx, id)
	if err != nil {
		return nil, err
	}

	return carView, nil
}

func (f *FabCarService) ListCars(ctx router.Context, _ *emptypb.Empty) (*Cars, error) {
	res, err := State(ctx).List(&Car{})
	if err != nil {
		return nil, err
	}

	return res.(*Cars), nil
}

func (f *FabCarService) UpdateCarOwners(ctx router.Context, req *UpdateCarOwnersRequest) (*CarOwners, error) {
	carView, err := f.GetCarView(ctx, &CarId{Id: req.CarId})
	if err != nil {
		return nil, err
	}

	mapper := f.carMapperByView(carView)
	if err = mapper.SetCarOwners(ctx, req.Owners); err != nil {
		return nil, err
	}

	if err = mapper.State.Apply(State(ctx), Event(ctx)); err != nil {
		return nil, err
	}

	return mapper.View().Owners, nil
}

func (f *FabCarService) DeleteCarOwner(ctx router.Context, id *CarOwnerId) (*CarOwner, error) {
	carOwner, err := f.GetCarOwner(ctx, id)
	if err != nil {
		return nil, err
	}

	if err = State(ctx).Delete(carOwner); err != nil {
		return nil, err
	}

	return carOwner, nil
}

func (f *FabCarService) GetCarOwner(ctx router.Context, id *CarOwnerId) (*CarOwner, error) {
	carOwner, err := State(ctx).Get(id, &CarOwner{})
	if err != nil {
		return nil, err
	}

	return carOwner.(*CarOwner), nil
}

func (f *FabCarService) ListCarOwners(ctx router.Context, id *CarId) (*CarOwners, error) {
	if err := router.ValidateRequest(id); err != nil {
		return nil, err
	}

	carOwner, err := State(ctx).ListWith(&CarOwner{}, id.Id)
	if err != nil {
		return nil, err
	}

	return carOwner.(*CarOwners), nil
}

func (f *FabCarService) UpdateCarDetails(ctx router.Context, req *UpdateCarDetailsRequest) (*CarDetails, error) {
	carView, err := f.GetCarView(ctx, &CarId{Id: req.CarId})
	if err != nil {
		return nil, err
	}

	mapper := f.carMapperByView(carView)
	if err = mapper.SetCarDetails(ctx, req.Details); err != nil {
		return nil, err
	}

	if err = mapper.State.Apply(State(ctx), Event(ctx)); err != nil {
		return nil, err
	}

	return mapper.View().Details, nil
}

func (f *FabCarService) DeleteCarDetail(ctx router.Context, id *CarDetailId) (*CarDetail, error) {
	carDetail, err := f.GetCarDetail(ctx, id)
	if err != nil {
		return nil, err
	}

	if err = State(ctx).Delete(carDetail); err != nil {
		return nil, err
	}

	return carDetail, nil
}

func (f *FabCarService) GetCarDetail(ctx router.Context, id *CarDetailId) (*CarDetail, error) {
	carDetail, err := State(ctx).Get(id, &CarDetail{})
	if err != nil {
		return nil, err
	}

	return carDetail.(*CarDetail), nil
}

func (f *FabCarService) ListCarDetails(ctx router.Context, id *CarId) (*CarDetails, error) {
	if err := router.ValidateRequest(id); err != nil {
		return nil, err
	}

	carDetail, err := State(ctx).ListWith(&CarDetail{}, id.Id)
	if err != nil {
		return nil, err
	}

	return carDetail.(*CarDetails), nil
}
