package main

import (
	"fmt"
	"os"

	"github.com/FabianRolfMatthiasNoll/Golang-CSharp-Compiler/src/lexer"
	"github.com/FabianRolfMatthiasNoll/Golang-CSharp-Compiler/src/parser"
	"github.com/FabianRolfMatthiasNoll/Golang-CSharp-Compiler/src/typecheck"
)

func main() {
	bytes, _ := os.ReadFile("./examples/standardTypes.lang")

	tokens := lexer.Tokenize(string(bytes))

	// for _, token := range tokens {
	// 	token.Debug()
	// }

	ast := parser.Parse(tokens)
	fmt.Println(ast)

	fmt.Println("=========================================")
	fmt.Println("Type checking...")
	fmt.Println("=========================================")

	typeChecker := typecheck.NewTypeChecker()
	typedAst := typeChecker.CheckProgram(&ast)
	fmt.Println(typedAst)

}
