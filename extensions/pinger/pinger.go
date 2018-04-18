// Package pinger contains structure and functions for checking chain code accessibility
package pinger

import (
	"encoding/json"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/pkg/errors"
	r "github.com/s7techlab/cckit/router"
)

const (
	// PingsStatePrefix prefix for PingInfo composite key in chain code state
	PingsStatePrefix = `PING`
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
	InvokerCert []byte
}

// ToBytes marshals PingInfo struct to json bytes
func (pi PingInfo) ToBytes() []byte {
	marshalled, _ := json.Marshal(pi)
	return marshalled
}

// PingInfoFromBytes unmarshal from bytes
func PingInfoFromBytes(marshalled []byte) (pingInfo PingInfo, e error) {
	pi := new(PingInfo)
	e = json.Unmarshal(marshalled, pi)
	return *pi, e
}

// PingsInfoFromBytes unmarshal from bytes
func PingsInfoFromBytes(marshalled []byte) (pingsInfo []PingInfo, e error) {
	pi := new([]PingInfo)
	e = json.Unmarshal(marshalled, pi)
	return *pi, e
}

// Ping chain code func for put tx creator information to chain code state
// checks endorsement policy is working
func Ping(c r.Context) peer.Response {
	id, err := cid.GetID(c.Stub())
	if err != nil {
		return c.Response().Error(err)
	}

	cert, err := cid.GetX509Certificate(c.Stub())
	if err != nil {
		return c.Response().Error(err)
	}

	t, err := c.Time()
	if err != nil {
		return c.Response().Error(err)
	}

	key, err := c.Stub().CreateCompositeKey(PingsStatePrefix, []string{id, t.String()})
	if err != nil {
		return c.Response().Error(err)
	}

	return c.Response().Create(`pinged`, c.State().Put(key, PingInfo{t, cert.Raw}))
	//cid.ClientIdentity()
}

// PingConstant chain code func only checks that chain code is installed and instantiated
func PingConstant(c r.Context) peer.Response {
	t, err := c.Time()
	if err != nil {
		return c.Response().Error(errors.Wrap(err, `failed to get tx timestamp`))
	}
	return c.Response().Success(t.String())
}

// Pings chain code func returns pings from chain code state
func Pings(c r.Context) peer.Response {
	return c.Response().Create(c.State().List(PingsStatePrefix, PingInfo{}))
}
