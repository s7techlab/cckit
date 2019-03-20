package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/s7techlab/cckit/examples/insurance/app"
)

var logger = shim.NewLogger("main")

func main() {
	logger.SetLevel(shim.LogInfo)

	err := shim.Start(new(app.SmartContract))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
