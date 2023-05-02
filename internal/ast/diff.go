package ast

import "go/ast"

type Change struct {
	Type             Type
	Reason           string
	Previous, Latest ast.Node
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
