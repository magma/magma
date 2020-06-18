// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package txgen

import (
	"errors"
	"go/types"
	"os"
	"path/filepath"
	ttemplates "text/template"

	"github.com/99designs/gqlgen/codegen"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/99designs/gqlgen/codegen/templates"
	"github.com/99designs/gqlgen/plugin"
)

type txgen struct {
	config.PackageConfig
}

// New returns a txgen plugin
func New(cfg config.PackageConfig) plugin.Plugin {
	if cfg.Package == "" {
		cfg.Package = "resolver"
	}
	if cfg.Filename == "" {
		cfg.Filename = "tx_generated.go"
	}
	return txgen{cfg}
}

func (txgen) Name() string {
	return "txgen"
}

func (t txgen) MutateConfig(*config.Config) error {
	err := os.Remove(t.Filename)
	if os.IsNotExist(err) {
		err = nil
	}
	return err
}

func (t txgen) GenerateCode(data *codegen.Data) error {
	var mutation *codegen.Object
	for _, object := range data.Objects {
		if object.Definition == data.Schema.Mutation {
			mutation = object
			break
		}
	}
	if mutation == nil {
		return errors.New("unable to find mutation object")
	}
	return templates.Render(templates.Options{
		PackageName: t.Package,
		Filename:    filepath.Join(t.Package, t.Filename),
		Data: &txgenData{
			Object: mutation,
			Type:   "txResolver",
		},
		Funcs: ttemplates.FuncMap{
			"ResultType": func(f *codegen.Field) string {
				result := templates.CurrentImports.LookupType(f.TypeReference.GO)
				if f.Object.Stream {
					result = "<-chan " + result
				}
				return result
			},
			"Package": func(f *codegen.Field) string {
				t := f.TypeReference
				for e := t.Elem(); e != nil; e = t.Elem() {
					t = e
				}
				if t, ok := t.GO.(*types.Named); ok {
					return t.Obj().Pkg().Name()
				}
				return ""
			},
		},
		GeneratedHeader: true,
		Packages:        data.Config.Packages,
	})
}

type txgenData struct {
	*codegen.Object
	Type string
}
