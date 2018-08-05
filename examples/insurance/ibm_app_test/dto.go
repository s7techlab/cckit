package main

import (
	"ibm_app"
	"time"
)

// ContractTypesDTO type used in in "Init" func (arg) in main.go and in "listContractTypes" (return) in "invoke_insurance.go"
type ContractTypesDTO []ContractTypeDTO

type ContractTypeDTO struct {
	UUID string `json:"uuid"`
	*ibm_app.ContractType
}

// CreateContractDTO type used in in "createContract" func (arg) in invoke_shop.go
type CreateContractDTO struct {
	UUID             string       `json:"uuid"`
	ContractTypeUUID string       `json:"contract_type_uuid"`
	Username         string       `json:"username"`
	Password         string       `json:"password"`
	FirstName        string       `json:"first_name"`
	LastName         string       `json:"last_name"`
	Item             ibm_app.Item `json:"item"`
	StartDate        time.Time    `json:"start_date"`
	EndDate          time.Time    `json:"end_date"`
}

// ContractCreateResponse type used in in "createContract" func in invoke_shop.go
type ContractCreateResponse struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LsContractTypeDTO type used in "listContractTypes" in invoke_insurance.go
type LsContractTypeDTO struct {
	ShopType string `json:"shop_type"`
}

// type used in in "getUser" func in invoke_insurance.go
type GetUserDTO struct {
	Username string `json:"username"`
}

// named type, from anonymous type used in in "getUser" func in invoke_insurance.go
type ResponseUserDTO struct {
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
