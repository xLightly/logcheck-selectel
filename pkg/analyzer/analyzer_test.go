package analyzer_test

import (
	"logcheck-selectel/pkg/analyzer"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, analyzer.Analyzer, "a", "b", "c", "d")
}
