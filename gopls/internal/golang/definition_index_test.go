// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package golang

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"testing"
)

func TestIndexOperatorDefinition(t *testing.T) {
	const src = `package p

type Vec struct {
	data []float64
}

func [](v *Vec, i int) float64 {
	return v.data[i]
}

func []=(v *Vec, i int, x float64) {
	v.data[i] = x
}

func use(v *Vec) {
	_ = v[0]
	v[0] = 1
}
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "p.go", src, 0)
	if err != nil {
		t.Fatal(err)
	}

	info := &types.Info{
		Uses:               make(map[*ast.Ident]types.Object),
		Defs:               make(map[*ast.Ident]types.Object),
		Selections:         make(map[*ast.SelectorExpr]*types.Selection),
		IndexOperatorCalls: make(map[ast.Expr]*ast.CallExpr),
		IndexAssignCalls:   make(map[ast.Expr]*ast.CallExpr),
	}
	if _, err := types.Config{IgnoreFuncBodies: false}.Check("p", fset, []*ast.File{f}, info); err != nil {
		t.Fatalf("typecheck: %v", err)
	}

	var readIndex, writeIndex *ast.IndexExpr
	ast.Inspect(f, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.AssignStmt:
			if len(n.Lhs) == 1 {
				if ix, ok := n.Lhs[0].(*ast.IndexExpr); ok {
					writeIndex = ix
				}
			}
		case *ast.IndexExpr:
			if _, ok := info.IndexOperatorCalls[n]; ok {
				readIndex = n
			}
		}
		return true
	})
	if readIndex == nil {
		t.Fatal("missing overloaded read index")
	}
	if writeIndex == nil {
		t.Fatal("missing overloaded write index")
	}
	if info.IndexOperatorCalls[readIndex] == nil {
		t.Fatal("IndexOperatorCalls missing read")
	}
	if info.IndexAssignCalls[writeIndex] == nil {
		t.Fatal("IndexAssignCalls missing write")
	}

	readObj := objectForOperatorCall(info, info.IndexOperatorCalls[readIndex])
	writeObj := objectForOperatorCall(info, info.IndexAssignCalls[writeIndex])
	if readObj == nil || readObj.Name() != "[]" {
		t.Fatalf("read operator obj = %v, want []", readObj)
	}
	if writeObj == nil || writeObj.Name() != "[]=" {
		t.Fatalf("write operator obj = %v, want []=", writeObj)
	}
	if !cursorOnIndexBrackets(readIndex.Lbrack, readIndex.Lbrack+1, readIndex) {
		t.Fatal("bracket selection failed for read")
	}
	if !cursorOnIndexBrackets(writeIndex.Lbrack, writeIndex.Lbrack+1, writeIndex) {
		t.Fatal("bracket selection failed for write")
	}
}

func objectForOperatorCall(info *types.Info, call *ast.CallExpr) types.Object {
	if call == nil {
		return nil
	}
	switch fun := call.Fun.(type) {
	case *ast.Ident:
		return info.ObjectOf(fun)
	case *ast.SelectorExpr:
		return info.ObjectOf(fun.Sel)
	default:
		return nil
	}
}
