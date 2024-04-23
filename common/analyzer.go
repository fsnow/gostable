package common

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const mongoPkgName = "go.mongodb.org/mongo-driver/mongo"

const optsPkgName = "go.mongodb.org/mongo-driver/mongo/options"

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

	unstableFunctions := map[string]map[string][]string{
		mongoPkgName: {
			"Client":     {"Watch"},
			"Collection": {"Distinct", "SearchIndexes", "Watch"},
			"Database":   {"Watch"},
		},
		optsPkgName: {
			"CreateCollectionOptions": {"SetCapped", "SetDefaultIndexOptions", "SetMaxDocuments", "SetSizeInBytes", "SetStorageEngine"},
			"FindOneAndDeleteOptions": {"SetMax", "SetMaxAwaitTime", "SetMin", "SetNoCursorTimeout", "SetOplogReplay", "SetReturnKey",
				"SetShowRecordID"},
			"FindOneAndReplaceOptions": {"SetMax", "SetMaxAwaitTime", "SetMin", "SetNoCursorTimeout", "SetOplogReplay", "SetReturnKey",
				"SetShowRecordID"},
			"FindOneAndUpdateOptions": {"SetMax", "SetMaxAwaitTime", "SetMin", "SetNoCursorTimeout", "SetOplogReplay", "SetReturnKey",
				"SetShowRecordID"},
			"FindOneOptions": {"SetMax", "SetMaxAwaitTime", "SetMin", "SetNoCursorTimeout", "SetOplogReplay", "SetReturnKey",
				"SetShowRecordID"},
			"FindOptions": {"SetCursorType", "SetMax", "SetMaxAwaitTime", "SetMin", "SetNoCursorTimeout", "SetOplogReplay", "SetReturnKey",
				"SetShowRecordID"},
			"IndexOptions": {"SetBackground", "SetBucketSize", "SetSparse", "SetStorageEngine"},
		},
	}

	unstableOptionsStructs := map[string][]string{
		"CreateCollectionOptions": {"Capped", "DefaultIndexOptions", "MaxDocuments", "SizeInBytes", "StorageEngine"},
		"FindOneAndDeleteOptions": {"Max", "MaxAwaitTime", "Min", "NoCursorTimeout", "OplogReplay", "ReturnKey",
			"ShowRecordID"},
		"FindOneAndReplaceOptions": {"Max", "MaxAwaitTime", "Min", "NoCursorTimeout", "OplogReplay", "ReturnKey",
			"ShowRecordID"},
		"FindOneAndUpdateOptions": {"Max", "MaxAwaitTime", "Min", "NoCursorTimeout", "OplogReplay", "ReturnKey",
			"ShowRecordID"},
		"FindOneOptions": {"Max", "MaxAwaitTime", "Min", "NoCursorTimeout", "OplogReplay", "ReturnKey",
			"ShowRecordID"},
		"FindOptions": {"CursorType", "Max", "MaxAwaitTime", "Min", "NoCursorTimeout", "OplogReplay", "ReturnKey",
			"ShowRecordID"},
		"IndexOptions": {"Background", "BucketSize", "Sparse", "StorageEngine"},
	}

	// TODO: options.Find().SetCursorType(options.TailableAwait), only form of FindOptions not supported as of 4/22

	restrictedStages := []string{"$currentOp", "$indexStats", "$listLocalSessions", "$listSessions", "$planCacheStats", "$search"}
	/*
		restrictedOperators := map[string][]string{
			"$group":   {"$sum", "$avg"},
			"$project": {"$add", "$multiply"},
		}
	*/

	/*
		if isPkgDotFunction(pass, call, mongoCollection, "Aggregate") {
			fmt.Println("\nPreorder, Aggregate")
			if hasRestrictedStages(pass, call, restrictedStages, restrictedOperators) {
				pass.Reportf(call.Pos(), "usage of restricted pipeline stages detected")
			}
		}
	*/

	callExprNodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	// loop over all function calls and check against unstableFunctions map
	inspect.Preorder(callExprNodeFilter, func(node ast.Node) {
		call := node.(*ast.CallExpr)
		for pkg, driverFnMap := range unstableFunctions {
			for driverType, fnNames := range driverFnMap {
				for _, fnName := range fnNames {
					fullPkg := pkg + "." + driverType
					if isPkgDotFunction(pass, call, fullPkg, fnName) {
						pass.Reportf(call.Pos(), "use of %v.%v is not supported by the MongoDB Stable API", driverType, fnName)
					}
				}
			}
		}
	})

	basicLitNodeFilter := []ast.Node{
		(*ast.BasicLit)(nil),
	}

	// looks for any of the restricted aggregation stages as a string
	inspect.Preorder(basicLitNodeFilter, func(node ast.Node) {
		switch x := node.(type) {
		case *ast.BasicLit:
			if x.Kind == token.STRING {
				//fmt.Printf("string value: %v\n", x.Value)
				for _, target := range restrictedStages {
					if strings.Contains(x.Value, target) {
						pass.Reportf(x.Pos(), "Aggregation stage '%s' is not supported by the MongoDB Stable API", target)
					}
				}
			}
		}
	})

	selectorExprNodeFilter := []ast.Node{
		(*ast.SelectorExpr)(nil),
	}

	inspect.Preorder(selectorExprNodeFilter, func(node ast.Node) {
		selExpr := node.(*ast.SelectorExpr)

		xIdent, ok := selExpr.X.(*ast.Ident)
		if !ok {
			return
		}

		if xIdent.Name != "CursorType" {
			return
		}

		switch selExpr.Sel.Name {
		case "Tailable", "TailableAwait":
			pass.Reportf(node.Pos(), "Usage of CursorType.%s", selExpr.Sel.Name)
		}
	})

	compositeLitNodeFilter := []ast.Node{
		(*ast.CompositeLit)(nil),
	}

	inspect.Preorder(compositeLitNodeFilter, func(n ast.Node) {
		compLit, ok := n.(*ast.CompositeLit)
		if !ok {
			return
		}

		packageName, structName, ok := getStructInfo(pass, compLit.Type)
		if !ok {
			return
		}

		if packageName != optsPkgName {
			return
		}

		members, ok := unstableOptionsStructs[structName]
		if !ok {
			return
		}

		for _, elt := range compLit.Elts {
			if kv, ok := elt.(*ast.KeyValueExpr); ok {
				if ident, ok := kv.Key.(*ast.Ident); ok {
					for _, member := range members {
						if ident.Name == member {
							pass.Reportf(n.Pos(), "%s.%s is not supported by the MongoDB Stable API", structName, member)
							break
						}
					}
				}
			}
		}
	})

	return nil, nil
}

func getStructInfo(pass *analysis.Pass, expr ast.Expr) (string, string, bool) {
	typ := pass.TypesInfo.TypeOf(expr)
	if typ == nil {
		return "", "", false
	}

	if named, ok := typ.(*types.Named); ok {
		obj := named.Obj()
		if pkg := obj.Pkg(); pkg != nil {
			return pkg.Path(), obj.Name(), true
		}
	}

	return "", "", false
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
		return false
	}

	typ := pass.TypesInfo.TypeOf(selExpr.X)
	if typ == nil {
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
