# Tests for IBM Blockchain insurance app chaincodes with CCKit

This example shows how to apply concepts of test-driven development to writing chaincode in Golang for Hyperledger Fabric 1.2
using shim MockStub to unit-test your chaincode without having to deploy it in a Blockchain network.
The use of described testing technique allowed us to find small bug in 
[IBM blockchain insurance app](https://github.com/IBM/build-blockchain-insurance-app/pull/44)

 
##  IBM blockchain insurance app  

Blockchain insurance application example https://developer.ibm.com/code/patterns/build-a-blockchain-insurance-app/  
shows how through distributed ledger and smart contracts blockchain can transform insurance processes:
 
>Blockchain presents a huge opportunity for the insurance industry. It offers the chance to innovate around the way
data is exchanged, claims are processed, and fraud is prevented. Blockchain can bring together developers from tech 
companies, regulators, and insurance companies to create a valuable new insurance management asset. 
 
![Architecture](app/images/arch-blockchain-insurance2.png)


### Source code 

https://github.com/IBM/build-blockchain-insurance-app

Chaincode from  https://github.com/IBM/build-blockchain-insurance-app/tree/master/web/chaincode/src/bcins is copied to  
the [app](app) directory for sample test creation.

## Test-driven development

Test-driven development (TDD) is a software development process that relies on the repetition of a very short development cycle:
requirements are turned into very specific test cases, then the software is improved to pass the new tests, 
only. This is opposed to software development that allows software to be added that is not proven to meet requirements.

When creating smart contracts on the Ethereum platform or chaincodes for Hyperledger Fabric, the developer programs a certain work logic 
that determines how the methods should change the state of the smart contract / chaincodes , what events are to be emitted,
when to return success or error response. This kind of software is ideal for development through testing.


## Why test with Mockstub

Creating blockchain network and deploying the chaincode(s) to it is quite cumbersome and slow process, especially if your 
code is constantly changing during developing process. 

The Hyperledger Fabric [shim package](https://github.com/hyperledger/fabric/tree/release-1.2/core/chaincode/shim) 
contains the [MockStub](https://github.com/hyperledger/fabric/blob/release-1.2/core/chaincode/shim/mockstub.go)
implementation that wraps chaincode and stub out the calls to 
[shim.ChaincodeStubInterface](https://github.com/hyperledger/fabric/blob/release-1.2/core/chaincode/shim/interfaces.go) 
in chaincode functions. Mockstub allows you to get the test results almost immediately.


MockStub contains implementation for most of `shim.ChaincodeStubInterface` function, but in the current version 
of Hyperledger Fabric (1.4), the MockStub has not implemented some of the important methods such
as `GetCreator`, for example. Since chaincode would use this method to get tx creator certificate
for access control, it's critical to be able to stub this method  in order to completely unit-test chaincode. 

`CCKit` contains extended [MockStub](../../testing/mockstub.go) with implementation of some of the unimplemented
methods and delegating existing ones to shim.MockStub. [CCkit MockStub]((../../testing/mockstub.go)) holds a reference
to the original MockStub and has additional methods and properties. 



## Tests for this application

The tests are located in [current](.) directory:

* [app_test.go](app_test.go) contains [Ginkgo](https://onsi.github.io/ginkgo/) based tests
* [dto.go](dto.go) contains named DTO (data transfer objects) from original chaincode source code
* [fixture.go](fixtures.go) contains fixtures for tests


## Getting started

Before you begin, be sure to get `CCkit`:

`git clone git@github.com:s7techlab/cckit.git`

and get dependencies using `go mod` command:

`go mod vendor`


CCkit repository contains several examples:

* [Cars](../cars)
* [Marbles](../marbles)
* And `insurance`

All examples and this tutorial uses the [Ginkgo](https://onsi.github.io/ginkgo/) testing library, it hooks into Go’s existing testing infrastructure. This allows you to run a Ginkgo suite using go test.
We import the ginkgo and gomega packages into the test’s top-level namespace by performing a dot-import. More about
testing with Ginkgo you can read on [onsi.github.io/ginkgo/](https://onsi.github.io/ginkgo/)


## Creating test suite 

Using separate [package with tests](.) instead of [chaincode package](app) allows us to respect
the encapsulation of the chaincode package: your tests will need to import chaincode and access it from the outside.

Test suite bootstrap example:

```go
package main

import (
	"fmt"
	"testing"
	"github.com/s7techlab/cckit/examples/insurance/app"
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
	cc := testcc.NewMockStub(`insurance`, new(app.SmartContract))
	
	...
	}
}
```

 
## What to test 

Chaincode interface functions described in file [main.go](main.go), so we can see all possible operations with chaincode data:

* list contract types
* create contract type
* activate contract type
* list contracts
* list claims...

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
Also chaincode has `init` function. So we can use them for testing insurance chaincode functionality.


## Data transfer objects and fixtures

In this example we use [Data Transfer Objects](dto.go) (DTO) for creating chaincode args and
[fixtures](fixtures.go)


## Test init function

```go
Describe("Chaincode initialization ", func() {
		It("Allow to provide contract types attributes  during chaincode creation [init]", func() {
			expectcc.ResponseOk(cc.Init(`init`, &ContractTypesDTO{ContractType1}))
		})
	})
````

## Test operations with contract types

```go
	Describe("Contract type ", func() {

		It("Allow to retrieve all contract type, added during chaincode init [contract_type_ls]", func() {
			contractTypes := expectcc.PayloadIs(
				cc.Invoke(`contract_type_ls`), &ContractTypesDTO{}).(ContractTypesDTO)

			Expect(len(contractTypes)).To(Equal(1))
			Expect(contractTypes[0].ShopType).To(Equal(ContractType1.ShopType))
		})

		It("Allow to create new contract type [contract_type_create]", func() {
			expectcc.ResponseOk(cc.Invoke(`contract_type_create`, &ContractType2))

			// get contract type list
			contractTypes := expectcc.PayloadIs(
				cc.Invoke(`contract_type_ls`), &ContractTypesDTO{}).(ContractTypesDTO)

			//expect now we have 2 contract type
			Expect(len(contractTypes)).To(Equal(2))
		})

		It("Allow to set active contract type [contract_type_set_active]", func() {
			Expect(ContractType2.Active).To(BeFalse())

			// active ContractType2
			expectcc.ResponseOk(cc.Invoke(`contract_type_set_active`, &ContractTypeActiveDTO{
				UUID: ContractType2.UUID, Active: true}))
		})

		It("Allow to retrieve filtered by shop type contract types [contract_type_ls]", func() {
			contractTypes := expectcc.PayloadIs(
				cc.Invoke(`contract_type_ls`, &ShopTypeDTO{ContractType2.ShopType}),
				&ContractTypesDTO{}).(ContractTypesDTO)

			Expect(len(contractTypes)).To(Equal(1))
			Expect(contractTypes[0].UUID).To(Equal(ContractType2.UUID))

			// Contract type 2 activated on previous step
			Expect(contractTypes[0].Active).To(BeTrue())
		})
	})
```

## Test operations with contracts

```go
Describe("Contract", func() {

		It("Allow everyone to create contract", func() {

			// get chaincode invoke response payload and expect returned payload is serialized instance of some structure
			contractCreateResponse := expectcc.PayloadIs(
				cc.Invoke(
					// invoke chaincode function
					`contract_create`,
					// with ContractDTO payload, it will automatically will be marshalled to json
					&Contract1,
				),
				// We expect than payload response is marshalled ContractCreateResponse structure
				&ContractCreateResponse{}).(ContractCreateResponse)

			Expect(contractCreateResponse.Username).To(Equal(Contract1.Username))
		})

		It("Allow every one to get user info", func() {

			// orininally was error https://github.com/IBM/build-blockchain-insurance-app/pull/44
			user := expectcc.PayloadIs(
				cc.Invoke(`user_get_info`, &GetUserDTO{
					Username: Contract1.Username,
				}), &ResponseUserDTO{}).(ResponseUserDTO)

			Expect(user.LastName).To(Equal(Contract1.LastName))
		})

	})
```