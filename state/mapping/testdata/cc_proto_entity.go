package testdata

import (
	"github.com/s7techlab/cckit/extensions/debug"
	"github.com/s7techlab/cckit/extensions/owner"
	"github.com/s7techlab/cckit/router"
	"github.com/s7techlab/cckit/router/param"
	"github.com/s7techlab/cckit/router/param/defparam"
	"github.com/s7techlab/cckit/state/mapping"
	"github.com/s7techlab/cckit/state/mapping/testdata/schema"
)

var (
	ProtoStateMapping = mapping.StateMappings{}.
				Add(&schema.ProtoEntity{},
			mapping.PKeySchema(&schema.ProtoEntityId{}),
			mapping.List(&schema.ProtoEntityList{}),
			mapping.UniqKey("ExternalId"),
		)

	ProtoEventMapping = mapping.EventMappings{}.
				Add(&schema.IssueProtoEntity{}).
				Add(&schema.IncrementProtoEntity{})
)

func NewProtoCC() *router.Chaincode {
	r := router.New("proto_test")
	r.Use(mapping.MapStates(ProtoStateMapping))
	r.Use(mapping.MapEvents(ProtoEventMapping))
	r.Init(owner.InvokeSetFromCreator)
	debug.AddHandlers(r, "debug", owner.Only)

	r.
		Query("list", queryList).
		Query("get", queryById, defparam.Proto(&schema.ProtoEntityId{})).
		Query("getByExternalId", queryByExternalId, param.String("externalId")).
		Invoke("issue", invokeIssue, defparam.Proto(&schema.IssueProtoEntity{})).
		Invoke("increment", invokeIncrement, defparam.Proto(&schema.IncrementProtoEntity{})).
		Invoke("delete", invokeDelte, defparam.Proto(&schema.ProtoEntityId{}))

	return router.NewChaincode(r)
}

func queryById(c router.Context) (interface{}, error) {
	return c.State().Get(c.Param().(*schema.ProtoEntityId))
}

func queryByExternalId(c router.Context) (interface{}, error) {
	externalId := c.ParamString("externalId")
	return c.State().(mapping.MappedState).GetByUniqKey(&schema.ProtoEntity{}, "ExternalId", []string{externalId})
}

func queryList(c router.Context) (interface{}, error) {
	return c.State().List(&schema.ProtoEntity{})
}

func invokeIssue(c router.Context) (interface{}, error) {
	issueData := c.Param().(*schema.IssueProtoEntity)
	entity := &schema.ProtoEntity{
		IdFirstPart:  issueData.IdFirstPart,
		IdSecondPart: issueData.IdSecondPart,
		Name:         issueData.Name,
		Value:        0,
		ExternalId:   issueData.ExternalId,
	}

	if err := c.Event().Set(issueData); err != nil {
		return nil, err
	}

	return entity, c.State().Insert(entity)
}

func invokeIncrement(c router.Context) (interface{}, error) {
	incrementData := c.Param().(*schema.IncrementProtoEntity)
	entity, _ := c.State().Get(
		&schema.ProtoEntityId{IdFirstPart: incrementData.IdFirstPart, IdSecondPart: incrementData.IdSecondPart},
		&schema.ProtoEntity{})

	protoEntity := entity.(*schema.ProtoEntity)

	protoEntity.Value = protoEntity.Value + 1

	if err := c.Event().Set(incrementData); err != nil {
		return nil, err
	}

	return protoEntity, c.State().Put(protoEntity)
}

func invokeDelte(c router.Context) (interface{}, error) {
	return nil, c.State().Delete(c.Param().(*schema.ProtoEntityId))
}
