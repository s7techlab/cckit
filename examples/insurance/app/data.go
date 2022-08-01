package app

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
)

// ContractType consists of prefix + UUID of the contract type
type ContractType struct {
	ShopType        string  `json:"shop_type"`
	FormulaPerDay   string  `json:"formula_per_day"`
	MaxSumInsured   float32 `json:"max_sum_insured"`
	TheftInsured    bool    `json:"theft_insured"`
	Description     string  `json:"description"`
	Conditions      string  `json:"conditions"`
	Active          bool    `json:"active"`
	MinDurationDays int32   `json:"min_duration_days"`
	MaxDurationDays int32   `json:"max_duration_days"`
}

// Contract consists of prefix + username + UUID of the contract
type Contract struct {
	Username         string    `json:"username"`
	Item             Item      `json:"item"`
	StartDate        time.Time `json:"start_date"`
	EndDate          time.Time `json:"end_date"`
	Void             bool      `json:"void"`
	ContractTypeUUID string    `json:"contract_type_uuid"`
	ClaimIndex       []string  `json:"claim_index,omitempty"`
}

// Item not persisted on its own
type Item struct {
	ID          int32   `json:"id"`
	Brand       string  `json:"brand"`
	Model       string  `json:"model"`
	Price       float32 `json:"price"`
	Description string  `json:"description"`
	SerialNo    string  `json:"serial_no"`
}

// Claim consists of prefix + UUID of the contract + UUID of the claim
type Claim struct {
	ContractUUID  string      `json:"contract_uuid"`
	Date          time.Time   `json:"date"`
	Description   string      `json:"description"`
	IsTheft       bool        `json:"is_theft"`
	Status        ClaimStatus `json:"status"`
	Reimbursable  float32     `json:"reimbursable"`
	Repaired      bool        `json:"repaired"`
	FileReference string      `json:"file_reference"`
}

// ClaimStatus the claim status indicates how the claim should be treated
type ClaimStatus int8

const (
	// ClaimStatusUnknown the claims status is unknown
	ClaimStatusUnknown ClaimStatus = iota
	// ClaimStatusNew the claim is new
	ClaimStatusNew
	// ClaimStatusRejected the claim has been rejected (either by the insurer, or by authorities
	ClaimStatusRejected
	// ClaimStatusRepair the item is up for repairs, or has been repaired
	ClaimStatusRepair
	// ClaimStatusReimbursement the customer should be reimbursed, or has already been
	ClaimStatusReimbursement
	// ClaimStatusTheftConfirmed the theft of the item has been confirmed by authorities
	ClaimStatusTheftConfirmed
)

func (s *ClaimStatus) UnmarshalJSON(b []byte) error {
	var value string
	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}

	switch strings.ToUpper(value) {
	default:
		*s = ClaimStatusUnknown
	case "N":
		*s = ClaimStatusNew
	case "J":
		*s = ClaimStatusRejected
	case "R":
		*s = ClaimStatusRepair
	case "F":
		*s = ClaimStatusReimbursement
	case "P":
		*s = ClaimStatusTheftConfirmed
	}

	return nil
}

func (s *ClaimStatus) MarshalJSON() ([]byte, error) {
	var value string

	switch s {
	default:
		fallthrough
	case ClaimStatusUnknown:
		value = ""
	case ClaimStatusNew:
		value = "N"
	case ClaimStatusRejected:
		value = "J"
	case ClaimStatusRepair:
		value = "R"
	case ClaimStatusReimbursement:
		value = "F"
	case ClaimStatusTheftConfirmed:
		value = "P"
	}

	return json.Marshal(value)
}

// User consists of prefix + username
type User struct {
	Username      string   `json:"username"`
	Password      string   `json:"password"`
	FirstName     string   `json:"first_name"`
	LastName      string   `json:"last_name"`
	ContractIndex []string `json:"contracts"`
}

// RepairOrder consists of prefix + UUID fo the repair order
type RepairOrder struct {
	ClaimUUID    string `json:"claim_uuid"`
	ContractUUID string `json:"contract_uuid"`
	Item         Item   `json:"item"`
	Ready        bool   `json:"ready"`
}

func (u *User) Contacts(stub shim.ChaincodeStubInterface) []Contract {
	contracts := make([]Contract, 0)

	// for each contractID in user.ContractIndex
	for _, contractID := range u.ContractIndex {

		c := &Contract{}

		// get contract
		contractAsBytes, err := stub.GetState(contractID)
		if err != nil {
			//res := "Failed to get state for " + contractID
			return nil
		}

		// parse contract
		err = json.Unmarshal(contractAsBytes, c)
		if err != nil {
			//res := "Failed to parse contract"
			return nil
		}

		// append to the contracts array
		contracts = append(contracts, *c)
	}

	return contracts
}

func (c *Contract) Claims(stub shim.ChaincodeStubInterface) ([]Claim, error) {
	var claims []Claim

	for _, claimKey := range c.ClaimIndex {
		claim := Claim{}

		claimAsBytes, err := stub.GetState(claimKey)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(claimAsBytes, &claim)
		if err != nil {
			return nil, err
		}

		claims = append(claims, claim)
	}

	return claims, nil
}

func (c *Contract) User(stub shim.ChaincodeStubInterface) (*User, error) {
	user := &User{}

	if len(c.Username) == 0 {
		return nil, errors.New("invalid user name in contract")
	}

	userKey, err := stub.CreateCompositeKey(prefixUser, []string{c.Username})
	if err != nil {
		return nil, err
	}

	userAsBytes, err := stub.GetState(userKey)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(userAsBytes, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (c *Claim) Contract(stub shim.ChaincodeStubInterface) (*Contract, error) {
	if len(c.ContractUUID) == 0 {
		return nil, nil
	}

	resultsIterator, err := stub.GetStateByPartialCompositeKey(prefixContract, []string{})
	if err != nil {
		return nil, err
	}
	defer func() { _ = resultsIterator.Close() }()

	for resultsIterator.HasNext() {
		kvResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		_, keyParams, err := stub.SplitCompositeKey(kvResult.Key)
		if len(keyParams) != 2 {
			continue
		}

		if keyParams[1] == c.ContractUUID {
			contract := &Contract{}
			err := json.Unmarshal(kvResult.Value, contract)
			if err != nil {
				return nil, err
			}
			return contract, nil
		}
	}
	return nil, nil
}
