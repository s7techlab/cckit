package testdata

import (
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/s7techlab/cckit/examples/cpaper/schema"
)

var (
	maturityDate1, _ = ptypes.TimestampProto(time.Now().AddDate(0, 1, 0))
	maturityDate2, _ = ptypes.TimestampProto(time.Now().AddDate(0, 2, 0))
	maturityDate3, _ = ptypes.TimestampProto(time.Now().AddDate(0, 3, 0))

	// CPapers commercial paper fixtures
	CPapers = []schema.IssueCommercialPaper{{
		Issuer:       `some-issuer-1`,
		PaperNumber:  `00000001`,
		IssueDate:    ptypes.TimestampNow(),
		MaturityDate: maturityDate1,
		FaceValue:    11111,
	}, {
		Issuer:       `some-issuer-2`,
		PaperNumber:  `00000002`,
		IssueDate:    ptypes.TimestampNow(),
		MaturityDate: maturityDate2,
		FaceValue:    22222,
	}, {
		Issuer:       `some-issuer-3`,
		PaperNumber:  `00000003`,
		IssueDate:    ptypes.TimestampNow(),
		MaturityDate: maturityDate3,
	}}
)
