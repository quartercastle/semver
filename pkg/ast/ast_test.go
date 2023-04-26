package ast

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"
)

func parse(a, b []string) (file1, file2 *ast.File, err error) {
	a = append([]string{"package foo"}, a...)
	b = append([]string{"package foo"}, b...)

	af, err := parser.ParseFile(token.NewFileSet(), "a.go", []byte(strings.Join(a, "\n")), 0)

	if err != nil {
		return nil, nil, fmt.Errorf("a: %w", err)
	}

	bf, err := parser.ParseFile(token.NewFileSet(), "b.go", []byte(strings.Join(b, "\n")), 0)

	if err != nil {
		return nil, nil, fmt.Errorf("b: %w", err)
	}

	return af, bf, nil
}

func TestCompare(t *testing.T) {
	tc := []struct {
		title            string
		previous, latest []string
		expected         Difference
	}{
		{
			"No difference",
			[]string{"func Foo()"},
			[]string{"func Foo()"},
			Patch,
		},
		{
			"additon of exported function",
			[]string{"func Foo()"},
			[]string{"func Foo()", "func Bar()"},
			Minor,
		},
		{
			"removal of exported function",
			[]string{"func Foo()", "func Bar()"},
			[]string{"func Foo()"},
			Major,
		},
		{
			"additon of internal function",
			[]string{"func Foo()"},
			[]string{"func Foo()", "func bar()"},
			Patch,
		},
		{
			"removal of internal function",
			[]string{"func Foo()", "func bar()"},
			[]string{"func Foo()"},
			Patch,
		},
		{
			"addition of argument in exported function",
			[]string{"func Foo()"},
			[]string{"func Foo(string)"},
			Major,
		},
		{
			"removal of argument in exported function",
			[]string{"func Foo(string)"},
			[]string{"func Foo()"},
			Major,
		},
		{
			"addition of return value in exported function",
			[]string{"func Foo()"},
			[]string{"func Foo() string"},
			Major,
		},
		{
			"removal of return value in exported function",
			[]string{"func Foo() string"},
			[]string{"func Foo()"},
			Major,
		},
		{
			"argument of different types in exported function",
			[]string{"func Foo(string)"},
			[]string{"func Foo(int)"},
			Major,
		},
		{
			"return value of different types in exported function",
			[]string{"func Foo() int"},
			[]string{"func Foo() string"},
			Major,
		},
		{
			"addition of argument in internal function",
			[]string{"func foo()"},
			[]string{"func foo(string)"},
			Patch,
		},
		{
			"removal of argument in internal function",
			[]string{"func foo(string)"},
			[]string{"func foo()"},
			Patch,
		},
		{
			"addition of return value in internal function",
			[]string{"func foo()"},
			[]string{"func foo() string"},
			Patch,
		},
		{
			"removal of return value in internal function",
			[]string{"func foo() string"},
			[]string{"func foo()"},
			Patch,
		},
		{
			"argument of different types in internal function",
			[]string{"func foo(string)"},
			[]string{"func foo(int)"},
			Patch,
		},
		{
			"return value of different types in internal function",
			[]string{"func foo() int"},
			[]string{"func foo() string"},
			Patch,
		},
	}

	for _, c := range tc {
		t.Run(c.title, func(t *testing.T) {
			previous, latest, err := parse(c.previous, c.latest)

			if err != nil {
				t.Error(err)
			}

			if actual := Compare(previous, latest); actual != c.expected {
				t.Errorf(
					"expected difference of %s; got %s",
					c.expected, actual,
				)
			}
		})
	}
}

func BenchmarkCompare(b *testing.B) {
	previous, latest, _ := parse(
		[]string{"func Foo()"},
		[]string{"func Foo(string)"},
	)
	for i := 0; i < b.N; i++ {
		Compare(previous, latest)
	}
}
