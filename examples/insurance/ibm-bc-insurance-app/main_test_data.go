package main

import "time"

// ContractDTO is used to create contract
type ContractDTO struct {
	UUID             string    `json:"uuid"`
	ContractTypeUUID string    `json:"contract_type_uuid"`
	Username         string    `json:"username"`
	Password         string    `json:"password"`
	FirstName        string    `json:"first_name"`
	LastName         string    `json:"last_name"`
	Item             item      `json:"item"`
	StartDate        time.Time `json:"start_date"`
	EndDate          time.Time `json:"end_date"`
}

type ContractCreateResponse struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
