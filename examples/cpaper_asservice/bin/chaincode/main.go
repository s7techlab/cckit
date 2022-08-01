package main

import (
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"

	"github.com/s7techlab/cckit/examples/cpaper_asservice"
)

func main() {
	cc, err := cpaper_asservice.NewCC()
	if err != nil {
		fmt.Printf("Error creating chaincode: %s", err)
	}

	err = shim.Start(cc)
	if err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
}
