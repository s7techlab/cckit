package defparam

import (
	"github.com/s7techlab/cckit/router"
	"github.com/s7techlab/cckit/router/param"
)

func Proto(target interface{}, argPoss ...int) router.MiddlewareFunc {
	return param.Proto(router.DefaultParam, target, argPoss...)
}

func String(argPoss ...int) router.MiddlewareFunc {
	return param.String(router.DefaultParam, argPoss...)
}
