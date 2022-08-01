package router

import "errors"

var (
	// ErrEmptyArgs occurs when trying to invoke chaincode method with empty args
	ErrEmptyArgs = errors.New(`empty args`)

	// ErrMethodNotFound occurs when trying to invoke non-existent chaincode method
	ErrMethodNotFound = errors.New(`chaincode method not found`)

	// ErrArgsNumMismatch occurs when the number of declared and the number of arguments passed does not match
	ErrArgsNumMismatch = errors.New(`chaincode method args count mismatch`)

	// ErrHandlerError error in handler
	ErrHandlerError = errors.New(`router handler error`)
)
