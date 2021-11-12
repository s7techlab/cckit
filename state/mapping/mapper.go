package mapping

import (
	"errors"
	"fmt"

	"github.com/s7techlab/cckit/state"
)

var (
	ErrEventPayloadEmpty = errors.New(`event payload empty`)
)

type (
	EntryMapper struct {
		Commands Commands
		Event    *Event
	}
	Commands []Command

	Command interface {
		Execute(state.State) error
		fmt.Stringer
	}

	Query interface {
		Query(state.State) error
	}

	CommandInsert struct {
		Entry interface{}
	}

	CommandPut struct {
		Entry interface{}
	}
)

func NewEntryMapper() *EntryMapper {
	return &EntryMapper{
		Event: &Event{},
	}
}

func (em *EntryMapper) Apply(state state.State, event state.Event) error {
	for _, cmd := range em.Commands {
		if err := cmd.Execute(state); err != nil {
			return fmt.Errorf(`execute command=%s: %w`, cmd, err)
		}
	}

	if em.Event != nil {
		if em.Event.Name != `` && em.Event.Payload == nil {
			return ErrEventPayloadEmpty
		}

		if em.Event.Payload != nil {
			if em.Event.Name != `` {
				return event.Set(em.Event.Name, em.Event.Payload)
			}
			return event.Set(em.Event.Payload)
		}
	}

	return nil
}

func (cc *Commands) Insert(entry interface{}) *Commands {
	*cc = append(*cc, &CommandInsert{
		Entry: entry,
	})
	return cc
}

func (cc *Commands) Put(entry interface{}) *Commands {
	*cc = append(*cc, &CommandPut{
		Entry: entry,
	})
	return cc
}

func (ci *CommandInsert) Execute(state state.State) error {
	return state.Insert(ci.Entry)
}

func (ci *CommandInsert) String() string {
	return fmt.Sprintf(`insert<%T>`, ci.Entry)
}

func (ci *CommandPut) Execute(state state.State) error {
	return state.Put(ci.Entry)
}

func (ci *CommandPut) String() string {
	return fmt.Sprintf(`put<%T>`, ci.Entry)
}
