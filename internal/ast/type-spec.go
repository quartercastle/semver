package ast

import (
	"go/ast"
)

func extractTypeSpec(node ast.Node) []*ast.TypeSpec {
	result := []*ast.TypeSpec{}

	if node == nil {
		return nil
	}

	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.TypeSpec:
			if !ast.IsExported(x.Name.String()) {
				return false
			}

			result = append(result, x)
		}
		return true
	})
	return result
}

func diffTypeSpec(a, b *ast.TypeSpec) Diff {
	var diff Diff
	if a == nil && b != nil {
		return diff.Add(Change{
			Type:   Minor,
			Reason: "struct has been added",
			Latest: b,
		})
	}

	if a != nil && b == nil {
		return diff.Add(Change{
			Type:     Major,
			Reason:   "struct has been removed",
			Previous: a,
		})
	}

	if !equalTypeSpec(a, b) {
		return diff.Add(Change{
			Type:     Major,
			Reason:   "struct has has changed signature",
			Previous: a,
			Latest:   b,
		})
	}

	return diff
}

func compareTypeSpec(a, b ast.Node) Diff {
	previous, latest := extractTypeSpec(a), extractTypeSpec(b)
	var diff Diff

	match := [][2]*ast.TypeSpec{}

	for _, p := range previous {
		match = append(match, [2]*ast.TypeSpec{p})
	}

	for _, l := range latest {
		var found bool
		for j, m := range match {
			p := m[0]

			if p == nil {
				break
			}

			if equalIdent(p.Name, l.Name) {
				match[j][1] = l
				found = true
				break
			}
		}

		if !found {
			match = append(match, [2]*ast.TypeSpec{nil, l})
		}
	}

	for _, m := range match {
		p, l := m[0], m[1]
		diff = diff.Merge(diffTypeSpec(p, l))
	}

	return diff
}