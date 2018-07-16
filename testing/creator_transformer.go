package testing

import (
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/msp"
	"github.com/s7techlab/cckit/identity"

	pmsp "github.com/hyperledger/fabric/protos/msp"
)

// TransformCreator transforms arbitrary tx creator (pmsp.SerializedIdentity etc)  to mspID string, certPEM []byte,
func TransformCreator(txCreator ...interface{}) (mspID string, certPEM []byte, err error) {
	if len(txCreator) == 1 {
		p := txCreator[0]
		switch p.(type) {

		case identity.CertIdentity:
			return p.(identity.CertIdentity).MspID, p.(identity.CertIdentity).GetPEM(), nil

		case *identity.CertIdentity:
			return p.(*identity.CertIdentity).MspID, p.(*identity.CertIdentity).GetPEM(), nil

		case pmsp.SerializedIdentity:
			return p.(pmsp.SerializedIdentity).Mspid, p.(pmsp.SerializedIdentity).IdBytes, nil

		case msp.SigningIdentity:

			serialized, err := p.(msp.SigningIdentity).Serialize()
			if err != nil {
				return ``, nil, err
			}

			sid := &pmsp.SerializedIdentity{}
			if err = proto.Unmarshal(serialized, sid); err != nil {
				return ``, nil, err
			}
			return sid.Mspid, sid.IdBytes, nil

		case [2]string:
			// array with 2 elements  - mspId and ca cert
			return p.([2]string)[0], []byte(p.([2]string)[1]), nil
		}
	} else if len(txCreator) == 2 {
		return txCreator[0].(string), txCreator[1].([]byte), nil
	}

	return ``, nil, ErrUnknownFromArgsType
}
