package sdk

// SDK interface for deal with Hyperledger Fabric SDK
// client from github.com/s7techlab/hlf-sdk-go implements this interface
type SDK interface {
	Invoker
	EventDelivery
}
