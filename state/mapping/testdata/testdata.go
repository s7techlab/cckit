package testdata

import "github.com/s7techlab/cckit/state/mapping/testdata/schema"

var (
	CreateEntityWithCompositeId = []*schema.CreateEntityWithCompositeId{{
		IdFirstPart:  "A",
		IdSecondPart: "1",
		Name:         "Lorem",
		Value:        1,
	}, {
		IdFirstPart:  "B",
		IdSecondPart: "1",
		Name:         "Ipsum",
		Value:        2,
	}, {
		IdFirstPart:  "B",
		IdSecondPart: "2",
		Name:         "Dolor",
		Value:        3,
	}}

	CreateEntityWithIndexes = []*schema.CreateEntityEntityWithIndexes{{
		Id:         `aaa`,
		ExternalId: `aaa_aaa`,
		Value:      1,
	}, {
		Id:                  `bbb`,
		ExternalId:          `bbb_bbb`,
		OptionalExternalIds: []string{`bbb_opt1`, `bbb_opt2`},
		Value:               1,
	}}
)
