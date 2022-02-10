package cpaper_proxy

import (
	"fmt"

	cpservice "github.com/s7techlab/cckit/examples/cpaper_asservice"
	"github.com/s7techlab/cckit/extensions/crosscc"
	"github.com/s7techlab/cckit/router"
)

type (
	CPaperProxyService struct {
		CPaperServiceResolver cpservice.CPaperServiceChaincodeResolver
	}
)

// NewServiceWithLocalCPaperResolver - crosscc service and cpaper service in one chaincode
func NewServiceWithLocalCPaperResolver(cpaperService cpservice.CPaperServiceChaincode) *CPaperProxyService {
	return &CPaperProxyService{
		CPaperServiceResolver: cpservice.NewCPaperServiceChaincodeLocalResolver(cpaperService),
	}
}

func NewServiceWithRemoteCPaperResolver(setting crosscc.SettingServiceChaincode) *CPaperProxyService {
	return &CPaperProxyService{
		CPaperServiceResolver: cpservice.NewCPaperServiceChaincodeResolver(crosscc.LocatorResolver(setting)),
	}
}

func (c *CPaperProxyService) GetFromCPaper(ctx router.Context, id *Id) (*InfoFromCPaper, error) {
	cpaperService, err := c.CPaperServiceResolver.Resolve(ctx)
	if err != nil {
		return nil, fmt.Errorf(`resolve Commercial Paper service: %w`, err)
	}
	// It can be cross chaincode invocation or local, if commercial paper service works in same chaincode
	cpaper, err := cpaperService.Get(ctx, &cpservice.CommercialPaperId{Issuer: id.Issuer, PaperNumber: id.PaperNumber})
	if err != nil {
		return nil, fmt.Errorf(`get commercial paper from service: %w`, err)
	}

	return &InfoFromCPaper{
		Issuer:      cpaper.Issuer,
		PaperNumber: cpaper.PaperNumber,
		Owner:       cpaper.Owner,
	}, nil
}
