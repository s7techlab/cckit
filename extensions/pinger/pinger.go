// Package pinger contains structure and functions for checking chain code accessibility
package pinger

import (
	"time"

	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/s7techlab/cckit/identity"
	r "github.com/s7techlab/cckit/router"
)

const (
	// PingsStatePrefix prefix for PingInfo composite key in chain code state
	PingKeyPrefix = `PING`
	// PingEvent event name
	PingEvent = `PING`
	// FuncPingConstant func name
	FuncPingConstant = `pingLocal`
	// FuncPing func name
	FuncPing = `ping`
	// FuncPings func name
	FuncPings = `pings`
)

// PingInfo stores time and certificate of ping tx creator
type PingInfo struct {
	Time        time.Time
	InvokerID   string
	InvokerCert []byte
}

func (p PingInfo) Key() ([]string, error) {
	return []string{PingKeyPrefix, p.InvokerID, p.Time.String()}, nil
}

// Ping chaincode func puts tx creator information into chaincode state
// can be used for checking endorsement policy is working
func Ping(c r.Context) (interface{}, error) {
	pingInfo, err := FromContext(c)
	if err != nil {
		return nil, err
	}

	c.SetEvent(PingEvent, pingInfo)
	return pingInfo, c.State().Put(pingInfo, pingInfo)
}

// FromContext create PingInfo struct with tx creator Id and certificate in PEM format
func FromContext(c r.Context) (*PingInfo, error) {
	id, err := cid.GetID(c.Stub())
	if err != nil {
		return nil, err
	}

	//take certificate from creator
	invoker, err := identity.FromStub(c.Stub())
	if err != nil {
		return nil, err
	}
	t, err := c.Time()
	if err != nil {
		return nil, err
	}
	return &PingInfo{Time: t, InvokerID: id, InvokerCert: invoker.GetPEM()}, nil
}

// PingConstant chaincode func returns invoker information
// can be used for testing that chain code is installed and instantiated
func PingConstant(c r.Context) (interface{}, error) {
	return FromContext(c)
}

// Pings chain code func returns pings from chain code state
func Pings(c r.Context) (interface{}, error) {
	return c.State().List(PingKeyPrefix, &PingInfo{})
}
