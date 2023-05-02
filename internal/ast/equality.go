package ast

import (
	"fmt"
	"go/ast"
)

func equalIdent(a, b *ast.Ident) bool {
	if a == nil && b == nil {
		return true
	}

	if (a == nil && b != nil) || (a != nil && b == nil) {
		return false
	}

	return a.Name == b.Name && equalObject(a.Obj, b.Obj)
}

func equalFuncType(a, b *ast.FuncType) bool {
	if a == nil && b == nil {
		return true
	}

	if (a == nil && b != nil) || (a != nil && b == nil) {
		return false
	}

	return equalFieldList(a.TypeParams, b.TypeParams) &&
		equalFieldList(a.Params, b.Params) &&
		equalFieldList(a.Results, b.Results)
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

	if a.Kind == ast.Con && b.Kind == ast.Con {
		if a.Data.(int) != b.Data.(int) {
			return false
		}
	}

	return a.Kind == b.Kind && a.Name == b.Name
}

func equalChanType(a, b *ast.ChanType) bool {
	if a == nil && b == nil {
		return true
	}

	if (a == nil && b != nil) || (a != nil && b == nil) {
		return false
	}

	return a.Arrow == b.Arrow &&
		a.Dir == b.Dir &&
		equalExpr(a.Value, b.Value)
}

func equalMapType(a, b *ast.MapType) bool {
	return equalExpr(a.Key, b.Key) && equalExpr(a.Value, b.Value)
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
			return equalMapType(t, v)
		}
	case *ast.SelectorExpr:
		if v, ok := b.(*ast.SelectorExpr); ok {
			return equalIdent(t.Sel, v.Sel)
		}
	case *ast.StructType:
		if v, ok := b.(*ast.StructType); ok {
			return equalFieldList(t.Fields, v.Fields)
		}
	case *ast.ChanType:
		if v, ok := b.(*ast.ChanType); ok {
			return equalChanType(t, v)
		}
	case *ast.BasicLit:
		if v, ok := b.(*ast.BasicLit); ok {
			return equalBasicLit(t, v)
		}
	}

	fmt.Printf("DEBUG: %#v -> %#v\n", a, b)
	return a == b
}

func equalValueSpec(a, b *ast.ValueSpec) bool {
	if a == nil && b == nil {
		return true
	}

	if (a == nil && b != nil) || (a != nil && b == nil) {
		return false
	}

	return equalIdents(a.Names, b.Names) &&
		equalExprs(a.Values, b.Values) &&
		equalExpr(a.Type, b.Type)
}

func equalExprs(a, b []ast.Expr) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if !equalExpr(a[i], b[i]) {
			return false
		}
	}
	return true
}

func equalIdents(a, b []*ast.Ident) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if !equalIdent(a[i], b[i]) {
			return false
		}
	}

	return true
}

func equalNames(a, b []*ast.Ident) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i].Name != b[i].Name {
			return false
		}
	}

	return true
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
