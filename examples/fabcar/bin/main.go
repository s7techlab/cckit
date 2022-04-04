package main

import (
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"

	"github.com/s7techlab/cckit/examples/fabcar"
)

func main() {
	cc, err := fabcar.New()
	if err != nil {
		fmt.Printf("error creating %s chaincode: %s", fabcar.ChaincodeName, err)
		return
	}

	if err = shim.Start(cc); err != nil {
		fmt.Printf("error starting %s chaincode: %s", fabcar.ChaincodeName, err)
	}
}
