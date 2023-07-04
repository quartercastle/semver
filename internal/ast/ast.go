package ast

import (
	"go/ast"
	"reflect"
)

type Node = ast.Node
type Package = ast.Package

type Type int

func (t Type) String() string {
	return types[t]
}

var (
	types = map[Type]string{
		Patch: "PATCH",
		Minor: "MINOR",
		Major: "MAJOR",
	}
)

const (
	Patch Type = iota
	Minor
	Major
)

type comparator func(previous, latest Node) Diff

func compose(previous, latest Node) func(comparators ...comparator) Diff {
	return func(comparators ...comparator) Diff {
		diff := Diff{}
		for _, comparator := range comparators {
			diff = diff.Merge(comparator(previous, latest))
		}
		return diff
	}
}

func Compare(previous, latest ast.Node) Diff {
	diff := Diff{}
	if (previous == nil || reflect.ValueOf(previous).IsNil()) && (latest != nil || !reflect.ValueOf(latest).IsNil()) {
		return diff.Add(Change{
			Type:   Minor,
			Reason: "package has been added",
			Latest: latest,
		})
	}

	if (previous != nil || !reflect.ValueOf(previous).IsNil()) && (latest == nil || reflect.ValueOf(latest).IsNil()) {
		return diff.Add(Change{
			Type:     Major,
			Reason:   "package has been removed",
			Previous: previous,
		})
	}

	if (previous == nil || reflect.ValueOf(previous).IsNil()) && (latest == nil || reflect.ValueOf(latest).IsNil()) {
		return diff
	}

	return diff.Merge(compose(previous, latest)(
		comparePackage,
		compareValueSpec,
		compareFuncDecl,
		compareTypeSpec,
	))
}
