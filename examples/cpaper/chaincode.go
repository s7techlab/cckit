package cpaper

import (
	"github.com/s7techlab/cckit/examples/cpaper/schema"
	"github.com/s7techlab/cckit/extensions/debug"
	"github.com/s7techlab/cckit/extensions/owner"
	"github.com/s7techlab/cckit/router"
	p "github.com/s7techlab/cckit/router/param"
	"github.com/s7techlab/cckit/state/mapping"
)

func NewCC() *router.Chaincode {

	r := router.New(`commercial_paper`)

	r.Use(mapping.MapState(mapping.SchemaMappings{}.
		Add(&schema.CommercialPaper{},
			[]string{`cpaper`},
			func(e interface{}) ([]string, error) {
				return []string{e.(*schema.CommercialPaper).GetPaper()}, nil
			})))

	r.Init(owner.InvokeSetFromCreator)

	debug.AddHandlers(r, `debug`, owner.Only)

	r.Group(`cpaper`).
		Invoke(`List`, cpaperList).
		//Invoke(`Get`, cpaperGet, p.String(`id`)).
		Invoke(`Insert`, cpaperInsert, p.Proto(`cpaper`, &schema.CommercialPaper{}))
	//Invoke(`Upsert`, cpaperUpsert, p.Struct(`cpaper`, &schema.CommercialPaper{})).
	//Invoke(`Delete`, cpaperDelete, p.String(`id`))

	return router.NewChaincode(r)
}

func cpaperList(c router.Context) (interface{}, error) {
	return c.State().List(&schema.CommercialPaper{})
}

func cpaperInsert(c router.Context) (interface{}, error) {
	cpaper := c.Param(`cpaper`)
	return cpaper, c.State().Insert(cpaper)
}

//func cpaperUpsert(c router.Context) (interface{}, error) {
//	cpaper := c.Param(`cpaper`)
//	return cpaper, c.State().Put(cpaper)
//}
//
//func cpaperGet(c router.Context) (interface{}, error) {
//	return c.State().Get(schema.cpaper{Id: c.ParamString(`id`)})
//}
//
//func cpaperDelete(c router.Context) (interface{}, error) {
//	return nil, c.State().Delete(schema.cpaper{Id: c.ParamString(`id`)})
//}
