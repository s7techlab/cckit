package param

import (
	"github.com/s7techlab/cckit/router"
)

func String(name string) router.MiddlewareFunc {
	return Param(name, ``)
}

func Int(name string) router.MiddlewareFunc {
	return Param(name, 1)
}

func Struct(name string, target interface{}) router.MiddlewareFunc {
	return Param(name, target)
}
