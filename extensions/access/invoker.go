package access

import (
	"crypto/x509"
	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/s7techlab/cckit/identity"
)

type Invoker struct {
	mspId string
	cert  *x509.Certificate
}

func (i Invoker) GetId() string {
	return identity.IdByCert(i.cert)
}

func (i Invoker) GetMSPId() string {
	return i.mspId
}

func (i Invoker) GetSubject() string {
	return identity.GetDN(&i.cert.Subject)
}

func (i Invoker) GetIssuer() string {
	return identity.GetDN(&i.cert.Issuer)
}

func (i Invoker) Is(id identity.Identity) bool {
	return i.mspId == id.GetMSPId() && i.GetSubject() == id.GetSubject()
}

func InvokerFromCert(mspId string, c []byte) (i identity.Identity, err error) {
	cert, err := identity.Certificate(c)
	if err != nil {
		return nil, err
	}

	return &Invoker{mspId, cert}, nil
}

func InvokerFromStub(stub shim.ChaincodeStubInterface) (i identity.Identity, err error) {
	clientIdentity, err := cid.New(stub)
	if err != nil {
		return
	}

	mspId, err := clientIdentity.GetMSPID()
	if err != nil {
		return
	}
	cert, err := clientIdentity.GetX509Certificate()
	if err != nil {
		return
	}

	return &Invoker{mspId, cert}, nil
}
