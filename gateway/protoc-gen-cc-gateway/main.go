package main

import (
	"flag"
	"io"
	"log"
	"os"
	"strings"

	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/grpc-ecosystem/grpc-gateway/codegenerator"
	"github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway/descriptor"
	"github.com/s7techlab/cckit/gateway/protoc-gen-cc-gateway/generator"
)

var (
	file = flag.String("file", "-", "where to load data from")
)

func main() {
	var err error
	flag.Parse()
	if err = flag.Lookup("logtostderr").Value.Set("true"); err != nil {
		log.Fatal(err)
	}

	reg := descriptor.NewRegistry()
	fs := os.Stdin
	if *file != "-" {
		if fs, err = os.Open(*file); err != nil {
			log.Fatal(err)
		}
	}
	req, err := codegenerator.ParseRequest(fs)
	if err != nil {
		log.Fatal(err)
	}

	if err = reg.Load(req); err != nil {
		emitError(err)
		return
	}

	g := generator.New(reg)

	for _, param := range strings.Split(req.GetParameter(), ",") {
		var value string
		if i := strings.Index(param, "="); i >= 0 {
			value = param[i+1:]
			param = param[0:i]
		}
		switch param {
		case "paths":
			switch value {
			case "source_relative":
				g.PathsSourceRelative = true
			}
		}
	}

	var (
		targets []*descriptor.File
		f       *descriptor.File
	)
	for _, target := range req.FileToGenerate {
		if f, err = reg.LookupFile(target); err != nil {
			log.Fatal(err)
		}
		targets = append(targets, f)
	}

	out, err := g.Generate(targets)
	if err != nil {
		emitError(err)
		return
	}
	emitFiles(os.Stdout, out)
}

func emitFiles(w io.Writer, out []*plugin.CodeGeneratorResponse_File) {
	emitResp(w, &plugin.CodeGeneratorResponse{File: out})
}

func emitError(err error) {
	emitResp(os.Stdout, &plugin.CodeGeneratorResponse{Error: proto.String(err.Error())})
}

func emitResp(out io.Writer, resp *plugin.CodeGeneratorResponse) {
	buf, err := proto.Marshal(resp)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := out.Write(buf); err != nil {
		log.Fatal(err)
	}
}
