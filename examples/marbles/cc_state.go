package marbles

import (
	"github.com/s7techlab/cckit/examples/marbles/schema"
	"github.com/s7techlab/cckit/state"
)

const MarbleEntity = `marble`

var Mappings = state.Mapping{
	Type: state.EntryMappings{
		`*schema.Marble`: {
			PrimaryKey: func(e interface{}) ([]string, error) {
				return []string{MarbleEntity, e.(*schema.Marble).GetID()}, nil
			}},
	}}
