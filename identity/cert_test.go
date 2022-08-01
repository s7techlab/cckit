package identity_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/s7techlab/cckit/identity"
	"github.com/s7techlab/cckit/identity/testdata"
)

func TestIdentity(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Router suite")
}

var (
	certA = testdata.Certificates[0].MustCertBytes()
	certB = testdata.Certificates[1].MustCertBytes()
)

var _ = Describe(`Cert`, func() {

	BeforeSuite(func() {

	})

	It(`Allow to compare certificate subject `, func() {

		certEq, err := identity.CertSubjEqual(certA, certB)

		Expect(certEq).To(BeFalse())
		Expect(err).NotTo(HaveOccurred())

		certEq, err = identity.CertSubjEqual(certA, certA)

		Expect(certEq).To(BeTrue())
		Expect(err).NotTo(HaveOccurred())
	})

})
