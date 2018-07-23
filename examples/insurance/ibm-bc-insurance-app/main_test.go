package main

import (
	"testing"

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
	cc := testcc.NewMockStub(`cars`, new(SmartContract))

	BeforeSuite(func() {
		// init chaincode
		expectcc.ResponseOk(cc.Init()) // init chaincode
	})

})
