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
		expected         Type
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
		{
			"variable argument names",
			[]string{"func Foo(foo int) string"},
			[]string{"func Foo(bar int) string"},
			Patch,
		},
		{
			"equal spread operator argument",
			[]string{"func Foo(...int)"},
			[]string{"func Foo(...int)"},
			Patch,
		},
		{
			"spread operator argument of different type",
			[]string{"func Foo(...int)"},
			[]string{"func Foo(...string)"},
			Major,
		},
		{
			"equal method",
			[]string{"func (Foo) Bar()"},
			[]string{"func (Foo) Bar()"},
			Patch,
		},
		{
			"addition of argument in exported method",
			[]string{"func (Foo) Bar()"},
			[]string{"func (Foo) Bar(string)"},
			Major,
		},
		{
			"removal of argument in exported method",
			[]string{"func (Foo) Bar(string)"},
			[]string{"func (Foo) Bar()"},
			Major,
		},
		{
			"addition of return value in exported method",
			[]string{"func (Foo) Bar()"},
			[]string{"func (Foo) Bar() string"},
			Major,
		},
		{
			"removal of return value in exported method",
			[]string{"func (Foo) Bar() string"},
			[]string{"func (Foo) Bar()"},
			Major,
		},
		{
			"addition of argument in internal method",
			[]string{"func (Foo) bar()"},
			[]string{"func (Foo) bar(string)"},
			Patch,
		},
		{
			"removal of argument in internal method",
			[]string{"func (Foo) bar(string)"},
			[]string{"func (Foo) bar()"},
			Patch,
		},
		{
			"addition of return value in internal method",
			[]string{"func (Foo) bar()"},
			[]string{"func (Foo) bar() string"},
			Patch,
		},
		{
			"removal of return value in internal method",
			[]string{"func (Foo) bar() string"},
			[]string{"func (Foo) bar()"},
			Patch,
		},
		{
			"receiver concrete type",
			[]string{"func (Foo) Foo()"},
			[]string{"func (Foo) Foo()"},
			Patch,
		},
		{
			"receiver pointer type",
			[]string{"func (*Foo) Foo()"},
			[]string{"func (*Foo) Foo()"},
			Patch,
		},
		{
			"receiver change from a concrete to a pointer type",
			[]string{"func (*Foo) Foo()"},
			[]string{"func (Foo) Foo()"},
			Major,
		},
		{
			"receiver change from a pointer to a concrete type",
			[]string{"func (Foo) Foo()"},
			[]string{"func (*Foo) Foo()"},
			Major,
		},
		{
			"change of receiver type",
			[]string{"func (Bar) Foo()"},
			[]string{"func (Foo) Foo()"},
			Major,
		},
		{
			"change of receiver variable",
			[]string{"func (a Foo) Foo()"},
			[]string{"func (b Foo) Foo()"},
			Patch,
		},
		{
			"addition of receiver function",
			[]string{"func (Foo) Bar()"},
			[]string{"func (Foo) Bar()", "func (Foo) Foo()"},
			Minor,
		},
		{
			"removal of receiver function",
			[]string{"func (Foo) Bar()", "func (Foo) Foo()"},
			[]string{"func (Foo) Bar()"},
			Major,
		},
		{
			"addition of exported function",
			[]string{""},
			[]string{"func Foo()"},
			Minor,
		},
		{
			"removal of exported function",
			[]string{"func Foo()"},
			[]string{""},
			Major,
		},
		{
			"typed paramterer argument",
			[]string{"func Foo[A any](A)"},
			[]string{"func Foo[A any](A)"},
			Patch,
		},
		{
			"addition of typed paramterer argument",
			[]string{"func Foo[A any]() A"},
			[]string{"func Foo[A, B any]() (A, B)"},
			Major,
		},
		{
			"removal of typed paramterer argument",
			[]string{"func Foo[A, B any]() (A, B)"},
			[]string{"func Foo[A any]() A "},
			Major,
		},
		{
			"typed paramterer return value",
			[]string{"func Foo[A any]() A"},
			[]string{"func Foo[A any]() A"},
			Patch,
		},
		{
			"addition of typed paramterer return value",
			[]string{"func Foo[A any]() A"},
			[]string{"func Foo[A, B any]() (A, B)"},
			Major,
		},
		{
			"removal of typed paramterer return value",
			[]string{"func Foo[A, B any]() (A, B)"},
			[]string{"func Foo[A any]() A"},
			Major,
		},
		{
			"changing internal receiver type to exported",
			[]string{"func (*foo) Bar()"},
			[]string{"func (*Foo) Bar()"},
			Minor,
		},
		{
			"equal callback function argument",
			[]string{"func Foo(func (int))"},
			[]string{"func Foo(func (int))"},
			Patch,
		},
		{
			"changing callback function arguments",
			[]string{"func Foo(func(string))"},
			[]string{"func Foo(func(int))"},
			Major,
		},
		{
			"change of type in map argument",
			[]string{"func Foo(map[string]string)"},
			[]string{"func Foo(map[string]int)"},
			Major,
		},
		{
			"function with a chan argument",
			[]string{"func Foo(chan int)"},
			[]string{"func Foo(chan int)"},
			Patch,
		},
		{
			"function with chan arguments of different directions",
			[]string{"func Foo(chan<- int)"},
			[]string{"func Foo(<-chan int)"},
			Major,
		},
		{
			"function with chan arguments of same directions",
			[]string{"func Foo(<-chan int)"},
			[]string{"func Foo(<-chan int)"},
			Patch,
		},
		{
			"function with chan arguments of different types",
			[]string{"func Foo(chan int)"},
			[]string{"func Foo(chan string)"},
			Major,
		},

		// consts
		{
			"exported const with same type",
			[]string{"const Foo int = iota"},
			[]string{"const Foo int = iota"},
			Patch,
		},
		{
			"exported const with inferred type",
			[]string{"const Foo = \"bar\""},
			[]string{"const Foo = \"bar\""},
			Patch,
		},
		{
			"addition of exported const with iota",
			[]string{"const Foo = iota"},
			[]string{
				"const (",
				"	Foo = iota",
				"   Bar",
				")",
			},
			Minor,
		},
		{
			"removal of exported const with iota",
			[]string{
				"const (",
				"	Foo = iota",
				"   Bar",
				")",
			},
			[]string{"const Foo = iota"},
			Major,
		},
		{
			"changing order of iota const",
			[]string{"const Foo = iota"},
			[]string{
				"const (",
				"	Bar = iota",
				"   Foo",
				")",
			},
			Major,
		},
		{
			"change of type for exported const",
			[]string{"const Test int = 0"},
			[]string{"const Test string = \"\""},
			Major,
		},
		{
			"ignore internal const",
			[]string{"const test int = 0"},
			[]string{"const test string = \"\""},
			Patch,
		},

		// type spec
		{
			"addition of exported type",
			[]string{""},
			[]string{"type Foo struct {}"},
			Minor,
		},
		{
			"removal of exported type",
			[]string{"type Foo struct {}"},
			[]string{""},
			Major,
		},
		{
			"changing type of exported type",
			[]string{"type Foo struct {}"},
			[]string{"type Foo func()"},
			Major,
		},
		{
			"changing order of fields in struct",
			[]string{
				"type Foo struct {",
				"	Foo int",
				"	Bar int",
				"}",
			},
			[]string{
				"type Foo struct {",
				"	Foo int",
				"	Baz int",
				"	Bar int",
				"}",
			},
			Major,
		},
		{
			"changing order of fields in struct",
			[]string{
				"type Foo struct {",
				"	Foo int",
				"	Bar int",
				"}",
			},
			[]string{
				"type Foo struct {",
				"	Foo int",
				"	Baz int",
				"	Bar int",
				"}",
			},
			Major,
		},
	}

	for _, c := range tc {
		t.Run(c.title, func(t *testing.T) {
			previous, latest, err := parse(c.previous, c.latest)

			if err != nil {
				t.Error(err)
			}

			if actual := Compare(previous, latest); actual.Type() != c.expected {
				t.Errorf(
					"expected difference of %s; got %s",
					c.expected, actual.Type(),
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
