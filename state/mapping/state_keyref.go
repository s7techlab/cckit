package mapping

import (
	"strings"

	"github.com/s7techlab/cckit/state"
	"github.com/s7techlab/cckit/state/schema"
)

// KeyRefNamespace namespace for uniq indexes
const KeyRefNamespace = `_idx`

// KeyRefIDKeyer keyer for KeyRef entity
var KeyRefIDKeyer = attrsKeyer([]string{`Schema`, `Idx`, `RefKey`})

var KeyRefMapper = &StateMapping{
	schema:       &schema.KeyRef{},
	namespace:    state.Key{KeyRefNamespace},
	primaryKeyer: KeyRefIDKeyer,
}

var KeyRefIDMapper = &StateMapping{
	schema:       &schema.KeyRefId{},
	namespace:    state.Key{KeyRefNamespace},
	primaryKeyer: KeyRefIDKeyer,
}

func NewKeyRef(target interface{}, idx string, refKey, pKey state.Key) *schema.KeyRef {
	return &schema.KeyRef{
		Schema: strings.Join(SchemaNamespace(target), `-`),
		Idx:    idx,
		RefKey: []string(refKey),
		PKey:   []string(pKey),
	}
}

func NewKeyRefID(target interface{}, idx string, refKey state.Key) *schema.KeyRefId {
	return &schema.KeyRefId{
		Schema: strings.Join(SchemaNamespace(target), `-`),
		Idx:    idx,
		RefKey: []string(refKey),
	}
}

func NewKeyRefInstance(target interface{}, idx string, refKey, pKey state.Key) *StateInstance {
	return NewStateInstance(NewKeyRef(target, idx, refKey, pKey), KeyRefMapper, DefaultSerializer)
}

func NewKeyRefIDInstance(target interface{}, idx string, refKey state.Key) *StateInstance {
	return NewStateInstance(
		NewKeyRefID(target, idx, refKey),
		KeyRefIDMapper,
		DefaultSerializer,
	)
}
