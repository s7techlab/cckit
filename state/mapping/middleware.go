package mapping

import (
	"github.com/s7techlab/cckit/router"
)

func MapStates(stateMappings StateMappings) router.MiddlewareFunc {
	return func(next router.HandlerFunc, pos ...int) router.HandlerFunc {
		return func(c router.Context) (interface{}, error) {
			c.UseState(WrapState(c.State(), stateMappings))
			return next(c)
		}
	}
}

func MapEvents(eventMappings EventMappings) router.MiddlewareFunc {
	return func(next router.HandlerFunc, pos ...int) router.HandlerFunc {
		return func(c router.Context) (interface{}, error) {
			c.UseEvent(NewEvent(c.Stub(), eventMappings))
			return next(c)
		}
	}
}
