package param

import (
	"github.com/s7techlab/cckit/router"
)

// StrictKnown allows passing arguments to chaincode func only if parameters are defined in router
func StrictKnown(next router.HandlerFunc, _ ...int) router.HandlerFunc {
	return func(c router.Context) (interface{}, error) {
		return next(c)
	}
}
