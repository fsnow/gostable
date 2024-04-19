package common

import (
	"fmt"
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const mongoCollection = "go.mongodb.org/mongo-driver/mongo.Collection"

var Analyzer = &analysis.Analyzer{
	Name: "gostable",
	Doc:  "ensures that all MongoDB go driver code adheres to the Stable API",
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	restrictedStages := []string{"$currentOp", "$indexStats", "$listLocalSessions", "$listSessions", "$planCacheStats", "$search"}
	restrictedOperators := map[string][]string{
		"$group":   {"$sum", "$avg"},
		"$project": {"$add", "$multiply"},
	}

	inspect.Preorder(nodeFilter, func(node ast.Node) {
		call := node.(*ast.CallExpr)
		if isPkgDotFunction(pass, call, mongoCollection, "Aggregate") {
			if hasRestrictedStages(call, restrictedStages, restrictedOperators) {
				pass.Reportf(call.Pos(), "usage of restricted pipeline stages detected")
			}
		}

		if isPkgDotFunction(pass, call, mongoCollection, "SearchIndexes") {
			pass.Reportf(call.Pos(), "usage of collection.SearchIndexes is not supported by the Stable API")
		}

		if isPkgDotFunction(pass, call, mongoCollection, "Watch") {
			pass.Reportf(call.Pos(), "usage of collection.Watch is not supported by the Stable API")
		}
	})

	return nil, nil
}

func isPkgDotFunction(pass *analysis.Pass, call *ast.CallExpr, packagePath, functionName string) bool {
	// Check if the call expression is a selector expression
	selExpr, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	// Check if the selector matches the function name
	if selExpr.Sel.Name != functionName {
		return false
	}

	// Check if the selector's X expression is a package identifier
	_, ok = selExpr.X.(*ast.Ident)
	if !ok {
		fmt.Println("1")
		return false
	}

	typ := pass.TypesInfo.TypeOf(selExpr.X)
	if typ == nil {
		fmt.Println("2")
		return false
	}

	typStr := typ.String()
	if len(typStr) > 0 && typStr[0] == '*' {
		typStr = typStr[1:]
	}

	if typStr == packagePath {
		return true
	}

	fmt.Println("false")
	return false
}

func hasRestrictedStages(call *ast.CallExpr, restrictedStages []string, restrictedOperators map[string][]string) bool {
	for _, arg := range call.Args {
		if bsonDoc, ok := arg.(*ast.CompositeLit); ok {
			for _, elt := range bsonDoc.Elts {
				if kv, ok := elt.(*ast.KeyValueExpr); ok {
					if key, ok := kv.Key.(*ast.BasicLit); ok {
						stageName := strings.Trim(key.Value, "\"")

						// Check if the stage is prohibited
						for _, stage := range restrictedStages {
							if stageName == stage {
								return true
							}
						}

						// Check if the stage has restricted operators
						if restrictedOperators, ok := restrictedOperators[stageName]; ok {
							if bsonDoc, ok := kv.Value.(*ast.CompositeLit); ok {
								for _, elt := range bsonDoc.Elts {
									if kv, ok := elt.(*ast.KeyValueExpr); ok {
										if key, ok := kv.Key.(*ast.BasicLit); ok {
											operatorName := strings.Trim(key.Value, "\"")

											for _, operator := range restrictedOperators {
												if operatorName == operator {
													return true
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return false
}
