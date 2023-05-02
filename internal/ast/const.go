package ast

import (
	"go/ast"
)

func extractConsts(node ast.Node) []*ast.ValueSpec {
	result := []*ast.ValueSpec{}

	if node == nil {
		return nil
	}

	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.ValueSpec:
			for _, ident := range x.Names {
				if !ast.IsExported(ident.Name) {
					return false
				}
			}

			result = append(result, x)
		}
		return true
	})
	return result
}

func compareValueSpec(a, b *ast.ValueSpec) Diff {
	var diff Diff
	if a == nil && b != nil {
		return diff.Add(Change{
			Type:   Minor,
			Reason: "constant has been added",
			Latest: b,
		})
	}

	if a != nil && b == nil {
		return diff.Add(Change{
			Type:     Major,
			Reason:   "constant has been removed",
			Previous: a,
		})
	}

	if !equalValueSpec(a, b) {
		return diff.Add(Change{
			Type:     Major,
			Reason:   "constant has has changed signature",
			Previous: a,
			Latest:   b,
		})
	}

	return diff
}

func compareConsts(a, b ast.Node) Diff {
	previous, latest := extractConsts(a), extractConsts(b)
	var diff Diff

	match := [][2]*ast.ValueSpec{}

	for _, p := range previous {
		match = append(match, [2]*ast.ValueSpec{p})
	}

	for _, l := range latest {
		var found bool
		for j, m := range match {
			p := m[0]

			if p == nil {
				break
			}

			if equalNames(p.Names, l.Names) {
				match[j][1] = l
				found = true
				break
			}
		}

		if !found {
			match = append(match, [2]*ast.ValueSpec{nil, l})
		}
	}

	for _, m := range match {
		p, l := m[0], m[1]
		diff = diff.Merge(compareValueSpec(p, l))
	}

	return diff
}
