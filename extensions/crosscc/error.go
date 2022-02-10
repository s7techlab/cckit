package crosscc

import "errors"

var (
	ErrServiceNotForLocalChaincodeResolver = errors.New("service not set for local chaincode resolver")
)
