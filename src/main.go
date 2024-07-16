package main

import (
	"fmt"
	"os"

	"github.com/FabianRolfMatthiasNoll/Golang-CSharp-Compiler/src/lexer"
	"github.com/FabianRolfMatthiasNoll/Golang-CSharp-Compiler/src/parser"
	"github.com/FabianRolfMatthiasNoll/Golang-CSharp-Compiler/src/typecheck"
)

func main() {
	fmt.Println("=========================================")
	fmt.Println("Reading File...")
	fmt.Println("=========================================")

	bytes, _ := os.ReadFile("./examples/standardTypes.lang")
	fmt.Println(bytes)

	fmt.Println("=========================================")
	fmt.Println("Parsing...")
	fmt.Println("=========================================")

	tokens := lexer.Tokenize(string(bytes))

	for _, token := range tokens {
		token.Debug()
	}

	fmt.Println("=========================================")
	fmt.Println("Parsing...")
	fmt.Println("=========================================")

	ast := parser.Parse(tokens)
	fmt.Println(ast)

	fmt.Println("=========================================")
	fmt.Println("Type checking...")
	fmt.Println("=========================================")

	typeChecker := typecheck.NewTypeChecker()
	typedAst := typeChecker.CheckProgram(&ast)
	fmt.Println(typedAst)

}
