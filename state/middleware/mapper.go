package middleware

import (
	"github.com/s7techlab/cckit/router"
	"github.com/s7techlab/cckit/state"
)

func Mapper(mapping state.Mapping) router.MiddlewareFunc {
	return func(next router.HandlerFunc, pos ...int) router.HandlerFunc {
		return func(c router.Context) (interface{}, error) {
			c.State().Mapping(mapping)
			return next(c)
		}
	}
}
