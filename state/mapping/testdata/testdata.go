package testdata

import "github.com/s7techlab/cckit/state/mapping/testdata/schema"

var (
	ProtoIssueMocks = []schema.IssueProtoEntity{{
		IdFirstPart:  "A",
		IdSecondPart: "1",
		Name:         "Lorem",
		ExternalId:   "EXT1",
	}, {
		IdFirstPart:  "B",
		IdSecondPart: "1",
		Name:         "Ipsum",
		ExternalId:   "EXT2",
	}, {
		IdFirstPart:  "B",
		IdSecondPart: "2",
		Name:         "Dolor",
		ExternalId:   "EXT3",
	}}

	ProtoIssueMockExistingExternal = schema.IssueProtoEntity{
		IdFirstPart:  "Z",
		IdSecondPart: "1",
		Name:         "Lorem",
		ExternalId:   "EXT1",
	}

	ProtoIssueMockExistingPrimary = schema.IssueProtoEntity{
		IdFirstPart:  "A",
		IdSecondPart: "1",
		Name:         "Lorem",
		ExternalId:   "EXT100",
	}
)
