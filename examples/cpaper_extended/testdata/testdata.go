package testdata

import (
	"time"

	"github.com/s7techlab/cckit/testing"

	"github.com/golang/protobuf/ptypes"
	"github.com/s7techlab/cckit/examples/cpaper_extended/schema"
)

var (
	// CPapers commercial paper fixtures
	CPapers = []schema.IssueCommercialPaper{{
		Issuer:       `some-issuer-1`,
		PaperNumber:  `00000001`,
		IssueDate:    ptypes.TimestampNow(),
		MaturityDate: testing.MustProtoTimestamp(time.Now().AddDate(0, 1, 0)),
		FaceValue:    11111,
		ExternalId:   `some-ext-id-1`,
	}, {
		Issuer:       `some-issuer-2`,
		PaperNumber:  `00000002`,
		IssueDate:    ptypes.TimestampNow(),
		MaturityDate: testing.MustProtoTimestamp(time.Now().AddDate(0, 2, 0)),
		FaceValue:    22222,
	}, {
		Issuer:       `some-issuer-3`,
		PaperNumber:  `00000003`,
		IssueDate:    ptypes.TimestampNow(),
		MaturityDate: testing.MustProtoTimestamp(time.Now().AddDate(0, 3, 0)),
	}}
)
