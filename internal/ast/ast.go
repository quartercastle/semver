package ast

import (
	"fmt"
	"go/ast"
	"go/token"
	"reflect"
)

type Node = ast.Node
type Package = ast.Package

type Change struct {
	Type             Type
	Reason           string
	Previous, Latest Snippet
}

type Snippet struct{ Pos, End token.Pos }

type Diff map[*Change]struct{}

func (d Diff) Set(c Change) Diff {
	if d == nil {
		d = Diff{}
	}
	d[&c] = struct{}{}
	return d
}

func (d Diff) Merge(q Diff) Diff {
	if d == nil {
		d = Diff{}
	}
	for k, v := range q {
		d[k] = v
	}
	return d
}

func (d Diff) Type() Type {
	diff := Patch
	for change := range d {
		if diff < change.Type {
			diff = change.Type
		}
	}
	return diff
}

type Type int

func (t Type) String() string {
	return types[t]
}

var (
	types = map[Type]string{
		None:  "NONE",
		Patch: "PATCH",
		Minor: "MINOR",
		Major: "MAJOR",
	}
)

const (
	None Type = iota
	Patch
	Minor
	Major
)

func newFuncs(node ast.Node) []*ast.FuncDecl {
	result := []*ast.FuncDecl{}

	if node == nil {
		return nil
	}

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

func compareFuncDecl(diff Diff, a, b *ast.FuncDecl) Diff {
	if !equalIdent(a.Name, b.Name) {
		return diff
	}

	if !equalFuncType(a.Type, b.Type) {
		return diff.Set(Change{
			Type:     Major,
			Reason:   "function signature has changed",
			Previous: Snippet{a.Type.Pos(), a.Type.End()},
			Latest:   Snippet{b.Type.Pos(), b.Type.End()},
		})
	}

	return diff
}

func compareFuncs(diff Diff, previous, latest Node) Diff {
	previousFuncs, latestFuncs := newFuncs(previous), newFuncs(latest)

	if len(latestFuncs) < len(previousFuncs) {
		return diff.Set(Change{
			Type:   Major,
			Reason: "an exported function has been removed",
		})
	}

	for i := range latestFuncs {
		for j := range previousFuncs {
			if i >= len(previousFuncs) {
				// an exported function has been added
				return diff.Set(Change{
					Type:   Minor,
					Reason: "an exported function has been added",
				})
			}

			diff.Merge(compareFuncDecl(diff, previousFuncs[j], latestFuncs[i]))
		}
	}

	return diff
}

type Comparator func(diff Diff, previous, latest Node) Diff

func compose(diff Diff, previous, latest Node) func(comparators ...Comparator) Diff {
	return func(comparators ...Comparator) Diff {
		for _, comparator := range comparators {
			diff.Merge(comparator(diff, previous, latest))
		}
		return diff
	}
}

func Compare(previous, latest ast.Node) Diff {
	diff := Diff{}

	if (previous == nil || reflect.ValueOf(previous).IsNil()) && (latest == nil || reflect.ValueOf(latest).IsNil()) {
		return diff
	}

	if (previous == nil || reflect.ValueOf(previous).IsNil()) && (latest != nil || !reflect.ValueOf(latest).IsNil()) {
		return diff.Set(Change{
			Type:   Major,
			Reason: "removal of package",
		})
	}

	if (previous != nil || !reflect.ValueOf(previous).IsNil()) && (latest == nil || reflect.ValueOf(latest).IsNil()) {
		return diff.Set(Change{
			Type:   Minor,
			Reason: "addition of package",
		})
	}

	return compose(diff, previous, latest)(
		compareFuncs,
	)
}
