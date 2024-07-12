package main

import (
	"fmt"
	"os"

	"github.com/FabianRolfMatthiasNoll/Golang-CSharp-Compiler/src/lexer"
	"github.com/FabianRolfMatthiasNoll/Golang-CSharp-Compiler/src/parser"
)

func main() {
	bytes, _ := os.ReadFile("./examples/controlFlow.lang")

	tokens := lexer.Tokenize(string(bytes))

	// for _, token := range tokens {
	// 	token.Debug()
	// }

	ast := parser.Parse(tokens)
	fmt.Println(ast)
}
