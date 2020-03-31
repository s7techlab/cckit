package testing

import (
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/msp"
	"github.com/s7techlab/cckit/identity"

	pmsp "github.com/hyperledger/fabric/protos/msp"
)

func CreatorFromSigningIdentity(creator msp.SigningIdentity) (mspID string, certPEM []byte, err error) {
	serialized, err := creator.Serialize()
	if err != nil {
		return ``, nil, err
	}

	sid := &pmsp.SerializedIdentity{}
	if err = proto.Unmarshal(serialized, sid); err != nil {
		return ``, nil, err
	}
	return sid.Mspid, sid.IdBytes, nil
}

// TransformCreator transforms arbitrary tx creator (pmsp.SerializedIdentity etc)  to mspID string, certPEM []byte,
func TransformCreator(txCreator ...interface{}) (mspID string, certPEM []byte, err error) {
	if len(txCreator) == 1 {
		switch c := txCreator[0].(type) {

		case identity.CertIdentity:
			return c.MspID, c.GetPEM(), nil

		case *identity.CertIdentity:
			return c.MspID, c.GetPEM(), nil

		case pmsp.SerializedIdentity:
			return c.Mspid, c.IdBytes, nil

		case IdentitySample:
			id := c.SigningIdentity()
			return CreatorFromSigningIdentity(id)

		case msp.SigningIdentity:
			return CreatorFromSigningIdentity(c)

		case [2]string:
			// array with 2 elements  - mspId and ca cert
			return c[0], []byte(c[1]), nil
		}
	} else if len(txCreator) == 2 {
		return txCreator[0].(string), txCreator[1].([]byte), nil
	}

	return ``, nil, ErrUnknownFromArgsType
}
