package insurance

import (
	"time"

	"github.com/s7techlab/cckit/examples/insurance/app"
)

var (
	// ContractType fixture
	ContractType1 = ContractTypeDTO{
		UUID: `12345`,
		ContractType: &app.ContractType{
			ShopType:        `shop-type-1`,
			FormulaPerDay:   `aaa`,
			MaxSumInsured:   12345,
			TheftInsured:    true,
			Description:     `some_description`,
			Conditions:      `some_conditions`,
			Active:          true,
			MinDurationDays: 1,
			MaxDurationDays: 5,
		}}

	ContractType2 = ContractTypeDTO{
		UUID: `7890`,
		ContractType: &app.ContractType{
			ShopType:        `shop-type-2`,
			FormulaPerDay:   `bbb`,
			MaxSumInsured:   777,
			TheftInsured:    false,
			Description:     `once more description`,
			Conditions:      `once more conditions`,
			Active:          false,
			MinDurationDays: 2,
			MaxDurationDays: 10,
		}}

	// Contract fixture
	Contract1 = CreateContractDTO{
		UUID:             `xxx-aaa-bbb`,
		ContractTypeUUID: `xxx-ddd-ccc`,
		Username:         `vitiko`,
		Password:         `Root123AsUsual`,
		FirstName:        `Victor`,
		LastName:         `Nosov`,
		Item: app.Item{
			ID:          1,
			Brand:       `NoName`,
			Model:       `Model-XYZ`,
			Price:       123.45,
			Description: `Coolest thing ever`,
			SerialNo:    `ooo-999-222`,
		},
		StartDate: time.Now(),
		EndDate:   time.Now(),
	}
)
