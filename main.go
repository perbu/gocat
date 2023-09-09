package main

import (
	"os"
)

func main() {
	path := "."
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	fset, pkgs := parseDir(path)
	imports, decls := extractImportsAndDecls(pkgs)
	combinedFile := constructCombinedFile(imports, decls)
	printCombinedFile(fset, combinedFile)
}
