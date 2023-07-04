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
	"path/filepath"
	"strings"
	"time"

	"github.com/quartercastle/semver/internal/ast"
)

func merge[T any](a, b map[string]T) map[string]struct{} {
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

func walk(origin, target string) (ast.Diff, error) {
	ignore := map[string]struct{}{
		".git":    {},
		".github": {},
	}

	previous := map[string]struct{}{}
	a, err := os.ReadDir(origin)

	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	for _, entry := range a {
		if !entry.IsDir() {
			continue
		}

		previous[entry.Name()] = struct{}{}
	}

	latest := map[string]struct{}{}
	b, err := os.ReadDir(target)

	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	for _, entry := range b {
		if !entry.IsDir() {
			continue
		}

		latest[entry.Name()] = struct{}{}
	}

	packages := merge(previous, latest)

	var diff ast.Diff
	for pkg := range packages {
		if _, ok := ignore[pkg]; ok {
			continue
		}

		d, err := walk(
			filepath.Join(origin, pkg),
			filepath.Join(target, pkg),
		)

		diff = diff.Merge(d)

		if err != nil {
			return diff, err
		}
	}

	d, err := compare(origin, target)
	return diff.Merge(d), err
}

func compare(origin, target string) (ast.Diff, error) {
	a := token.NewFileSet()
	previous, err := parser.ParseDir(a, origin, func(f fs.FileInfo) bool {
		return !strings.Contains(f.Name(), "_test.go")
	}, 0)

	if err != nil && !os.IsNotExist(err) {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	b := token.NewFileSet()
	latest, err := parser.ParseDir(b, target, func(f fs.FileInfo) bool {
		return !strings.Contains(f.Name(), "_test.go")
	}, 0)

	if err != nil && !os.IsNotExist(err) {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	all := merge(latest, previous)

	var diff ast.Diff
	for pkg := range all {
		diff = diff.Merge(ast.Compare(previous[pkg], latest[pkg]))
	}

	if !explain {
		return diff, nil
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
			if !strings.Contains(buffer.String(), grep) {
				continue
			}
		}

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

	return diff, nil
}

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, "invalid arguments")
		os.Exit(1)
	}

	start := time.Now()
	diff, err := walk(args[0], args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(diff.Type(), time.Since(start))
}
