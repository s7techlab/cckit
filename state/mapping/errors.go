package mapping

import "errors"

var (
	ErrEntryTypeMappingNotSupported = errors.New(`entry type mapping not supported`)

	ErrEntryMappingNotDefined = errors.New(`entry mapping not defined`)
)
