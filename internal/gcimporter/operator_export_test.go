// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gcimporter_test

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"golang.org/x/tools/internal/gcimporter"
	. "go/types"
)

// Shallow export data preserves named types but not operator overload metadata.
// gopls must type-check such packages from source instead of export cache.
func TestShallowExportDropsOperatorOverloads(t *testing.T) {
	const matrixSrc = `package matrix

struct Matrix {
	rows, cols int
}

func +(l, r *Matrix) *Matrix {
	return l
}
`
	const mainSrc = `package main

import "matrix"

func demo(a, b *matrix.Matrix) {
	_ = a + b
}
`
	fset := token.NewFileSet()
	matrixFile, err := parser.ParseFile(fset, "matrix/matrix.go", matrixSrc, 0)
	if err != nil {
		t.Fatal(err)
	}
	conf := &Config{GoVersion: "go1.27"}
	matrixPkg, err := conf.Check("matrix", fset, []*ast.File{matrixFile}, nil)
	if err != nil {
		t.Fatalf("matrix Check failed: %v", err)
	}

	data, err := gcimporter.IExportShallow(fset, matrixPkg, nil)
	if err != nil {
		t.Fatalf("IExportShallow failed: %v", err)
	}
	imports := make(map[string]*Package)
	imported, err := gcimporter.IImportShallow(fset, gcimporter.GetPackagesFromMap(imports), data, "matrix", nil)
	if err != nil {
		t.Fatalf("IImportShallow failed: %v", err)
	}

	mainFile, err := parser.ParseFile(fset, "main.go", mainSrc, 0)
	if err != nil {
		t.Fatal(err)
	}
	conf.Importer = matrixImporter(func(path string) (*Package, error) {
		if path == "matrix" {
			return imported, nil
		}
		return nil, fmt.Errorf("unknown import %q", path)
	})
	if _, err := conf.Check("main", fset, []*ast.File{mainFile}, nil); err == nil {
		t.Fatal("main Check succeeded, want operator overload error after shallow import")
	}
}

type matrixImporter func(path string) (*Package, error)

func (f matrixImporter) Import(path string) (*Package, error) { return f(path) }
