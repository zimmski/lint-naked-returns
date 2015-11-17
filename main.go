package main

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/token"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/kisielk/gotool"
	"golang.org/x/tools/go/loader"
)

func foo() (hey string) {
	return
}

func main() {
	var opts struct {
		Tags []string `long:"tag" description:"Build tags"`
	}

	pkgs, err := flags.Parse(&opts)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)

		os.Exit(1)
	}

	if len(pkgs) == 0 {
		pkgs = append(pkgs, ".")
	}

	if err := lint(pkgs, opts.Tags); err != nil {
		fmt.Printf("ERROR: %s\n", err)

		os.Exit(1)
	}
}

func lint(pkgs []string, tags []string) error {
	ctx := build.Default
	for _, tag := range tags {
		ctx.BuildTags = append(ctx.BuildTags, tag)
	}

	conf := loader.Config{
		Build: &ctx,
	}

	gotoolCtx := gotool.Context{
		BuildContext: ctx,
	}

	for _, name := range gotoolCtx.ImportPaths(pkgs) {
		conf.ImportWithTests(name)
	}

	prog, err := conf.Load()
	if err != nil {
		return err
	}

	foundError := false
	for _, pkg := range prog.InitialPackages() {
		if pkg.Errors != nil {
			fmt.Printf("%s: %#v", pkg.String(), pkg.Errors)

			foundError = true
		}
	}
	if foundError {
		return fmt.Errorf("Found some initializing errors")
	}

	for _, pkg := range prog.InitialPackages() {
		for _, file := range pkg.Files {
			ast.Walk(walkForFunctions{
				FileSet: prog.Fset,
			}, file)
		}
	}

	return nil
}

type walkForFunctions struct {
	FileSet *token.FileSet
}

func (w walkForFunctions) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}

	if f, ok := n.(*ast.FuncDecl); ok && hasNamedReturnArgument(f) {
		ast.Walk(walkForNakedReturn{
			FileSet: w.FileSet,
		}, f.Body)
	}

	return w
}

func hasNamedReturnArgument(f *ast.FuncDecl) bool {
	if f.Type.Results == nil {
		return false
	}

	for _, r := range f.Type.Results.List {
		if len(r.Names) > 0 && r.Names[0].String() != "" {
			return true
		}
	}

	return false
}

type walkForNakedReturn struct {
	FileSet *token.FileSet
}

func (w walkForNakedReturn) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}

	if r, ok := n.(*ast.ReturnStmt); ok && len(r.Results) == 0 {
		fmt.Printf("%s: Naked return in function with named return arguments\n", w.FileSet.Position(r.Pos()))
	}

	return w
}
