package testdata

import (
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/s7techlab/cckit/examples/cpaper/schema"
)

var maturityDate1, _ = ptypes.TimestampProto(time.Now().AddDate(0, 1, 0))

var CPapers = []schema.CommercialPaper{{
	Paper:        `00000001`,
	Issuer:       `some-issuer`,
	Owner:        `some-issuer`,
	IssueDate:    ptypes.TimestampNow(),
	MaturityDate: maturityDate1,
}}
