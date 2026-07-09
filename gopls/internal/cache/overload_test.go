// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cache

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"strings"
	"testing"
)

func typecheckOverload(t *testing.T, src string, ignoreBodies bool) error {
	t.Helper()
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "p.go", src, 0)
	if err != nil {
		t.Fatal(err)
	}
	cfg := &types.Config{IgnoreFuncBodies: ignoreBodies}
	_, err = cfg.Check("p", fset, []*ast.File{f}, nil)
	return err
}

func TestOverloadTypecheck(t *testing.T) {
	const src = `package p

func add(int) int { return 0 }
func add(int64) int64 { return 0 }

type S struct{}
func (S) m(int) int { return 0 }
func (S) m(int64) int64 { return 0 }

func use() {
	_ = add(1)
	_ = add(int64(2))
	var s S
	_ = s.m(1)
	_ = s.m(int64(3))
}
`
	if err := typecheckOverload(t, src, true); err != nil {
		t.Fatalf("overload typecheck failed: %v", err)
	}
}

func TestOverloadErrors(t *testing.T) {
	t.Run("no_match", func(t *testing.T) {
		const src = `package p
func add(int) int { return 0 }
func add(int64) int64 { return 0 }
func f() { add("x") }
`
		err := typecheckOverload(t, src, false)
		if err == nil || !strings.Contains(err.Error(), "no matching overload") {
			t.Fatalf("got %v, want no matching overload", err)
		}
	})

	t.Run("ambiguous", func(t *testing.T) {
		const src = `package p
type A int
type B int
func add(A) int { return 0 }
func add(B) int { return 0 }
func f() { add(1) }
`
		err := typecheckOverload(t, src, false)
		if err == nil || !strings.Contains(err.Error(), "ambiguous overloaded") {
			t.Fatalf("got %v, want ambiguous overloaded", err)
		}
	})

	t.Run("duplicate_signature", func(t *testing.T) {
		const src = `package p
func dup(int) {}
func dup(int) {}
`
		err := typecheckOverload(t, src, true)
		if err == nil || !strings.Contains(err.Error(), "redeclared function dup") {
			t.Fatalf("got %v, want redeclared function dup", err)
		}
	})
}
