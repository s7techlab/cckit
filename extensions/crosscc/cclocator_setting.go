package crosscc

import (
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/s7techlab/cckit/gateway"
	"github.com/s7techlab/cckit/router"
)

type SettingService struct{}

var _ SettingServiceChaincode = &SettingService{}

func NewSettingService() *SettingService {
	return &SettingService{}
}

func LocatorResolver(c SettingServiceChaincode) gateway.ChaincodeLocatorResolver {
	return func(ctx router.Context, service string) (*gateway.ChaincodeLocator, error) {
		locator, err := c.ServiceLocatorGet(ctx, &ServiceLocatorId{Service: service})
		if err != nil {
			return nil, fmt.Errorf("chaincode locator not found, service=%s d: %w", service, err)
		}

		return &gateway.ChaincodeLocator{Channel: locator.Channel, Chaincode: locator.Chaincode}, nil
	}
}

func (c *SettingService) LocatorResolver() gateway.ChaincodeLocatorResolver {
	return LocatorResolver(c)
}

func (c *SettingService) ServiceLocatorSet(ctx router.Context, locatorSet *ServiceLocatorSetRequest) (*ServiceLocator, error) {
	if err := router.ValidateRequest(locatorSet); err != nil {
		return nil, err
	}

	locator := &ServiceLocator{
		Service:   locatorSet.Service,
		Channel:   locatorSet.Channel,
		Chaincode: locatorSet.Chaincode,
	}

	return locator, State(ctx).Put(locator)
}

func (c *SettingService) ServiceLocatorGet(ctx router.Context, id *ServiceLocatorId) (*ServiceLocator, error) {
	if err := router.ValidateRequest(id); err != nil {
		return nil, err
	}

	if res, err := State(ctx).Get(id, &ServiceLocator{}); err != nil {
		return nil, err
	} else {
		return res.(*ServiceLocator), nil
	}
}

func (c *SettingService) ListServiceLocators(ctx router.Context, _ *empty.Empty) (*ServiceLocators, error) {
	if res, err := State(ctx).List(&ServiceLocator{}); err != nil {
		return nil, err
	} else {
		return res.(*ServiceLocators), nil
	}
}

func (c *SettingService) PingService(context router.Context, id *ServiceLocatorId) (*PingServiceResponse, error) {
	panic("implement me")
}

func (c *SettingService) PingServices(context router.Context, empty *empty.Empty) (*PingServiceResponses, error) {
	panic("implement me")
}
