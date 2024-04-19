package main

import (
	"gostable/common"

	"golang.org/x/tools/go/analysis"
)

type analyzerPlugin struct{}

func (*analyzerPlugin) GetAnalyzers() []*analysis.Analyzer {
	return []*analysis.Analyzer{
		common.Analyzer,
	}
}

var AnalyzerPlugin analyzerPlugin
