package cpaper_asservice

import (
	"github.com/s7techlab/cckit/extensions/encryption"
	"github.com/s7techlab/cckit/extensions/owner"
	"github.com/s7techlab/cckit/router"
)

func CCRouter(name string) (*router.Group, error) {
	r := router.New(name)
	// Store on the ledger the information about chaincode instantiation
	r.Init(owner.InvokeSetFromCreator)

	if err := RegisterCPaperServiceChaincode(r, &CPaperService{}); err != nil {
		return nil, err
	}

	return r, nil
}

func NewCC() (*router.Chaincode, error) {
	r, err := CCRouter(`CommercialPaper`)
	if err != nil {
		return nil, err
	}

	return router.NewChaincode(r), nil
}

func NewCCEncrypted() (*router.Chaincode, error) {
	r, err := CCRouter(`CommercialPaperEncrypted`)
	if err != nil {
		return nil, err
	}

	r.
		// encryption key in transient map and encrypted args required
		Pre(encryption.ArgsDecrypt).
		// default Context replaced with EncryptedStateContext only if key is provided in transient map
		Use(encryption.EncStateContext).
		// invoke response will be encrypted because it will be placed in blocks
		After(encryption.EncryptInvokeResponse())

	return router.NewChaincode(r), nil
}
