package parser

import (
	"github.com/FabianRolfMatthiasNoll/Golang-CSharp-Compiler/src/ast"
	"github.com/FabianRolfMatthiasNoll/Golang-CSharp-Compiler/src/lexer"
)

func parseStatement(p *parser) ast.Stmt {
	stmt_fn, exists := stmtTable[p.currentTokenKind()]
	if exists {
		return stmt_fn(p)
	}

	expression := parseExpression(p, DEFAULT)
	p.expect(lexer.SEMICOLON)

	return ast.ExpressionStmt{Expression: expression}
}
