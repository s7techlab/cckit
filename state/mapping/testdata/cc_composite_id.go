package testdata

import (
	"github.com/s7techlab/cckit/extensions/debug"
	"github.com/s7techlab/cckit/extensions/owner"
	"github.com/s7techlab/cckit/router"
	"github.com/s7techlab/cckit/router/param/defparam"
	"github.com/s7techlab/cckit/state/mapping"
	"github.com/s7techlab/cckit/state/mapping/testdata/schema"
)

var (
	EntityWithCompositeIdStateMapping = mapping.StateMappings{}.
		Add(&schema.EntityWithCompositeId{},
			mapping.PKeySchema(&schema.EntityCompositeId{}),
			mapping.List(&schema.EntityWithCompositeIdList{}))
)

func NewCompositeIdCC() *router.Chaincode {
	r := router.New("composite_id")

	r.Use(mapping.MapStates(EntityWithCompositeIdStateMapping))

	r.Use(mapping.MapEvents(mapping.EventMappings{}.
		Add(&schema.CreateEntityWithCompositeId{}).
		Add(&schema.UpdateEntityWithCompositeId{})))

	r.Init(owner.InvokeSetFromCreator)
	debug.AddHandlers(r, "debug", owner.Only)

	r.
		Query("list", queryListComposite).
		Query("get", queryByIdComposite, defparam.Proto(&schema.EntityCompositeId{})).
		Invoke("create", invokeCreateComposite, defparam.Proto(&schema.CreateEntityWithCompositeId{})).
		Invoke("update", invokeUpdateComposite, defparam.Proto(&schema.UpdateEntityWithCompositeId{})).
		Invoke("delete", invokeDeleteComposite, defparam.Proto(&schema.EntityCompositeId{}))

	return router.NewChaincode(r)
}

func queryByIdComposite(c router.Context) (interface{}, error) {
	return c.State().Get(c.Param().(*schema.EntityCompositeId))
}

func queryListComposite(c router.Context) (interface{}, error) {
	return c.State().List(&schema.EntityWithCompositeId{})
}

func invokeCreateComposite(c router.Context) (interface{}, error) {
	create := c.Param().(*schema.CreateEntityWithCompositeId)
	entity := &schema.EntityWithCompositeId{
		IdFirstPart:  create.IdFirstPart,
		IdSecondPart: create.IdSecondPart,
		Name:         create.Name,
		Value:        create.Value,
	}

	if err := c.Event().Set(create); err != nil {
		return nil, err
	}

	return entity, c.State().Insert(entity)
}

func invokeUpdateComposite(c router.Context) (interface{}, error) {
	update := c.Param().(*schema.UpdateEntityWithCompositeId)
	entity, _ := c.State().Get(
		&schema.EntityCompositeId{IdFirstPart: update.IdFirstPart, IdSecondPart: update.IdSecondPart},
		&schema.EntityWithCompositeId{})

	e := entity.(*schema.EntityWithCompositeId)

	e.Name = update.Name
	e.Value = update.Value

	if err := c.Event().Set(update); err != nil {
		return nil, err
	}

	return e, c.State().Put(e)
}

func invokeDeleteComposite(c router.Context) (interface{}, error) {
	return nil, c.State().Delete(c.Param().(*schema.EntityCompositeId))
}
