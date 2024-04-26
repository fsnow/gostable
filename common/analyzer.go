package common

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const mongoPkgName = "go.mongodb.org/mongo-driver/mongo"

const fullClientPkg = mongoPkgName + ".Client"
const fullDbPkg = mongoPkgName + ".Database"

const optsPkgName = "go.mongodb.org/mongo-driver/mongo/options"

var StableAnalyzer = &analysis.Analyzer{
	Name: "gostable",
	Doc:  "ensures that all MongoDB go driver code adheres to the Stable API",
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

var unstableFunctions = map[string]map[string][]string{
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

var unstableOptionsStructs = map[string][]string{
	"CreateCollectionOptions": {"Capped", "DefaultIndexOptions", "MaxDocuments", "SizeInBytes", "StorageEngine"},
	"CursorType":              {"Tailable", "TailableAwait"},
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

var restrictedStages = []string{"$currentOp", "$indexStats", "$listLocalSessions", "$listSessions", "$planCacheStats", "$search"}

// commands that are supported without limitations or caveats
var stableCommands = []string{"count", "abortTransaction", "authenticate", "bulkWrite", "collMod", "commitTransaction",
	"delete", "drop", "dropDatabase", "dropIndexes", "endSessions", "findAndModify", "getMore", "insert", "hello",
	"killCursors", "listCollections", "listDatabases", "listIndexes", "ping", "refreshSessions", "update",
}

//var semistableCommands = []string{"aggregate", "create", "createIndexes", "explain", "find"}

/*
// The Stable API command list will change a lot less than the complete list of commands.
// I stopped adding at the sharding commands.
// https://www.mongodb.com/docs/manual/reference/command/#sharding-commands

	var unstableCommands = []string{
		// Aggregation Commands
		"distinct", "mapReduce",
		// Geospatial Commands
		"geoSearch",
		// Query and Write Operation Commands
		"resetError",
		// Query Plan Cache Commands
		"planCacheClear", "planCacheClearFilters", "planCacheListFilters", "planCacheSetFilter",
		// Authentication Commands
		"logout",
		// User Management Commands
		"createUser", "dropAllUsersFromDatabase", "dropUser", "grantRolesToUser", "revokeRolesFromUser", "updateUser", "usersInfo",
		// Role Management Commands
		"createRole", "dropRole", "dropAllRolesFromDatabase", "grantPrivilegesToRole", "grantRolesToRole", "invalidateUserCache",
		"revokePrivilegesFromRole", "revokeRolesFromRole", "rolesInfo", "updateRole",
		// Replication Commands
		"applyOps", "replSetAbortPrimaryCatchUp", "replSetFreeze", "replSetGetConfig", "replSetGetStatus", "replSetInitiate",
		"replSetMaintenance", "replSetReconfig", "replSetResizeOplog", "replSetStepDown", "replSetSyncFrom",
	}
*/
func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	/*
		if isPkgDotFunction(pass, call, mongoCollection, "Aggregate") {
			fmt.Println("\nPreorder, Aggregate")
			if hasRestrictedStages(pass, call, restrictedStages, restrictedOperators) {
				pass.Reportf(call.Pos(), "usage of restricted pipeline stages detected")
			}
		}
	*/

	nodeFilter := []ast.Node{
		(*ast.BasicLit)(nil),
		(*ast.CallExpr)(nil),
		(*ast.CompositeLit)(nil),
		(*ast.SelectorExpr)(nil),
	}

	inspect.WithStack(nodeFilter, func(node ast.Node, push bool, stack []ast.Node) bool {
		if push {
			return true
		}
		switch x := node.(type) {

		// look for any of the restricted aggregation stages as a string
		case *ast.BasicLit:
			if x.Kind == token.STRING {
				//fmt.Printf("string value: %v\n", x.Value)
				for _, target := range restrictedStages {
					if strings.Contains(x.Value, target) {
						pass.Reportf(x.Pos(), "Aggregation stage '%s' is not supported by the MongoDB Stable API", target)
					}
				}
			}

		// loop over all function calls and check against unstableFunctions map
		case *ast.CallExpr:
			call := node.(*ast.CallExpr)
			callPkgName, callFnName := pkgPathDotTypeAndFunction(pass, call)
			if (callPkgName == fullClientPkg || callPkgName == fullDbPkg) && callFnName == "RunCommand" {
				pass.Reportf(call.Pos(), "Any use of RunCommand should be reviewed against the MongoDB Stable API command list")
				analyzeRunCommand(pass, call, stack)
			} else {
				for pkg, driverFnMap := range unstableFunctions {
					for driverType, fnNames := range driverFnMap {
						for _, fnName := range fnNames {
							fullPkg := pkg + "." + driverType
							if isPkgDotFunction(pass, call, fullPkg, fnName) {
								pass.Reportf(call.Pos(), "Function %v.%v is not supported by the MongoDB Stable API", driverType, fnName)
							}
						}
					}
				}
			}

		// look for any unsupported struct fields
		case *ast.CompositeLit:
			compLit, ok := node.(*ast.CompositeLit)
			if !ok {
				return false
			}

			packageName, structName, ok := getStructInfo(pass, compLit.Type)
			if !ok {
				return false
			}

			if packageName != optsPkgName {
				return false
			}

			members, ok := unstableOptionsStructs[structName]
			if !ok {
				return false
			}

			for _, elt := range compLit.Elts {
				if kv, ok := elt.(*ast.KeyValueExpr); ok {
					if ident, ok := kv.Key.(*ast.Ident); ok {
						for _, member := range members {
							if ident.Name == member {
								pass.Reportf(ident.Pos(), "Struct field %s.%s is not supported by the MongoDB Stable API", structName, member)
								break
							}
						}
					}
				}
			}

		case *ast.SelectorExpr:
			selExpr := node.(*ast.SelectorExpr)

			switch selExpr.Sel.Name {
			case "Tailable", "TailableAwait":
				xIdent, ok := selExpr.X.(*ast.Ident)
				if !ok {
					return false
				}

				typ := pass.TypesInfo.TypeOf(xIdent)
				if typ == nil {
					return false
				}

				if named, ok := typ.(*types.Named); ok {
					obj := named.Obj()
					pkg := obj.Pkg()
					if pkg == nil || pkg.Path() != optsPkgName {
						return false
					}
				}

				pass.Reportf(node.Pos(), "Struct field CursorType.%s is not supported by the MongoDB Stable API", selExpr.Sel.Name)
			}
		}
		return false
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

func pkgPathDotTypeAndFunction(pass *analysis.Pass, call *ast.CallExpr) (string, string) {
	// Check if the call expression is a selector expression
	selExpr, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return "", ""
	}

	// Check if the selector's X expression is a package identifier
	_, ok = selExpr.X.(*ast.Ident)
	if !ok {
		return "", selExpr.Sel.Name
	}

	typ := pass.TypesInfo.TypeOf(selExpr.X)
	if typ == nil {
		return "", selExpr.Sel.Name
	}

	typStr := typ.String()
	if len(typStr) > 0 && typStr[0] == '*' {
		typStr = typStr[1:]
	}

	return typStr, selExpr.Sel.Name
}

func analyzeRunCommand(pass *analysis.Pass, call *ast.CallExpr, stack []ast.Node) {
	//fmt.Println("In analyzeRunCommand")
	// Get the command argument (second argument)
	if len(call.Args) < 2 {
		return
	}
	cmdArg := call.Args[1]

	// Check if the command argument is a bson.D literal
	if bsonDLit, ok := cmdArg.(*ast.CompositeLit); ok {
		//fmt.Println("CompositeLit")
		if isBsonDType(pass.TypesInfo.TypeOf(bsonDLit)) {
			// Analyze the bson.D literal
			//fmt.Println("analyzeCommandLiteral 1")
			analyzeCommandLiteral(pass, bsonDLit)
			return
		}
	} else {
		// Check if the command argument is a variable
		if ident, ok := cmdArg.(*ast.Ident); ok {
			//fmt.Println("Ident")
			if isBsonDType(pass.TypesInfo.ObjectOf(ident).Type()) {
				//fmt.Println("isBsonDType")
				// Find the variable declaration and analyze its value
				if assignStmt := findVariableAssignment(pass, ident, stack); assignStmt != nil {
					//fmt.Println("assignStmt != nil")
					if bsonDLit, ok := assignStmt.Rhs[0].(*ast.CompositeLit); ok {
						//fmt.Println("analyzeCommandLiteral 2")
						//fmt.Printf("bsonDLit %v\n", bsonDLit.Type)

						analyzeCommandLiteral(pass, bsonDLit)
					}
				}
			}
		}
	}
}

func isBsonDType(typ types.Type) bool {
	//fmt.Printf("typ: %v\n", typ.String())
	return typ.String() == "bson.D" || typ.String() == "go.mongodb.org/mongo-driver/bson/primitive.D"
}

func findVariableAssignment(pass *analysis.Pass, ident *ast.Ident, stack []ast.Node) *ast.AssignStmt {
	var assignStmt *ast.AssignStmt

	/*
		fmt.Println("Stack")
		for _, stackNode := range stack {
			fmt.Printf("%T\n", stackNode)
			//ast.Print(nil, stackNode)
		}
		fmt.Println("End Stack")

		for i := len(stack) - 1; i >= 0; i-- {
			fmt.Printf("%T\n", stack[i])
		}
		fmt.Println("End Reverse")
	*/

	// Search for the assignment, going up the call stack until found
	for i := len(stack) - 1; i >= 0; i-- {
		//fmt.Printf("%T\n", stack[i])

		ast.Inspect(stack[i], func(node ast.Node) bool {
			switch stmt := node.(type) {
			case *ast.AssignStmt:
				for _, lhs := range stmt.Lhs {
					if lhsIdent, ok := lhs.(*ast.Ident); ok && lhsIdent.Name == ident.Name {
						assignStmt = stmt
						return false
					}
				}
			}
			return true
		})

		if assignStmt != nil {
			break
		}
	}

	// The previous loop ended with the ast.File containing the runCommand call.
	// If an assignment of the command variable was not found in that file, the following
	// looks for that declaration in all Files. Note that this might find some false positives.
	if assignStmt == nil {
		for _, file := range pass.Files {
			ast.Inspect(file, func(node ast.Node) bool {
				switch stmt := node.(type) {
				case *ast.AssignStmt:
					for _, lhs := range stmt.Lhs {
						if lhsIdent, ok := lhs.(*ast.Ident); ok && lhsIdent.Name == ident.Name {
							assignStmt = stmt
							return false
						}
					}
				}
				return true
			})

			if assignStmt != nil {
				break
			}
		}
	}

	return assignStmt
}

func findVariableAssignment_save(pass *analysis.Pass, ident *ast.Ident) *ast.AssignStmt {
	var assignStmt *ast.AssignStmt

	// Find the function declaration containing the identifier
	funcDecl := findFunctionDeclaration(pass, ident)

	// First, search for the assignment statement within the function
	if funcDecl != nil {
		ast.Inspect(funcDecl, func(node ast.Node) bool {
			switch stmt := node.(type) {
			case *ast.AssignStmt:
				for _, lhs := range stmt.Lhs {
					if lhsIdent, ok := lhs.(*ast.Ident); ok && lhsIdent.Name == ident.Name {
						assignStmt = stmt
						return false
					}
				}
			}
			return true
		})
	}

	// If the assignment statement is not found within the function, search in outer scopes
	if assignStmt == nil {
		for _, file := range pass.Files {
			ast.Inspect(file, func(node ast.Node) bool {
				switch stmt := node.(type) {
				case *ast.AssignStmt:
					for _, lhs := range stmt.Lhs {
						if lhsIdent, ok := lhs.(*ast.Ident); ok && lhsIdent.Name == ident.Name {
							assignStmt = stmt
							return false
						}
					}
				}
				return true
			})

			if assignStmt != nil {
				break
			}
		}
	}

	return assignStmt
}

func findFunctionDeclaration(pass *analysis.Pass, ident *ast.Ident) *ast.FuncDecl {
	var funcDecl *ast.FuncDecl

	// Find the function declaration containing the identifier
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			if fd, ok := node.(*ast.FuncDecl); ok {
				if containsIdent(fd, ident) {
					funcDecl = fd
					return false
				}
			}
			return true
		})

		if funcDecl != nil {
			break
		}
	}

	return funcDecl
}

func containsIdent(node ast.Node, ident *ast.Ident) bool {
	var found bool
	ast.Inspect(node, func(n ast.Node) bool {
		if id, ok := n.(*ast.Ident); ok && id.Name == ident.Name {
			found = true
			return false
		}
		return true
	})
	return found
}

func analyzeCommandLiteral(pass *analysis.Pass, x *ast.CompositeLit) {
	//fmt.Println("Inside analyzeCommandLiteral")

	if x.Type != nil {
		switch t := x.Type.(type) {
		case *ast.SelectorExpr:
			//fmt.Printf("Sel.Name: %v\n", t.Sel.Name)
			if t.Sel.Name == "D" {
				if pkg, ok := t.X.(*ast.Ident); ok && pkg.Name == "bson" {
					// Access the bson.D type
					bsonDType := pass.TypesInfo.TypeOf(x)
					//fmt.Println("ACL 6")
					// Check if the bsonDType is *types.Named
					if named, ok := bsonDType.(*types.Named); ok {
						//fmt.Println("ACL 5")
						if named.Obj().Name() == "D" && named.Obj().Pkg().Name() == "primitive" {
							//fmt.Println("ACL 4")
							// Get the first element of the bson.D slice
							if len(x.Elts) > 0 {
								//fmt.Println("ACL 2")
								if compositeElt, ok := x.Elts[0].(*ast.CompositeLit); ok {
									//fmt.Println("ACL 1")
									// Extract the command name from the composite literal
									commandName := getCommandName(compositeElt)
									//fmt.Printf("commandName: %v\n", commandName)
									if commandName != "" {
										if isStableCommand(commandName) {
											//fmt.Printf("Stable: %v\n", commandName)
											// Perform further analysis on the value if needed
											// You can customize this based on your specific requirements
										} else {
											// Report a violation if the command is not stable
											pass.Reportf(compositeElt.Pos(), "Command %s is not supported by the MongoDB Stable API", commandName)
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

func getCommandName(compositeElt *ast.CompositeLit) string {
	//ast.Print(nil, compositeElt)
	if len(compositeElt.Elts) > 0 {
		//fmt.Printf("type: %T\n", compositeElt.Elts[0])
		if kv, ok := compositeElt.Elts[0].(*ast.KeyValueExpr); ok {
			if ident, ok := kv.Key.(*ast.Ident); ok {
				if ident.Name == "Key" {
					if val, ok := kv.Value.(*ast.BasicLit); ok {
						if val.Kind == token.STRING {
							return stripQuotes(val.Value)
						}
					}
				}
			}
		}
	}
	return ""
}

func stripQuotes(s string) string {
	if len(s) > 0 && s[0] == '"' {
		s = s[1:]
	}
	if len(s) > 0 && s[len(s)-1] == '"' {
		s = s[:len(s)-1]
	}
	return s
}

func isStableCommand(cmd string) bool {
	// Check if the command is part of the MongoDB Stable API
	// You can customize this based on your specific requirements

	for _, stableCmd := range stableCommands {
		if cmd == stableCmd {
			return true
		}
	}

	return false
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

// ---------------------------------------------------
// Might use this stuff later...
// ---------------------------------------------------

func hasRestrictedStages(pass *analysis.Pass, call *ast.CallExpr, restrictedStages []string, restrictedOperators map[string][]string) bool {
	for _, arg := range call.Args {
		var bsonDocs []*ast.CompositeLit

		switch arg := arg.(type) {
		case *ast.Ident:
			bsonDocs = findPipelineVariable(pass, arg.Name, call.Pos())
			//fmt.Println("ast.Ident case")
			//fmt.Printf("len: %v\n", len(bsonDocs))
		case *ast.CompositeLit:
			if arg.Type == nil {
				bsonDocs = extractBSONDocs(arg.Elts)
				//fmt.Println("ast.CompositeLit case, if")
			} else if ident, ok := arg.Type.(*ast.ArrayType); ok && ident.Elt != nil {
				if _, ok := ident.Elt.(*ast.StructType); ok {
					bsonDocs = extractBSONDocs(arg.Elts)
					//fmt.Println("ast.CompositeLit case, else if if")
				}
			}
		case *ast.CallExpr:
			if sel, ok := arg.Fun.(*ast.SelectorExpr); ok && sel.Sel.Name == "Pipeline" {
				bsonDocs = extractBSONDocsFromPipeline(arg.Args)
				//fmt.Println("ast.CallExpr case")
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

	//fmt.Printf("findPipelineVariable, varName=%v\n", varName)
	//fmt.Printf("pass.Files len: %v\n", len(pass.Files))

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
								//fmt.Println("before append 1")
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
				//fmt.Println("before append 2")
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
