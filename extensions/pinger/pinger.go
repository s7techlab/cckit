package pinger

import (
	"encoding/json"
	"github.com/hyperledger/fabric/protos/peer"
	r "github.com/s7techlab/cckit/router"
	"time"
	"github.com/pkg/errors"
	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
)

const (
	KeyPrefix        = `PING`
	FuncPingConstant = `pingLocal`
	FuncPing         = `ping`
	FuncPings        = `pings`
)

type PingInfo struct {
	Time        time.Time
	InvokerCert []byte
}

func (pi PingInfo) ToBytes() []byte {
	marshalled, _ := json.Marshal(pi)
	return marshalled
}

func PingInfoFromBytes(marshalled []byte) (pingInfo PingInfo, e error) {
	pi := new(PingInfo)
	e = json.Unmarshal(marshalled, pi)
	return *pi, e
}

func PingsInfoFromBytes(marshalled []byte) (pingsInfo []PingInfo, e error) {
	pi := new([]PingInfo)
	e = json.Unmarshal(marshalled, pi)
	return *pi, e
}

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

	key, err := c.Stub().CreateCompositeKey(KeyPrefix, []string{id, t.String()})
	if err != nil {
		return c.Response().Error(err)
	}

	return c.Response().Create(`pinged`, c.State().Put(key, PingInfo{t, cert.Raw}))
	//cid.ClientIdentity()
}

func PingConstant(c r.Context) peer.Response {
	t, err := c.Time()
	if err != nil {
		return c.Response().Error(errors.Wrap(err, `failed to get tx timestamp`))
	}
	return c.Response().Success(t.String())
}

func Pings(c r.Context) peer.Response {
	return c.Response().Create(c.State().List(KeyPrefix, PingInfo{}))
}
