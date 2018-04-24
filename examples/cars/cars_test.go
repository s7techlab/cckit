package main

import (
	"testing"

	testcc "github.com/s7techlab/cckit/testing"
)

func TestCars(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cars Suite")
}

var _ = Describe(`Cars`, func() {

	var cc *testcc.MockStub

	BeforeSuite(func() {

		cc = testcc.NewMockStub(`cars`, New())

	})

})
