package events

import (
	"io"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/s7techlab/cckit/gateway"
)

//go:generate make generate

type ChaincodeEventGateway struct {
	gateway.Chaincode
}

// ApiDef returns service definition
func (c *ChaincodeEventGateway) ApiDef() gateway.ServiceDef {
	return gateway.ServiceDef{
		Desc:                        &_ChaincodeEvent_serviceDesc,
		Service:                     c,
		HandlerFromEndpointRegister: RegisterChaincodeEventHandlerFromEndpoint,
	}
}

func (s *ChaincodeEventGateway) EventStream(_ *empty.Empty, stream ChaincodeEvent_EventStreamServer) error {
	sub, err := s.Chaincode.Events(stream.Context())
	if err != nil {
		return err
	}

	defer sub.Close()
	for ev := range sub.Events() {
		if ev != nil {
			errS := stream.Send(ev)
			if errS == io.EOF {
				return nil
			}
		}
	}

	return nil
}
