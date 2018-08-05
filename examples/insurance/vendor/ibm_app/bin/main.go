package main

import (
	"fmt"
	"ibm_app"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var logger = shim.NewLogger("main")

func main() {
	logger.SetLevel(shim.LogInfo)

	err := shim.Start(new(ibm_app.SmartContract))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
