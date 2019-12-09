package testdata

import (
	"github.com/s7techlab/cckit/extensions/owner"
	"github.com/s7techlab/cckit/router"
	"github.com/s7techlab/cckit/router/param/defparam"
	"github.com/s7techlab/cckit/state/mapping"
	"github.com/s7techlab/cckit/state/mapping/testdata/schema"
)

var (
	EntityWithIndexesStateMapping = mapping.StateMappings{}.
		Add(&schema.EntityWithIndexes{},
			mapping.PKeyId(),
			mapping.List(&schema.EntityWithIndexesList{}),
			mapping.UniqKey(`ExternalId`),
			mapping.WithIndex(&mapping.StateIndexDef{
				Name:     `OptionalExternalIds`,
				Required: false,
				Multi:    true,
			}))
)

func NewIndexesCC() *router.Chaincode {
	r := router.New("indexes")

	r.Use(mapping.MapStates(EntityWithIndexesStateMapping))

	r.Init(owner.InvokeSetFromCreator)

	r.
		Query("list", queryListIndexes).
		Query("get", queryByIdIndexes, defparam.String()).
		Query("getByExternalId", queryByExternalId, defparam.String()).
		Query("getByOptMultiExternalId", queryByOptMultiExternalId, defparam.String()).
		Invoke("create", invokeCreateIndexes, defparam.Proto(&schema.CreateEntityWithIndexes{})).
		Invoke("update", invokeUpdateIndexes, defparam.Proto(&schema.UpdateEntityWithIndexes{})).
		Invoke("delete", invokeDeleteIndexes, defparam.String())

	return router.NewChaincode(r)
}

func queryByIdIndexes(c router.Context) (interface{}, error) {
	return c.State().Get(&schema.EntityWithIndexes{Id: c.Param().(string)})
}

func queryListIndexes(c router.Context) (interface{}, error) {
	return c.State().List(&schema.EntityWithIndexes{})
}

func invokeCreateIndexes(c router.Context) (interface{}, error) {
	create := c.Param().(*schema.CreateEntityWithIndexes)
	entity := &schema.EntityWithIndexes{
		Id:                  create.Id,
		ExternalId:          create.ExternalId,
		RequiredExternalIds: create.RequiredExternalIds,
		OptionalExternalIds: create.OptionalExternalIds,
		Value:               create.Value,
	}

	return entity, c.State().Insert(entity)
}

func invokeUpdateIndexes(c router.Context) (interface{}, error) {
	update := c.Param().(*schema.UpdateEntityWithIndexes)
	entity := &schema.EntityWithIndexes{
		Id:                  update.Id,
		ExternalId:          update.ExternalId,
		RequiredExternalIds: update.RequiredExternalIds,
		OptionalExternalIds: update.OptionalExternalIds,
		Value:               update.Value,
	}

	return entity, c.State().Put(entity)
}

func invokeDeleteIndexes(c router.Context) (interface{}, error) {
	return nil, c.State().(mapping.MappedState).Delete(&schema.EntityWithIndexes{Id: c.Param().(string)})
}

func queryByExternalId(c router.Context) (interface{}, error) {
	externalId := c.Param().(string)
	return c.State().(mapping.MappedState).GetByKey(
		&schema.EntityWithIndexes{}, "ExternalId", []string{externalId})
}

func queryByOptMultiExternalId(c router.Context) (interface{}, error) {
	externalId := c.Param().(string)
	return c.State().(mapping.MappedState).GetByKey(
		&schema.EntityWithIndexes{}, "OptionalExternalIds", []string{externalId})
}
