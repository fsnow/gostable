package main

import (
	"gostable/common"

	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(common.StableAnalyzer)
}
