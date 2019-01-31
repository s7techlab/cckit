package state

import (
	"github.com/pkg/errors"
)

var (
	// ErrUnableToCreateKey can occurs while creating composite key for entry
	ErrUnableToCreateKey = errors.New(`unable to create state key`)

	// ErrKeyAlreadyExists can occurs when trying to insert entry with existing key
	ErrKeyAlreadyExists = errors.New(`state key already exists`)

	// ErrrKeyNotFound key not found in chaincode state
	ErrKeyNotFound = errors.New(`state entry not found`)

	// ErrAllowOnlyOneValue can occurs when trying to call Insert or Put with more than 2 arguments
	ErrAllowOnlyOneValue = errors.New(`allow only one value`)

	// ErrKeyNotSupportKeyerInterface can occurs when trying to Insert or Put struct without providing key and struct not support Keyer interface
	ErrKeyNotSupportKeyerInterface = errors.New(`key not support keyer interface`)

	// ErrKeyPartsLength can occurs when trying to create key consisting of zero parts
	ErrKeyPartsLength = errors.New(`key parts length must be greater than zero`)

	ErrEntryTypeMappingNotSupported = errors.New(`entry type mapping not supported`)

	ErrEntryMappingNotDefined = errors.New(`entry mapping not defined`)
)
