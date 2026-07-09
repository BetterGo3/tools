// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gcimporter_test

import (
	"go/ast"
	goparser "go/parser"
	"go/token"
	"go/types"
	"testing"

	"golang.org/x/tools/internal/gcimporter"
)

func TestIExportShallowOptional(t *testing.T) {
	const src = `package p

type Person struct {
	Title string?
}
`
	fset := token.NewFileSet()
	f, err := goparser.ParseFile(fset, "p.go", src, 0)
	if err != nil {
		t.Fatal(err)
	}
	pkg, err := new(types.Config).Check("p", fset, []*ast.File{f}, nil)
	if err != nil {
		t.Fatal(err)
	}

	data, err := gcimporter.IExportShallow(fset, pkg, nil)
	if err != nil {
		t.Fatalf("IExportShallow: %v", err)
	}

	imports := make(map[string]*types.Package)
	pkg2, err := gcimporter.IImportShallow(fset, gcimporter.GetPackagesFromMap(imports), data, "p", nil)
	if err != nil {
		t.Fatalf("IImportShallow: %v", err)
	}

	person := pkg2.Scope().Lookup("Person").Type().Underlying().(*types.Struct)
	title := person.Field(0).Type()
	opt, ok := title.(*types.Optional)
	if !ok {
		t.Fatalf("Title type = %T, want *types.Optional", title)
	}
	if opt.Elem() != types.Typ[types.String] {
		t.Fatalf("Elem() = %v, want string", opt.Elem())
	}
}
