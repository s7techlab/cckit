package service

import "errors"

var (
	// ErrChaincodeNotExists occurs when attempting to invoke a nonexostent external chaincode
	ErrChaincodeNotExists = errors.New(`chaincode not exists`)

	// ErrSignerNotDefinedInContext msp.SigningIdentity is not defined in context
	ErrSignerNotDefinedInContext = errors.New(`signer is not defined in context`)

	// ErrUnknownInvocationType query or invoke
	ErrUnknownInvocationType = errors.New(`unknown invocation type`)
)
