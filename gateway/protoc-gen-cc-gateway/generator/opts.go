package generator

import (
	"strings"
)

const (
	ParamPaths                        = `paths`
	ParamPathsSourceRelative          = `source_relative`
	ParamEmbedSwagger                 = `embed_swagger`
	ParamChaincodeMethodServicePrefix = `service_name_method_prefix`
	ParamServiceChaincodeResolver     = `service_resolver`
)

// Opts by default all opts are disabled
type Opts struct {
	PathsSourceRelative          bool
	EmbedSwagger                 bool // generate var with embed annotation to include generated swagger fie
	ChaincodeMethodServicePrefix bool // add prefix with service name to chaincode method
	ServiceChaincodeResolver     bool
}

func isOptEnabled(paramValue string) bool {
	if paramValue == `0` || paramValue == `false` {
		return false
	}

	return true
}

func OptsFromParams(params string) Opts {
	opts := Opts{}
	for _, param := range strings.Split(params, ",") {
		var value string
		if i := strings.Index(param, "="); i >= 0 {
			value = param[i+1:]
			param = param[0:i]
		}
		switch param {
		case ParamPaths:
			switch value {
			case ParamPathsSourceRelative:
				opts.PathsSourceRelative = true
			}

		case ParamEmbedSwagger:
			opts.EmbedSwagger = isOptEnabled(value)

		case ParamChaincodeMethodServicePrefix:
			opts.ChaincodeMethodServicePrefix = isOptEnabled(value)

		case ParamServiceChaincodeResolver:
			opts.ServiceChaincodeResolver = isOptEnabled(value)
		}

	}

	return opts
}
