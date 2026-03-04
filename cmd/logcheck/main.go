package main

import (
	"logcheck-selectel/pkg/analyzer"

	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(analyzer.Analyzer)
}
