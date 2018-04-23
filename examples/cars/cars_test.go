package main

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	testcc "github.com/s7techlab/cckit/testing"
)

func TestRefueling(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cars Suite")
}

var _ = Describe(`Cars`, func() {

	var cc *testcc.MockStub

	BeforeSuite(func() {

		cc = testcc.NewMockStub(`cars`, New())

	})

})
