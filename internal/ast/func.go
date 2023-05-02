package ast

import "go/ast"

func extractFuncs(node ast.Node) []*ast.FuncDecl {
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

func compareFuncDecl(a, b *ast.FuncDecl) Diff {
	diff := Diff{}

	if a == nil {
		b.Body = nil
		return diff.Add(Change{
			Type:   Minor,
			Reason: "function has been added",
			Latest: b,
		})
	}

	if b == nil {
		a.Body = nil

		if a.Recv != nil {
			for _, field := range a.Recv.List {
				// internal receiver, not breaking
				var name string
				switch t := field.Type.(type) {
				case *ast.Ident:
					name = t.Name
				case *ast.StarExpr:
					if v, ok := t.X.(*ast.Ident); ok {
						name = v.Name
					}
				}

				if !ast.IsExported(name) {
					return diff
				}
			}
		}

		return diff.Add(Change{
			Type:     Major,
			Reason:   "function has been removed",
			Previous: a,
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

func compareFuncs(a, b Node) Diff {
	previous, latest := extractFuncs(a), extractFuncs(b)
	diff := Diff{}

	match := [][2]*ast.FuncDecl{}

	for _, p := range previous {
		match = append(match, [2]*ast.FuncDecl{p})
	}

	for _, l := range latest {
		var found bool
		for j, m := range match {
			p := m[0]

			if p == nil {
				break
			}

			if equalIdent(p.Name, l.Name) && equalFieldList(p.Recv, l.Recv) {
				match[j][1] = l
				found = true
			}
		}
		if !found {
			match = append(match, [2]*ast.FuncDecl{nil, l})
		}
	}

	for _, m := range match {
		p, l := m[0], m[1]
		diff = diff.Merge(compareFuncDecl(p, l))
	}

	return diff
}
