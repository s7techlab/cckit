// Package pinger contains structure and functions for checking chain code accessibility
package pinger

import (
	"time"

	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	r "github.com/s7techlab/cckit/router"
)

const (
	// PingsStatePrefix prefix for PingInfo composite key in chain code state
	PingKeyPrefix = `PING`
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

// Ping chain code func for put tx creator information to chain code state
// checks endorsement policy is working
func Ping(c r.Context) (interface{}, error) {
	pingInfo, err := FromContext(c)
	if err != nil {
		return nil, err
	}

	return pingInfo, c.State().Put([]string{PingKeyPrefix, pingInfo.InvokerID, pingInfo.Time.String()}, pingInfo)
}

func FromContext(c r.Context) (*PingInfo, error) {
	id, err := cid.GetID(c.Stub())
	if err != nil {
		return nil, err
	}
	cert, err := cid.GetX509Certificate(c.Stub())
	if err != nil {
		return nil, err
	}
	t, err := c.Time()
	if err != nil {
		return nil, err
	}
	return &PingInfo{Time: t, InvokerID: id, InvokerCert: cert.Raw}, nil
}

// PingConstant chain code func only checks that chain code is installed and instantiated
func PingConstant(c r.Context) (interface{}, error) {
	return FromContext(c)
}

// Pings chain code func returns pings from chain code state
func Pings(c r.Context) (interface{}, error) {
	return c.State().List(PingKeyPrefix, &PingInfo{})
}
