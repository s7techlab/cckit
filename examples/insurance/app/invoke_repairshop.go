package app

import (
	"encoding/json"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
)

func listRepairOrders(stub shim.ChaincodeStubInterface, _ []string) pb.Response {
	resultsIterator, err := stub.GetStateByPartialCompositeKey(prefixRepairOrder, []string{})
	if err != nil {
		return shim.Error(err.Error())
	}
	defer func() { _ = resultsIterator.Close() }()

	var results []interface{}
	for resultsIterator.HasNext() {
		kvResult, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		repairOrder := RepairOrder{}
		err = json.Unmarshal(kvResult.Value, &repairOrder)
		if err != nil {
			return shim.Error(err.Error())
		}
		if repairOrder.Ready {
			continue
		}

		result := struct {
			UUID         string `json:"uuid"`
			ClaimUUID    string `json:"claim_uuid"`
			ContractUUID string `json:"contract_uuid"`
			Item         Item   `json:"item"`
		}{}
		err = json.Unmarshal(kvResult.Value, &result)
		if err != nil {
			return shim.Error(err.Error())
		}
		prefix, keyParts, err := stub.SplitCompositeKey(kvResult.Key)
		if err != nil {
			return shim.Error(err.Error())
		}
		if len(keyParts) == 0 {
			result.UUID = prefix
		} else {
			result.UUID = keyParts[0]
		}
		results = append(results, result)
	}

	resultsAsBytes, err := json.Marshal(results)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(resultsAsBytes)
}

func completeRepairOrder(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Invalid argument count.")
	}

	input := struct {
		UUID string `json:"uuid"`
	}{}
	err := json.Unmarshal([]byte(args[0]), &input)
	if err != nil {
		return shim.Error(err.Error())
	}

	repairOrderKey, err := stub.CreateCompositeKey(prefixRepairOrder, []string{input.UUID})
	if err != nil {
		return shim.Error(err.Error())
	}

	repairOrderBytes, _ := stub.GetState(repairOrderKey)
	if len(repairOrderBytes) == 0 {
		return shim.Error("Could not find the repair order")
	}

	repairOrder := RepairOrder{}
	err = json.Unmarshal(repairOrderBytes, &repairOrder)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Marking repair order as ready
	repairOrder.Ready = true

	repairOrderBytes, err = json.Marshal(repairOrder)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(repairOrderKey, repairOrderBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Reflect changes in the corresponding claim
	claimKey, err := stub.CreateCompositeKey(prefixClaim, []string{repairOrder.ContractUUID, repairOrder.ClaimUUID})
	if err != nil {
		return shim.Error(err.Error())
	}
	claimBytes, _ := stub.GetState(claimKey)
	if claimBytes != nil {
		claim := Claim{}
		err := json.Unmarshal(claimBytes, &claim)
		if err != nil {
			return shim.Error(err.Error())
		}

		claim.Repaired = true
		claimBytes, err = json.Marshal(claim)
		if err != nil {
			return shim.Error(err.Error())
		}

		err = stub.PutState(claimKey, claimBytes)
		if err != nil {
			return shim.Error(err.Error())
		}
	}

	return shim.Success(nil)
}
