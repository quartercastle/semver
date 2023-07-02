package main

import (
	"bytes"
	"flag"
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

var (
	filter  string
	grep    string
	explain bool
)

func init() {
	flag.BoolVar(&explain, "explain", false, "explain reason behind decision")
	flag.StringVar(&filter, "filter", "", "filter between changes: patch, minor, major")
	flag.StringVar(&grep, "grep", "", "grep output")
}

func main() {
	flag.Parse()
	args := flag.Args()

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

	if !explain {
		fmt.Println(diff.Type(), time.Since(start))
		os.Exit(0)
	}

	for _, change := range diff {
		if filter != "" {
			if change.Type.String() != strings.ToUpper(filter) {
				continue
			}
		}

		if grep != "" {
			buffer := new(bytes.Buffer)
			printer.Fprint(buffer, token.NewFileSet(), change.Previous)
			printer.Fprint(buffer, token.NewFileSet(), change.Latest)
			if !strings.Contains(string(buffer.Bytes()), grep) {
				continue
			}
		}

		fmt.Printf("%s: %s\n", change.Type, change.Reason)

		if change.Previous != nil {
			fmt.Print("- ")
		}
		printer.Fprint(os.Stdout, token.NewFileSet(), change.Previous)
		if change.Previous != nil {
			fmt.Println()
		}
		if change.Latest != nil {
			fmt.Print("+ ")
		}
		printer.Fprint(os.Stdout, token.NewFileSet(), change.Latest)
		fmt.Println()
		if change.Latest != nil {
			fmt.Println()
		}
	}

	fmt.Println(diff.Type(), time.Since(start))
}
