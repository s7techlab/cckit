package mapping

import "errors"

var (
	// ErrEntryTypeNotSupported entry type has no appropriate mapper type
	ErrEntryTypeNotSupported = errors.New(`entry type not supported for mapping`)

	// ErrEntryTypeNotDefined mapping for entry type not defined
	ErrEntryTypeNotDefined = errors.New(`mapping for entry type not defined`)
)
