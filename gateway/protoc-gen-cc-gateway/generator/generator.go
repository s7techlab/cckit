package generator

import (
	"bytes"
	"fmt"
	"go/format"
	"log"
	"path"
	"path/filepath"
	"strings"

	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway/descriptor"
)

type Generator struct {
	reg                 *descriptor.Registry
	imports             []descriptor.GoPackage // common imports
	PathsSourceRelative bool
}

// New returns a new generator which generates handler wrappers.
func New(reg *descriptor.Registry) *Generator {
	return &Generator{
		reg: reg,
	}
}

func (g *Generator) Generate(targets []*descriptor.File) ([]*plugin.CodeGeneratorResponse_File, error) {
	var files []*plugin.CodeGeneratorResponse_File
	for _, file := range targets {
		if len(file.Services) == 0 {
			continue
		}

		if code, err := g.generateCC(file); err == nil {
			files = append(files, code)
		} else {
			return nil, err
		}
	}

	return files, nil
}

func (g *Generator) generateCC(file *descriptor.File) (*plugin.CodeGeneratorResponse_File, error) {
	code, err := g.getCCTemplate(file)
	if err != nil {
		return nil, err
	}

	formatted, err := format.Source([]byte(code))
	if err != nil {
		log.Printf("%v: %s", err, annotateString(code))
		return nil, err
	}

	name := filepath.Base(file.GetName())
	ext := filepath.Ext(name)
	base := strings.TrimSuffix(name, ext)

	basePath := file.GoPkg.Name
	if !g.PathsSourceRelative {
		basePath = file.GoPkg.Path
	}

	output := fmt.Sprintf(filepath.Join(basePath, "%s.pb.cc.go"), base)
	output = filepath.Clean(output)

	return &plugin.CodeGeneratorResponse_File{
		Name:    proto.String(output),
		Content: proto.String(string(formatted)),
	}, nil
}

func (g *Generator) getCCTemplate(f *descriptor.File) (string, error) {
	pkgSeen := make(map[string]bool)
	var imports []descriptor.GoPackage
	for _, pkg := range g.imports {
		pkgSeen[pkg.Path] = true
		imports = append(imports, pkg)
	}

	pkgs := [][]string{
		{"context", "context"},
		{"github.com/s7techlab/cckit/gateway", "cckit_gateway"},
		{"github.com/s7techlab/cckit/gateway/service", "cckit_ccservice"},
		{"github.com/s7techlab/cckit/router", "cckit_router"},
		{"github.com/s7techlab/cckit/router/param/defparam", "cckit_defparam"},
		{"github.com/s7techlab/cckit/router/param", "cckit_param"},
	}

	for _, pkg := range pkgs {
		pkgSeen[pkg[0]] = true
		imports = append(imports, g.newGoPackage(pkg[0], pkg[1]))
	}

	for _, svc := range f.Services {
		for _, m := range svc.Methods {
			checkedAppend := func(pkg descriptor.GoPackage) {
				// Add request type package to imports if needed
				if m.Options == nil || pkg == f.GoPkg || pkgSeen[pkg.Path] {
					return
				}
				pkgSeen[pkg.Path] = true

				// always generate alias for external packages, when types used in req/resp object
				//if pkg.Alias == "" {
				//	pkg.Alias = pkg.Name
				//	pkgSeen[pkg.Path] = false
				//}

				imports = append(imports, pkg)
			}

			checkedAppend(m.RequestType.File.GoPkg)
			checkedAppend(m.ResponseType.File.GoPkg)
		}
	}

	p := param{File: f, Imports: imports}
	return applyTemplate(p)
}

func (g *Generator) newGoPackage(pkgPath string, aalias ...string) descriptor.GoPackage {
	gopkg := descriptor.GoPackage{
		Path: pkgPath,
		Name: path.Base(pkgPath),
	}
	alias := gopkg.Name
	if len(aalias) > 0 {
		alias = aalias[0]
		gopkg.Alias = alias
	}

	reference := alias
	if reference == "" {
		reference = gopkg.Name
	}

	for i := 0; ; i++ {
		if err := g.reg.ReserveGoPackageAlias(alias, gopkg.Path); err == nil {
			break
		}
		alias = fmt.Sprintf("%s_%d", gopkg.Name, i)
		gopkg.Alias = alias
	}

	pkg[reference] = alias

	return gopkg
}

func applyTemplate(p param) (string, error) {
	w := bytes.NewBuffer(nil)
	if err := headerTemplate.Execute(w, p); err != nil {
		return "", err
	}

	if err := ccTemplate.Execute(w, p); err != nil {
		return "", err
	}

	if err := gatewayTemplate.Execute(w, p); err != nil {
		return "", err
	}

	return w.String(), nil
}

func annotateString(str string) string {
	strs := strings.Split(str, "\n")
	for pos := range strs {
		strs[pos] = fmt.Sprintf("%v: %v", pos, strs[pos])
	}
	return strings.Join(strs, "\n")
}
