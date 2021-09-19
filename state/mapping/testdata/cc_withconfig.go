package testdata

import (
	"github.com/s7techlab/cckit/extensions/owner"
	"github.com/s7techlab/cckit/router"
	"github.com/s7techlab/cckit/router/param/defparam"
	m "github.com/s7techlab/cckit/state/mapping"
	"github.com/s7techlab/cckit/state/mapping/testdata/schema"
)

func NewCCWithConfig() *router.Chaincode {
	r := router.New(`withConfig`)

	// Mappings for chaincode state
	r.Use(m.MapStates(m.StateMappings{}.
		//key will be <`config`>
		Add(&schema.Config{}, m.WithConstPKey())))

	r.Init(owner.InvokeSetFromCreator)

	r.Group(`config`).
		Invoke(`Set`, configSet, defparam.Proto(&schema.Config{})).
		Invoke(`Get`, configGet)

	return router.NewChaincode(r)
}

func configGet(c router.Context) (interface{}, error) {
	return c.State().Get(&schema.Config{}, &schema.Config{})
}

func configSet(c router.Context) (interface{}, error) {
	conf := c.Param()
	return conf, c.State().Put(conf)
}
