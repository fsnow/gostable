package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
)

func main() {
	// Check if the file path argument is provided
	if len(os.Args) < 2 {
		fmt.Println("Please provide the file path of a Go file as an argument.")
		return
	}

	// Get the file path from the command-line argument
	filePath := os.Args[1]

	// Read the code file
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, filePath, nil, parser.AllErrors)
	if err != nil {
		fmt.Println("Error parsing file:", err)
		return
	}

	// Print the AST
	ast.Print(fileSet, file)
}
