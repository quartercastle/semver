package ast

import (
	"fmt"
	"go/ast"
)

func equalIdent(a, b *ast.Ident) bool {
	if !equalObject(a.Obj, b.Obj) {
		return false
	}

	return a.Name == b.Name
}

func equalFuncType(a, b *ast.FuncType) bool {
	if a == nil && b == nil {
		return true
	}

	if (a == nil && b != nil) || (a != nil && b == nil) {
		return false
	}

	if !equalFieldList(a.TypeParams, b.TypeParams) {
		return false
	}

	if !equalFieldList(a.Params, b.Params) {
		return false
	}

	if !equalFieldList(a.Results, b.Results) {
		return false
	}

	return true
}

func equalFieldList(a, b *ast.FieldList) bool {
	if a == nil && b == nil {
		return true
	}

	if (a == nil && b != nil) || (a != nil && b == nil) {
		return false
	}

	if len(a.List) != len(b.List) {
		return false
	}

	for i := range a.List {
		if !equalField(a.List[i], b.List[i]) {
			return false
		}
	}

	return true
}

func equalBasicLit(a, b *ast.BasicLit) bool {
	if a == nil && b == nil {
		return true
	}

	if (a == nil && b != nil) || (a != nil && b == nil) {
		return false
	}

	return a.Kind == b.Kind && a.Value == b.Value
}

func equalObject(a, b *ast.Object) bool {
	if a == nil && b == nil {
		return true
	}

	if (a == nil && b != nil) || (a != nil && b == nil) {
		return false
	}

	return a.Kind == b.Kind && a.Name == b.Name
}

func equalExpr(a, b ast.Expr) bool {
	switch t := a.(type) {
	case *ast.Ident:
		if v, ok := b.(*ast.Ident); ok {
			return equalIdent(t, v)
		}
	case *ast.Ellipsis:
		if v, ok := b.(*ast.Ellipsis); ok {
			return equalExpr(t.Elt, v.Elt)
		}
	case *ast.StarExpr:
		if v, ok := b.(*ast.StarExpr); ok {
			return equalExpr(t.X, v.X)
		}
	case *ast.ArrayType:
		if v, ok := b.(*ast.ArrayType); ok {
			return equalExpr(t.Elt, v.Elt)
		}
	case *ast.FuncType:
		if v, ok := b.(*ast.FuncType); ok {
			return equalFuncType(t, v)
		}
	case *ast.MapType:
		if v, ok := b.(*ast.MapType); ok {
			return equalExpr(t.Key, v.Key) && equalExpr(t.Value, v.Value)
		}
	case *ast.SelectorExpr:
		if v, ok := b.(*ast.SelectorExpr); ok {
			return equalIdent(t.Sel, v.Sel)
		}
	case *ast.StructType:
		if v, ok := b.(*ast.StructType); ok {
			return equalFieldList(t.Fields, v.Fields)
		}
	}

	fmt.Printf("DEBUG: %#v -> %#v\n", a, b)
	return false
}

func equalField(a, b *ast.Field) bool {
	/*if len(a.Names) != len(b.Names) {
		return false
	}

	for i := range a.Names {
		if a.Names[i].Obj.Kind == ast.Var {
			continue
		}

		if !equalIdent(a.Names[i], b.Names[i]) {
			return false
		}
	}*/

	if !equalExpr(a.Type, b.Type) {
		return false
	}

	if !equalBasicLit(a.Tag, b.Tag) {
		return false
	}

	return true
}
