package cache

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"testing"
)

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
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "p.go", src, 0)
	if err != nil {
		t.Fatal(err)
	}
	cfg := &types.Config{IgnoreFuncBodies: true}
	_, err = cfg.Check("p", fset, []*ast.File{f}, nil)
	if err != nil {
		t.Fatalf("overload typecheck failed: %v", err)
	}
}
