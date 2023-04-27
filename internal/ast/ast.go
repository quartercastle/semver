package ast

import (
	"fmt"
	"go/ast"
	"reflect"
)

type Node = ast.Node
type Package = ast.Package

type Change struct {
	Type             Type
	Reason           string
	Previous, Latest Node
}

type Diff []Change

func (d Diff) Add(c ...Change) Diff {
	if d == nil {
		d = Diff{}
	}
	return append(d, c...)
}

func (d Diff) Merge(q Diff) Diff {
	if d == nil {
		d = Diff{}
	}
	return append(d, q...)
}

func (d Diff) Type() Type {
	diff := Patch
	for _, change := range d {
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
	}

	return fmt.Sprint(a) == fmt.Sprint(b)
}

func equalField(a, b *ast.Field) bool {
	if len(a.Names) != len(b.Names) {
		return false
	}

	for i := range a.Names {
		if a.Names[i].Obj.Kind == ast.Var {
			continue
		}

		if !equalIdent(a.Names[i], b.Names[i]) {
			return false
		}
	}

	if !equalExpr(a.Type, b.Type) {
		return false
	}

	if !equalBasicLit(a.Tag, b.Tag) {
		return false
	}

	return true
}

func compareFuncDecl(a, b *ast.FuncDecl) Diff {
	diff := Diff{}

	if a == nil {
		b.Body = nil
		return diff.Add(Change{
			Type:     Minor,
			Reason:   "function has been added",
			Previous: b,
		})
	}

	if b == nil {
		a.Body = nil
		return diff.Add(Change{
			Type:   Major,
			Reason: "function has been removed",
			Latest: a,
		})
	}

	if equalIdent(a.Name, b.Name) && !equalFieldList(a.Recv, b.Recv) {
		a.Body, b.Body = nil, nil
		return diff.Add(Change{
			Type:     Major,
			Reason:   "receiver type has changed",
			Previous: a,
			Latest:   b,
		})
	}

	if !equalFuncType(a.Type, b.Type) {
		a.Body, b.Body = nil, nil
		return diff.Add(Change{
			Type:     Major,
			Reason:   "function signature has changed",
			Previous: a,
			Latest:   b,
		})
	}

	return diff
}

func compareFuncs(previous, latest Node) Diff {
	previousFuncs, latestFuncs := newFuncs(previous), newFuncs(latest)
	diff := Diff{}

	matchPrevious := map[*ast.FuncDecl]bool{}

	for _, p := range previousFuncs {
		matchPrevious[p] = false
	}

	for i := range latestFuncs {
		for j := range previousFuncs {
			p, l := previousFuncs[j], latestFuncs[i]

			if matchPrevious[p] {
				break
			}

			if i >= len(previousFuncs) {
				// an exported function has been added
				return diff.Merge(compareFuncDecl(nil, l))
			}

			if equalIdent(p.Name, l.Name) && equalFieldList(p.Recv, l.Recv) {
				matchPrevious[p] = true
			}

			diff = diff.Merge(compareFuncDecl(p, l))
		}
	}

	for _, p := range previousFuncs {
		if !matchPrevious[p] {
			// an exported function has been removed
			return diff.Merge(compareFuncDecl(p, nil))
		}
	}

	return diff
}

type Comparator func(previous, latest Node) Diff

func compose(previous, latest Node) func(comparators ...Comparator) Diff {
	return func(comparators ...Comparator) Diff {
		diff := Diff{}
		for _, comparator := range comparators {
			diff = diff.Merge(comparator(previous, latest))
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
		return diff.Add(Change{
			Type:   Major,
			Reason: "removal of package",
		})
	}

	if (previous != nil || !reflect.ValueOf(previous).IsNil()) && (latest == nil || reflect.ValueOf(latest).IsNil()) {
		return diff.Add(Change{
			Type:   Minor,
			Reason: "addition of package",
		})
	}

	return diff.Merge(compose(previous, latest)(
		compareFuncs,
	))
}
