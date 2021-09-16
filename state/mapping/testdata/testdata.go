package testdata

import (
	"github.com/s7techlab/cckit/state/mapping/testdata/schema"
	"github.com/s7techlab/cckit/testing"
)

var (
	Dates = []string{`2021-02-15`, `2021-02-16`, `2021-03-15`}

	CreateEntityWithComplextId = []*schema.EntityWithComplexId{{
		Id: &schema.EntityComplexId{
			IdPart1: []string{`aaa`, `bb`},
			IdPart2: `ccc`,
			IdPart3: testing.MustTime(`2020-01-28T17:00:00Z`),
		},
	}}

	CreateEntityWithCompositeId = []*schema.CreateEntityWithCompositeId{{
		IdFirstPart:  "A",
		IdSecondPart: "1",
		IdThirdPart:  testing.MustTime(Dates[0] + `T00:00:00Z`),
		Name:         "Lorem",
		Value:        1,
	}, {
		IdFirstPart:  "B",
		IdSecondPart: "1",
		IdThirdPart:  testing.MustTime(Dates[1] + `T00:00:00Z`),
		Name:         "Ipsum",
		Value:        2,
	}, {
		IdFirstPart:  "B",
		IdSecondPart: "2",
		IdThirdPart:  testing.MustTime(Dates[2] + `T00:00:00Z`),
		Name:         "Dolor",
		Value:        3,
	}}

	CreateEntityWithIndexes = []*schema.CreateEntityWithIndexes{{
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
