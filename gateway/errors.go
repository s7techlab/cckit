package gateway

import "errors"

var (
	ErrEventChannelClosed = errors.New(`event channel is closed`)

	// ErrChaincodeNotExists occurs when attempting to invoke a nonexistent external chaincode
	ErrChaincodeNotExists = errors.New(`chaincode not exists`)

	// ErrSignerNotDefinedInContext msp.SigningIdentity is not defined in context
	ErrSignerNotDefinedInContext = errors.New(`signer is not defined in context`)

	// ErrUnknownInvocationType query or invoke
	ErrUnknownInvocationType = errors.New(`unknown invocation type`)
)
