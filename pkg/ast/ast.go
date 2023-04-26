package ast

import (
	"fmt"
	"go/ast"
)

type Node = ast.Node

type Difference int

func (d Difference) String() string {
	return difference[d]
}

var (
	difference = map[Difference]string{
		None:  "NONE",
		Patch: "PATCH",
		Minor: "MINOR",
		Major: "MAJOR",
	}
)

const (
	None Difference = iota
	Patch
	Minor
	Major
)

func newFuncs(node ast.Node) []*ast.FuncDecl {
	result := []*ast.FuncDecl{}
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			if !ast.IsExported(x.Name.Name) {
				return true
			}
			result = append(result, x)
		}
		return true
	})
	return result
}

func equalIdent(a, b *ast.Ident) bool {
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

func equalExpr(a, b ast.Expr) bool {
	return fmt.Sprint(a) == fmt.Sprint(b)
}

func equalField(a, b *ast.Field) bool {
	/*if len(a.Names) != len(b.Names) {
		return false
	}

	for i := range a.Names {
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

func compareFuncDecl(a, b *ast.FuncDecl) Difference {
	diff := None

	if !equalIdent(a.Name, b.Name) {
		return diff
	}

	if !equalFuncType(a.Type, b.Type) {
		return set(diff, Major)
	}

	return diff
}

func compareFuncs(previous, latest Node) Difference {
	previousFuncs, latestFuncs := newFuncs(previous), newFuncs(latest)
	diff := Patch

	if len(latestFuncs) < len(previousFuncs) {
		// a exported function has been removed
		return Major
	}

	for i := range latestFuncs {
		for j := range previousFuncs {
			if i >= len(previousFuncs) {
				// an exported function has been added
				return set(diff, Minor)
			}

			diff = set(
				diff,
				compareFuncDecl(previousFuncs[j], latestFuncs[i]),
			)
		}
	}

	return diff
}

type Comparator func(previous, latest Node) Difference

func set(current, latest Difference) Difference {
	if current < latest {
		return latest
	}
	return current
}

func compose(previous, latest Node) func(comparators ...Comparator) Difference {
	return func(comparators ...Comparator) Difference {
		diff := None
		for _, comparator := range comparators {
			diff = set(diff, comparator(previous, latest))
		}
		return diff
	}
}

func Compare(previous, latest ast.Node) Difference {
	return compose(previous, latest)(
		compareFuncs,
	)
}
