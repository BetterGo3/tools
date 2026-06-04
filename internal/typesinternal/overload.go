package typesinternal

import (
	"go/ast"
	"go/types"
)

// OverloadsForExpr returns all overload candidates for a function call
// expression's Fun, or nil if the callee is not overloaded.
func OverloadsForExpr(info *types.Info, fun ast.Expr) []*types.Func {
	if info == nil {
		return nil
	}
	if info.CallOverloads != nil {
		if fns, ok := info.CallOverloads[fun]; ok && len(fns) > 1 {
			return fns
		}
	}
	switch fun := fun.(type) {
	case *ast.Ident:
		if info.FuncOverloads != nil {
			if fns := info.FuncOverloads[fun.Name]; len(fns) > 1 {
				return fns
			}
		}
	case *ast.SelectorExpr:
		if info.MethodOverloads == nil {
			return nil
		}
		recvName := recvBaseName(info, fun.X)
		if recvName == "" {
			return nil
		}
		if fns := info.MethodOverloads[types.MethodOverloadKey{RecvName: recvName, Name: fun.Sel.Name}]; len(fns) > 1 {
			return fns
		}
	}
	return nil
}

func recvBaseName(info *types.Info, x ast.Expr) string {
	t := info.TypeOf(x)
	if t == nil {
		return ""
	}
	if p, ok := t.(*types.Pointer); ok {
		t = p.Elem()
	}
	if n, ok := t.(*types.Named); ok && n.Obj() != nil {
		return n.Obj().Name()
	}
	return ""
}
