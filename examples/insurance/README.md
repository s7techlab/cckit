# Tests for IBM Blockchain insurance app with CCKit

This example shows how to apply concepts of test-driven development to writing chaincode in Golang for Hyperledger Fabric 
using shim MockStub to unit-test your chaincode without having to deploy it in a Blockchain network

 
##  IBM blockchain insurance app  

Blockchain insurance application example https://developer.ibm.com/code/patterns/build-a-blockchain-insurance-app/  
shows how through distributed ledger and smart contracts blockchain can transform insurance processes:
 
>Blockchain presents a huge opportunity for the insurance industry. It offers the chance to innovate around the way
data is exchanged, claims are processed, and fraud is prevented. Blockchain can bring together developers from tech 
companies, regulators, and insurance companies to create a valuable new insurance management asset. 
 
![Architecture](vendor/ibm_app/images/arch-blockchain-insurance2.png)


### Source code 

https://github.com/IBM/build-blockchain-insurance-app

Chaincode from  https://github.com/IBM/build-blockchain-insurance-app/tree/master/web/chaincode/src/bcins is copied to  
the [vendor/ibm-app](vendor/ibm_app) directory for sample test creation.

## Why test with Mockstub

Creating blockchain network and deploying the chaincode(s) to it is quite cumbersome and slow process, especially if your 
code is constantly changing during developing process. 

The Hyperledger Fabric [shim package](https://github.com/hyperledger/fabric/tree/release-1.2/core/chaincode/shim) 
contains the [MockStub](https://github.com/hyperledger/fabric/blob/release-1.2/core/chaincode/shim/mockstub.go)
implementation that wraps chaincode and stub out the calls to 
[shim.ChaincodeStubInterface](https://github.com/hyperledger/fabric/blob/release-1.2/core/chaincode/shim/interfaces.go) 
in chaincode functions. Mockstub allows you to get the test results almost immediately.


MockStub contains implementation for most of `shim.ChaincodeStubInterface` function, but in the current version 
of Hyperledger Fabric (1.2), the MockStub has not implemented some of the important methods such
as `GetCreator`, for example. Since chaincode would use this method to get tx creator certificate
for access control, it's critical to be able to stub this method  in order to completely unit-test chaincode. 

`CCKit` contains extended [MockStub](../../testing/mockstub.go) with implementation of some of the unimplemented
methods and delegating existing ones to shim.MockStub.

## Tests for this application

The tests are located in [ibm_app_test](ibm_app_test) directory:

* [app_test.go](ibm_app_test/app_test.go) contains [Ginkgo](https://onsi.github.io/ginkgo/) based tests
* [dto.go](ibm_app_test/dto.go) contains named DTO (data transfer objects) from original chaincode source code
* [fixture.go](ibm_app_test/fixtures.go) contains fixtures for tests


## Getting started

Ginkgo hooks into Go’s existing testing infrastructure. This allows you to run a Ginkgo suite using go test.
We import the ginkgo and gomega packages into the test’s top-level namespace by performing a dot-import. More about
testing with Ginkgo you can read on [onsi.github.io/ginkgo/](https://onsi.github.io/ginkgo/)

Using separate [package with tests](ibm_app_test) instead of [chaincode package](vendor/ibm_app) allows us to respect
the encapsulation of the chaincode package: your tests will need to import chaincode and access it from the outside.

```go
package main

import (
	"fmt"
	"testing"

	"ibm_app"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	testcc "github.com/s7techlab/cckit/testing"
	expectcc "github.com/s7techlab/cckit/testing/expect"
)

func TestInsuranceApp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Insurance app suite")
}

var _ = Describe(`Insurance`, func() {

	//Create chaincode mock
	cc := testcc.NewMockStub(`insurance`, new(ibm_app.SmartContract))
	
	...
	}
}
```

 
## Chaincode interface 

Chaincode functions described in file [main.go](main.go):

```go
var bcFunctions = map[string]func(shim.ChaincodeStubInterface, []string) pb.Response{
	// Insurance Peer
	"contract_type_ls":         listContractTypes,
	"contract_type_create":     createContractType,
	"contract_type_set_active": setActiveContractType,
	"contract_ls":              listContracts,
	"claim_ls":                 listClaims,
	"claim_file":               fileClaim,
	"claim_process":            processClaim,
	"user_authenticate":        authUser,
	"user_get_info":            getUser,
	// Shop Peer
	"contract_create": createContract,
	"user_create":     createUser,
	// Repair Shop Peer
	"repair_order_ls":       listRepairOrders,
	"repair_order_complete": completeRepairOrder,
	// Police Peer
	"theft_claim_ls":      listTheftClaims,
	"theft_claim_process": processTheftClaim,
}
```

So we can use them for testing insurance chaincode functionality.