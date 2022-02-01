package generator

import (
	"strings"
	"text/template"

	"github.com/golang/protobuf/protoc-gen-go/generator"
	"github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway/descriptor"
)

var (
	funcMap = template.FuncMap{
		"goTypeName":      goTypeName,
		"hasBindings":     hasBindings,
		"hasGetBinding":   hasGetBinding,
		"removeExtension": removeExtension,
	}
)

func hasBindings(service *descriptor.Service) bool {
	for _, m := range service.Methods {
		if len(m.Bindings) > 0 {
			return true
		}
	}
	return false
}

func hasGetBinding(method *descriptor.Method) bool {
	for _, b := range method.Bindings {
		if b.HTTPMethod == "GET" {
			return true
		}
	}
	return false
}

func removeExtension(fileName string) string {
	startPos := 0
	slashPos := strings.LastIndex(fileName, `/`)

	if slashPos != -1 {
		startPos = slashPos + 1
	}
	return fileName[startPos:strings.LastIndex(fileName, `.`)]
}

func goTypeName(s string) string {
	toks := strings.Split(s, ".")
	i := 0
	if len(toks) > 1 {
		i = 1
	}
	for pos := range toks[i:] {
		toks[pos+i] = generator.CamelCase(toks[pos+i])
	}
	return strings.Join(toks, ".")
}
