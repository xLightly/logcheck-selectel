package main

import (
	"logcheck-selectel/pkg/analyzer"

	"golang.org/x/tools/go/analysis"
)

type AnalyzerPlugin struct{}

func (a *AnalyzerPlugin) GetAnalyzers() []*analysis.Analyzer {
	return []*analysis.Analyzer{analyzer.Analyzer}
}

func New(_ interface{}) ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{analyzer.Analyzer}, nil
}

var AnalyzerPluginVar = AnalyzerPlugin{}
