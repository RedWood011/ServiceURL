package analisys

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestOsExitCheckAnalyzer_Analysistest(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), OsExitCheckAnalyzer, "./...")
}
