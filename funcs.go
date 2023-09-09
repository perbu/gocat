package main

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
)

func parseDir(path string) (*token.FileSet, map[string]*ast.Package) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, path, nil, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse directory: %v\n", err)
		os.Exit(1)
	}
	return fset, pkgs
}

func extractImportsAndDecls(pkgs map[string]*ast.Package) (map[string]struct{}, []ast.Decl) {
	imports := make(map[string]struct{})
	var decls []ast.Decl

	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				switch d := decl.(type) {
				case *ast.GenDecl:
					if d.Tok == token.IMPORT {
						for _, spec := range d.Specs {
							importSpec, ok := spec.(*ast.ImportSpec)
							if !ok {
								continue
							}
							imports[importSpec.Path.Value] = struct{}{}
						}
					} else {
						decls = append(decls, decl)
					}
				default:
					decls = append(decls, decl)
				}
			}
		}
	}

	return imports, decls
}

func constructCombinedFile(imports map[string]struct{}, decls []ast.Decl) *ast.File {
	// Prepare unique imports for the output
	importDecl := &ast.GenDecl{
		Tok: token.IMPORT,
	}
	for imp := range imports {
		importDecl.Specs = append(importDecl.Specs, &ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: imp,
			},
		})
	}
	decls = append([]ast.Decl{importDecl}, decls...)

	// Construct the final file AST
	return &ast.File{
		Name:  ast.NewIdent("main"),
		Decls: decls,
	}
}

func printCombinedFile(fset *token.FileSet, file *ast.File) {
	if err := format.Node(os.Stdout, fset, file); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to format code: %v\n", err)
		os.Exit(1)
	}
}
