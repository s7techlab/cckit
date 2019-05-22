# Test Hyperledger Fabric golang chaincode: check blockchain state, events and permissions

Testing stage is a critical requirement for software quality assurance, doesn't matter is this  
web application or a smart contract. Tests must be fast enough to run on every commit to repository. 
[CCKit](https://github.com/s7techlab/cckit/), programming toolkit for developing and testing hyperledger fabric golang 
chaincodes, enhances the development experience with extended version of MockStub for chaincode testing. 


## Chaincode testing techniques 

### Chaincode (smart contract) testing 

A smart contract defines the different states of a business object and governs the processes that move the object between 
these different states. Smart contracts are important because they allow architects and smart contract developers to define the 
key business processes and data that are shared across the different organizations collaborating in a blockchain network.

Tests must ensure that chaincode works as expected:

* particular input payload leads to particular business object state change
* particular input payload leads to validation or other errors
* particular object state allow subset of state transitions 

Any software testing (chaincode or web application for example) may either be a manual or an automated process. Manual 
software testing is led by a team or individual who will manually operate a software product and ensure it behaves as expected.
In case of chaincode tests you can manually invoke chaincode via `peer` cli tools.

Automated software testing is the practice of instrumenting input and output correctness checks for individual units of code.
During automated testing, code are executed in a test environment with simulated input.

### Steps in chaincode development process

The job of a smart contract developer is to take an existing business process and express it as a smart contract in a programming language. Steps of chaincode development:
    
*    Define chaincode model - schema for state entries, input payload and events
*    Define chaincode interface
*    Implement chaincode instantiate method
*    Implement chaincode methods with business logic
*    Create tests

Test driven development (TDD) or Behavioral Driven Development, possibly, single way to develop smart contracts.

### Running chaincode

Deploying chaincode to blockchain network isn't the quickest thing in the world, there's a lot of time that can be saved 
with testing. Also, more importantly, since blockchain is immutable and supposed to be secure because the code is on the
network, we rather not leave flaws in our code. 
           
During chaincode development we can divide testing to multiple stage - fast stage, when testing only smart contract logic,
and more complicated stage, when we do integration testing with live blockchain network, multiple peers, 
deployed on-chain code (smart contracts) and off-chain application, that uses SDK to connect with
blockchain network peers. 

### Chaincode DEV mode

Deploying a Hyperledger Fabric blockchain network, chaincode installing and initializing is quite complicated to set up and
a long procedure. Time to re-install / upgrade the code of a smart contract can be reduced by using 
[chaincode dev mode](https://hyperledger-fabric.readthedocs.io/en/latest/peer-chaincode-devmode.html). Normally chaincodes 
are started and maintained by peer. In “dev” mode, chaincode is built and started by the user. 
This mode is useful during chaincode development phase for rapid code/build/run/debug cycle turnaround. However, the process 
of updating the code will still be slow.

### MockStub

The [shim](https://github.com/hyperledger/fabric/tree/master/core/chaincode/shim) package contains a 
[MockStub](https://github.com/hyperledger/fabric/blob/master/core/chaincode/shim/mockstub.go) implementation 
that wraps calls to a chaincode, simulating its behavior in the HLF peer environment. MockStub does not need to start 
multiple docker containers with peer, world state database, chaincodes and allows to get test results almost immediately.
MockStub essentially replaces the SDK and peer enviroment and allows to test chaincode without actually starting your network. 
It implements almost every function the actual stub does, but in memory.

![mockstub](../docs/img/mockstub-hlf-peer.png)

MockStub includes implementation for most of `shim.ChaincodeStubInterface` function, but in the current version 
of Hyperledger Fabric (1.4), the MockStub has not implemented some of the important methods such
as `GetCreator` or method for work with private state range, for example. Since chaincode would use `GetCreator`  method 
to get tx creator certificate for access control, it's critical to be able to stub this method in order to completely unit-test chaincode.

### CCKit MockStub

CCKit [testing](.) package contains:

* [MockStub](mockstub.go) with implemented `GetCreator`, `GetTransient` methods and event subscription feature
* Test [identity](identity.go) creation helpers
* Chaincode [expect](expect) helpers

## Testing in Go

Go has a built-in testing command called `go test` and a package testing which combine to give a minimal but complete 
testing experience. In our example we use [Ginkgo](https://github.com/onsi/ginkgo) - BDD-style Go testing framework, builds on Go’s testing package, 
and allows you to write expressive tests in an efficient and effective manner. It is best paired with the [Gomega](https://github.com/onsi/gomega) 
matcher library but is designed to be matcher-agnostic.

As with popular BDD frameworks in other languages, Ginkgo allows you to group tests in `Describe` and `Context` container blocks. 
Ginkgo provides the `It` and `Specify` blocks which can hold your assertions. It also comes with handy structural utilities
such as `BeforeSuite`, `AfterSuite` etc that allow you to separate test configuration from test creation, and improve code reuse.

Ginkgo also comes with support for writing asynchronous tests. This makes testing code that use channels with chaincode events 
as easy as testing synchronous code.

# `Commercial Paper` chaincode

## Scenario

Official hyperledger fabric documentation contain detailed 
[chaincode example](https://hyperledger-fabric.readthedocs.io/en/release-1.4/developapps/scenario.html) -  commercial paper 
smart contract that defines the valid states for commercial paper, and the transaction logic that transition 
a paper from one state to another. We will test [commercial paper extended chaincode example](../examples/cpaper_extended) 
based on CCKit with protobuf state, described in [article](https://medium.com/coinmonks/hyperledger-fabric-smart-contract-data-model-protobuf-to-chaincode-state-mapping-191cdcfa0b78).

We can represent the lifecycle of a commercial paper using a state transition diagram: commercial papers transition 
between issued, trading and redeemed states by means of the issue, buy and redeem transactions.

![state](../examples/cpaper/img/state.png)

## Requirements 

To produce tests first we need to define requirements to tested application. Let’s start by listing our requirements 
for commercial paper chaincode:

* It should allow the issuer to issue commercial paper
* It should allow the participant to buy commercial paper
* It should allow the owner to redeem commercial paper

Chaincode interface functions described in file [chaincode.go](../examples/cpaper_extended/chaincode.go), so we can see all possible operations with chaincode data.

```go 
    Query("list", queryCPapers).

    // Get method has 2 params - commercial paper primary key components
    Query("get", queryCPaper, defparam.Proto(&schema.CommercialPaperId{})).
    Query("getByExternalId", queryCPaperGetByExternalId, param.String("externalId")).

    // txn methods
    Invoke("issue", invokeCPaperIssue, defparam.Proto(&schema.IssueCommercialPaper{})).
    Invoke("buy", invokeCPaperBuy, defparam.Proto(&schema.BuyCommercialPaper{})).
    Invoke("redeem", invokeCPaperRedeem, defparam.Proto(&schema.RedeemCommercialPaper{})).
    Invoke("delete", invokeCPaperDelete, defparam.Proto(&schema.CommercialPaperId{}))

```

## Getting started

Before you begin, be sure to get `CCKit`:

`git clone git@github.com:s7techlab/cckit.git`

This will fetch and install the CCKit package with [examples](../examples). After that we need to install the dependencies 
using command:

`go mod vendor`

## Creating test suite 

### Test package 

To write a new test suite, create a file whose name ends _test.go that contains the TestXxx functions, in our case will 
be [cpaper_extended/chaincode_test.go](../examples/cpaper_extended/chaincode_test.go)

Using separate package with tests [cpaper_extended_test](../examples/cpaper_extended/chaincode_test.go) instead of 
[cpaper_extended](../examples/cpaper_extended/chaincode.go) allows us to respect
the encapsulation of the chaincode package: your tests will need to import chaincode and access it from the outside. 
You cannot fiddle around with the internals, instead you focus on the exposed chaincode interface.

### Import matchers and helpers

To get started, we need to import the `matcher` functionality from the Ginkgo testing package 
so we can use different comparison mechanisms like comparing response objects or status codes.

We import the `ginkgo` and `gomega` packages with the . so that we can use functions from these packages without the package prefix.
In short, this allows us to use `Describe` instead of `ginkgo.Describe`, and `Equal` instead of `gomega.Equal`.


### Bootstrap

The call to `RegisterFailHandler` registers a handler, the `Fail` function from the ginkgo package, with `Gomega`. 
This creates the coupling between `Ginkgo` and `Gomega`.

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

func TestCommercialPaper(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Commercial paper suite")
}

var _ = Describe(`Commercial paper`, func() {
	

}
```

### Test structure
This particular test specification can be written using `Ginkgo` as follows:

```go
var _ = Describe(`CommercialPaper`, func() {
	
		Describe("Commercial Paper lifecycle", func() {
			
			It("Allow issuer to issue new commercial paper", func() { ... }
			
			It("Allow issuer to get commercial paper by composite primary key", func() { ... }
			
		    It("Allow issuer to get commercial paper by unique key", func() { ... }
		    
		    It("Allow issuer to get a list of commercial papers", func() { ... }
		    
		    It("Allow buyer to buy commercial paper", func() { ... }
		    
		    It("Allow buyer to redeem commercial paper", func() { ... }
			
		    It("Allow issuer to delete commercial paper", func() { ... }
		}
}
```

### Implementing tests
    
Now we go in depth to see how to create test functions specifically for chaincode development using `MockStub` functionalities.

Tests suite starts with creating a new instance of chaincode, or we can also instantiate a new chaincode instance before every test spec. 
This depends on how and what we want to test. In this example we instantiate a global `commercial paper` chaincode  that can 
be used in multiple test specs.

####  Chaincode invocation results

All chaincode invocation (via SDK to blockchain peer or to MockStub) resulted as 
[peer.Response](https://github.com/hyperledger/fabric/blob/release-1.4/protos/peer/proposal_response.pb.go) structure:

```go
type Response struct {
	// A status code that should follow the HTTP status codes.
	Status int32 
	// A message associated with the response code.
	Message string 
	// A payload that can be used to include metadata with this response.
	Payload              []byte   
}
```

During tests we can check `Response` attribute:

* Status (error or success)
* Message string ( contains error description)
* Payload contents ( marshaled json or protobuf)


#### Testing helpers

`Testing` package [contains](expect/matcher.go) multiple helpers / wrappers on ginkgo `expect` functions.

Most frequently used helpers are:

* `ResponseOk` (*response* **peer.Response**) expects that peer response contains ok code (`200`)
* `ResponseError` (*response* **peer.Response**)   expects that peer response contains error code (`500`). Optionally
you can pass expected error substring
* `PayloadIs`(*response* **peer.Response**, target interface{}) 

### Init function

When the chaincode is initialised, we want to test the execution status and make sure everything was
 successful. The equal method of the `expect` functionality comes in handy to compare the status.

## Running test

To run the test suite you have to simply run the command in the repository where the test suite is located:

`go test`


## Conclusion

Chaincode MockStub is really useful as it allows a developer to test his chaincode without starting the network every time. 
This reduces development time as he can use a test driven development (TDD) approach where he doesn’t 
need to start the network (this takes +- 40-80 seconds depending on the specs of the computer).