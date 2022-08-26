package cpaper_proxy

import (
	"github.com/s7techlab/cckit/examples/cpaper_asservice"
	"github.com/s7techlab/cckit/extensions/crosscc"
	"github.com/s7techlab/cckit/router"
)

func NewCCWithLocalCPaper() (*router.Chaincode, error) {
	r := router.New(`crosscc_local`)

	cPaperService := cpaper_asservice.NewService()
	crossCCService := NewServiceWithLocalCPaperResolver(cPaperService)

	// 2 services in one chaincode
	// both CPaper and CrossCC in one chaincode, that is why used local resolver
	if err := cpaper_asservice.RegisterCPaperServiceChaincode(r, cPaperService); err != nil {
		return nil, err
	}

	if err := RegisterCPaperProxyServiceChaincode(r, crossCCService); err != nil {
		return nil, err
	}

	return router.NewChaincode(r), nil
}

func NewCCWithRemoteCPaper() (*router.Chaincode, error) {
	r := router.New(`crosscc_remote`)

	crossCCSettingService := crosscc.NewSettingService()
	crossCCService := NewServiceWithRemoteCPaperResolver(crossCCSettingService)

	// crossCC service and CPaper service - in separate chaincodes
	// in crossCC chaincode there are two services:
	// 1. CrossCC itself
	// 2. Setting service to store information where (channel, chaincode) CPaper service located
	if err := crosscc.RegisterSettingServiceChaincode(r, crossCCSettingService); err != nil {
		return nil, err
	}

	if err := RegisterCPaperProxyServiceChaincode(r, crossCCService); err != nil {
		return nil, err
	}

	return router.NewChaincode(r), nil
}
