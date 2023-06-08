// Пакет для статистического анализа кода.
// Запустить static-lint <path_to_files>, где path_to_files путь до файлов проверки
// Все проверки SA static-check.io,
// Go-critic and nil err linters
// - OsExitAnalyzer  проверка os.Exit, что не вызывается в пакете main
package main

import (
	"log"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"
	"staticlint/analisys"
	"staticlint/config"
)

func main() {
	cfg := config.NewConfig()
	err := config.ReadConfig(cfg)
	if err != nil {
		log.Fatal(err)
	}
	// Изначальные проверки
	myChecks := []*analysis.Analyzer{
		analisys.OsExitCheckAnalyzer,
		shadow.Analyzer,
		printf.Analyzer,
		structtag.Analyzer,
	}

	checks := make(map[string]bool)
	for _, v := range cfg.StaticCheck {
		checks[v] = true
	}

	// Добавить необходимые проверки
	for _, v := range staticcheck.Analyzers {
		myChecks = append(myChecks, v.Analyzer)
	}

	multichecker.Main(myChecks...)
}
