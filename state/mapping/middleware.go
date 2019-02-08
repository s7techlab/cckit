package mapping

import (
	"github.com/s7techlab/cckit/router"
)

func MapStates(stateMappings StateMappings) router.MiddlewareFunc {
	return func(next router.HandlerFunc, pos ...int) router.HandlerFunc {
		return func(c router.Context) (interface{}, error) {
			c.UseState(NewState(c.Stub(), stateMappings))
			return next(c)
		}
	}
}

func MapEvents(events Events) router.MiddlewareFunc {
	return func(next router.HandlerFunc, pos ...int) router.HandlerFunc {
		return func(c router.Context) (interface{}, error) {
			// not really mapped yet
			return next(c)
		}
	}
}
