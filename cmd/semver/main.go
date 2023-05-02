package main

import (
	"fmt"
	"go/parser"
	"go/printer"
	"go/token"
	"io/fs"
	"os"
	"strings"
	"time"

	"github.com/quartercastle/semver/internal/ast"
)

func merge(a, b map[string]*ast.Package) map[string]struct{} {
	result := map[string]struct{}{}
	for k := range a {
		result[k] = struct{}{}
	}
	for k := range b {
		result[k] = struct{}{}
	}
	return result
}

func main() {
	args := os.Args[1:]

	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, "invalid arguments")
		os.Exit(1)
	}

	start := time.Now()
	a := token.NewFileSet()
	previous, err := parser.ParseDir(a, args[0], func(f fs.FileInfo) bool {
		return !strings.Contains(f.Name(), "_test.go")
	}, 0)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	b := token.NewFileSet()
	latest, err := parser.ParseDir(b, args[1], func(f fs.FileInfo) bool {
		return !strings.Contains(f.Name(), "_test.go")
	}, 0)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	all := merge(latest, previous)

	var diff ast.Diff
	for pkg := range all {
		diff = diff.Merge(ast.Compare(previous[pkg], latest[pkg]))
	}

	for _, change := range diff {
		fmt.Printf("%s: %s\n", change.Type, change.Reason)
		if change.Previous != nil {
			fmt.Println(a.Position(change.Previous.Pos()))
			fmt.Print("- ")
		}
		printer.Fprint(os.Stdout, token.NewFileSet(), change.Previous)
		if change.Previous != nil {
			fmt.Println()
		}
		if change.Latest != nil {
			fmt.Println(b.Position(change.Latest.Pos()))
			fmt.Print("+ ")
		}
		printer.Fprint(os.Stdout, token.NewFileSet(), change.Latest)
		fmt.Println()
		if change.Latest != nil {
			fmt.Println()
		}
	}

	fmt.Println(diff.Type(), "in:", time.Since(start))
}
