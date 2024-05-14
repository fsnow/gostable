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

		// Look for any of the restricted aggregation stages as a string.
		// The strings are sufficiently specific that this is not likely to have false positives,
		// i.e. same string but not the MongoDB agg stage.
		case *ast.BasicLit:
			if x.Kind == token.STRING {
				//fmt.Printf("string value: %v\n", x.Value)
				for _, target := range restrictedStages {
					if strings.Contains(x.Value, target) {
						pass.Reportf(x.Pos(), "Aggregation stage '%s' is not supported by the MongoDB Stable API", target)
					}
				}
			}

		// Look at all function calls
		case *ast.CallExpr:
			call := node.(*ast.CallExpr)
			callPkgName, callFnName := pkgPathDotTypeAndFunction(pass, call)
			// Make a general warning about direct use of RunCommand. We might not catch all possible unsupported command constructions.
			if (callPkgName == fullClientPkg || callPkgName == fullDbPkg) && callFnName == "RunCommand" {
				pass.Reportf(call.Pos(), "Any use of RunCommand should be reviewed against the MongoDB Stable API command list")
				// and also try to find the actual command passed to RunCommand
				analyzeRunCommand(pass, call, stack)
			} else {
				// Check against unstableFunctions map
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

		// Look for any unsupported struct fields.
		// We will assume that simply referring to them or setting them is unsupported.
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

		// Look for unsupported cursor types
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

	// debugging, print out stack
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
										} else {
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
