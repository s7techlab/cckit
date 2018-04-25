package param

import (
	"github.com/s7techlab/cckit/router"
)

func String(name string, argPoss ...int) router.MiddlewareFunc {
	return Param(name, ``, argPoss...)
}

func Int(name string, argPoss ...int) router.MiddlewareFunc {
	return Param(name, 1, argPoss...)
}

func Bool(name string, argPoss ...int) router.MiddlewareFunc {
	return Param(name, true, argPoss...)
}

func Struct(name string, target interface{}, argPoss ...int) router.MiddlewareFunc {
	return Param(name, target, argPoss...)
}
