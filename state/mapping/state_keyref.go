package mapping

import (
	"github.com/s7techlab/cckit/state"
	"github.com/s7techlab/cckit/state/schema"
)

const KeyRefNamespace = `_idx`

var KeyRefIdKeyer = attrsPKeyer([]string{`Schema`, `Idx`, `RefKey`})

var KeyRefMapper = &StateMapping{
	schema:       &schema.KeyRef{},
	namespace:    state.Key{KeyRefNamespace},
	primaryKeyer: KeyRefIdKeyer,
}

var KeyRefIdMapper = &StateMapping{
	schema:       &schema.KeyRefId{},
	namespace:    state.Key{KeyRefNamespace},
	primaryKeyer: KeyRefIdKeyer,
}

func NewKeyRef(target interface{}, idx string, refKey, pKey state.Key) *schema.KeyRef {
	return &schema.KeyRef{
		Schema: mapKey(target),
		Idx:    idx,
		RefKey: []string(refKey),
		PKey:   []string(pKey),
	}
}

func NewKeyRefId(target interface{}, idx string, refKey state.Key) *schema.KeyRefId {
	return &schema.KeyRefId{
		Schema: mapKey(target),
		Idx:    idx,
		RefKey: []string(refKey),
	}
}

func NewKeyRefMapped(target interface{}, idx string, refKey, pKey state.Key) *ProtoStateMapped {
	return NewProtoStateMapped(NewKeyRef(target, idx, refKey, pKey), KeyRefMapper)
}

func NewKeyRefIdMapped(target interface{}, idx string, refKey state.Key) *ProtoStateMapped {
	return NewProtoStateMapped(NewKeyRefId(target, idx, refKey), KeyRefIdMapper)
}
