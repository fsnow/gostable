package common

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const mongoPkgName = "go.mongodb.org/mongo-driver/mongo"

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

	unstableFunctions := map[string][]string{
		"Client":     {"Watch"},
		"Collection": {"Distinct", "SearchIndexes", "Watch"},
		"Database":   {"Watch"},
	}

	restrictedStages := []string{"$currentOp", "$indexStats", "$listLocalSessions", "$listSessions", "$planCacheStats", "$search"}
	restrictedOperators := map[string][]string{
		"$group":   {"$sum", "$avg"},
		"$project": {"$add", "$multiply"},
	}

	mongoCollection := mongoPkgName + ".Collection"
	inspect.Preorder(nodeFilter, func(node ast.Node) {
		call := node.(*ast.CallExpr)
		if isPkgDotFunction(pass, call, mongoCollection, "Aggregate") {
			fmt.Println("\nPreorder, Aggregate")
			if hasRestrictedStages(pass, call, restrictedStages, restrictedOperators) {
				pass.Reportf(call.Pos(), "usage of restricted pipeline stages detected")
			}
		}

		for driverType, functions := range unstableFunctions {
			for _, fnName := range functions {
				pkgName := mongoPkgName + "." + driverType
				if isPkgDotFunction(pass, call, pkgName, fnName) {
					pass.Reportf(call.Pos(), "use of %v.%v is not supported by the MongoDB Stable API", driverType, fnName)
				}
			}
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

	return false
}

func hasRestrictedStages(pass *analysis.Pass, call *ast.CallExpr, restrictedStages []string, restrictedOperators map[string][]string) bool {
	for _, arg := range call.Args {
		var bsonDocs []*ast.CompositeLit

		switch arg := arg.(type) {
		case *ast.Ident:
			bsonDocs = findPipelineVariable(pass, arg.Name, call.Pos())
			fmt.Println("ast.Ident case")
			fmt.Printf("len: %v\n", len(bsonDocs))
		case *ast.CompositeLit:
			if arg.Type == nil {
				bsonDocs = extractBSONDocs(arg.Elts)
				fmt.Println("ast.CompositeLit case, if")
			} else if ident, ok := arg.Type.(*ast.ArrayType); ok && ident.Elt != nil {
				if _, ok := ident.Elt.(*ast.StructType); ok {
					bsonDocs = extractBSONDocs(arg.Elts)
					fmt.Println("ast.CompositeLit case, else if if")
				}
			}
		case *ast.CallExpr:
			if sel, ok := arg.Fun.(*ast.SelectorExpr); ok && sel.Sel.Name == "Pipeline" {
				bsonDocs = extractBSONDocsFromPipeline(arg.Args)
				fmt.Println("ast.CallExpr case")
			}
		}

		for _, bsonDoc := range bsonDocs {
			if hasRestrictedStagesInDoc(bsonDoc, restrictedStages, restrictedOperators) {
				return true
			}
		}
	}
	return false
}

func findPipelineVariable(pass *analysis.Pass, varName string, pos token.Pos) []*ast.CompositeLit {
	var bsonDocs []*ast.CompositeLit

	fmt.Printf("findPipelineVariable, varName=%v\n", varName)
	fmt.Printf("pass.Files len: %v\n", len(pass.Files))

	var spec *ast.ValueSpec
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			if n == nil || n.Pos() > pos {
				return false
			}
			switch n := n.(type) {
			case *ast.ValueSpec:
				for _, name := range n.Names {
					if name.Name == varName {
						spec = n
						return false
					}
				}
			case *ast.AssignStmt:
				for _, lhs := range n.Lhs {
					if ident, ok := lhs.(*ast.Ident); ok && ident.Name == varName {
						if len(n.Rhs) == 1 {
							if compositeLit, ok := n.Rhs[0].(*ast.CompositeLit); ok {
								fmt.Println("before append 1")
								bsonDocs = append(bsonDocs, extractBSONDocs(compositeLit.Elts)...)
							}
						}
						return false
					}
				}
			}
			return true
		})
		if spec != nil {
			break
		}
	}

	if spec != nil {
		for _, value := range spec.Values {
			if compositeLit, ok := value.(*ast.CompositeLit); ok {
				fmt.Println("before append 2")
				bsonDocs = append(bsonDocs, extractBSONDocs(compositeLit.Elts)...)
			}
		}
	}

	return bsonDocs
}

func hasRestrictedStagesInDoc(bsonDoc *ast.CompositeLit, restrictedStages []string, restrictedOperators map[string][]string) bool {
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
	return false
}

func extractBSONDocs(elts []ast.Expr) []*ast.CompositeLit {
	var bsonDocs []*ast.CompositeLit
	for _, elt := range elts {
		if bsonDoc, ok := elt.(*ast.CompositeLit); ok {
			bsonDocs = append(bsonDocs, bsonDoc)
		}
	}
	return bsonDocs
}

func extractBSONDocsFromPipeline(args []ast.Expr) []*ast.CompositeLit {
	var bsonDocs []*ast.CompositeLit
	for _, arg := range args {
		if unaryExpr, ok := arg.(*ast.UnaryExpr); ok {
			if callExpr, ok := unaryExpr.X.(*ast.CallExpr); ok {
				if len(callExpr.Args) == 1 {
					if bsonDoc, ok := callExpr.Args[0].(*ast.CompositeLit); ok {
						bsonDocs = append(bsonDocs, bsonDoc)
					}
				}
			} else if bsonDoc, ok := unaryExpr.X.(*ast.CompositeLit); ok {
				bsonDocs = append(bsonDocs, bsonDoc)
			}
		}
	}
	return bsonDocs
}
