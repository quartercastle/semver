package ast

import (
	"fmt"
	"go/ast"
	"reflect"
)

func comparePackage(previous, latest Node) Diff {
	var diff Diff

	if (previous != nil || !reflect.ValueOf(previous).IsNil()) && (latest == nil || reflect.ValueOf(latest).IsNil()) {
		if v, ok := previous.(*ast.Package); ok {
			return diff.Add(Change{
				Type:   Major,
				Reason: fmt.Sprintf("removal of package %s", v.Name),
			})
		}
	}

	if (previous == nil || reflect.ValueOf(previous).IsNil()) && (latest != nil || !reflect.ValueOf(latest).IsNil()) {
		if v, ok := latest.(*ast.Package); ok {
			return diff.Add(Change{
				Type:   Minor,
				Reason: fmt.Sprintf("addition of package %s", v.Name),
			})
		}
	}

	return diff
}
