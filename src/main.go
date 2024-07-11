package main

import (
	"os"

	"github.com/FabianRolfMatthiasNoll/Golang-CSharp-Compiler/src/lexer"
	"github.com/FabianRolfMatthiasNoll/Golang-CSharp-Compiler/src/parser"
	"github.com/sanity-io/litter"
)

func main() {
	bytes, _ := os.ReadFile("./examples/02.lang")

	tokens := lexer.Tokenize(string(bytes))
	
	for _, token := range tokens {
		token.Debug()
	}

	ast := parser.Parse(tokens)
	litter.Dump(ast)
}
