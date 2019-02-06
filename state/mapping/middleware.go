package mapping

import (
	"github.com/s7techlab/cckit/router"
)

func MapState(mappings Mappings) router.MiddlewareFunc {
	return func(next router.HandlerFunc, pos ...int) router.HandlerFunc {
		return func(c router.Context) (interface{}, error) {
			c.UseState(NewState(c.Stub(), mappings))
			return next(c)
		}
	}
}
