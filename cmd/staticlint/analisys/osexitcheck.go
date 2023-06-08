// Package analisys для статических анализаторов кода.
package analisys

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// OsExitCheckAnalyzer   проверяет
// вызовы os.Exit в пакете main, функции main.
var OsExitCheckAnalyzer = &analysis.Analyzer{
	Name: "OsExitCheckAnalyzer",
	Doc:  "check for os.exit in main package",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			switch x := node.(type) {
			case *ast.File:
				if x.Name.Name != "main" {
					return false
				}
			case *ast.SelectorExpr:
				if x.Sel.Name == "Exit" {
					pass.Reportf(x.Pos(), "expression has os.Exit call in main package")

				}
			}
			return true
		})
	}
	return nil, nil
}
