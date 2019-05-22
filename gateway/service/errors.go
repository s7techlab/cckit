package service

import "errors"

var (
	// ErrChaincodeNotExists occurs when attempting to invoke a nonexostent external chaincode
	ErrChaincodeNotExists = errors.New(`chaincode not exists`)

	ErrSignerNotDefinedInContext = errors.New(`signer is not defined in context`)
)
