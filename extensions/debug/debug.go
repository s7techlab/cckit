package debug

import (
	"github.com/hyperledger/fabric-chaincode-go/shim"
)

// DelStateByPrefixes deletes from state entries with matching key prefix
func DelStateByPrefixes(stub shim.ChaincodeStubInterface, prefixes []string) (map[string]int, error) {
	prefixMatches := make(map[string]int)
	for _, prefix := range prefixes {
		iter, err := stub.GetStateByPartialCompositeKey(prefix, []string{})
		prefixMatches[prefix] = 0

		if err != nil {
			return nil, err
		}
		for iter.HasNext() {
			v, err := iter.Next()
			if err != nil {
				return nil, err
			}

			if err := stub.DelState(v.Key); err != nil {
				return nil, err
			}
			prefixMatches[prefix]++
		}
		iter.Close()
	}

	return prefixMatches, nil
}
