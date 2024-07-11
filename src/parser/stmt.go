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

func parseVarDeclStmt(p *parser) ast.Stmt {
	typeOrModifier := p.advance()

	modifier := lexer.PUBLIC // Default modifier is public
	dataType := lexer.VAR

	// TODO: Implement more Types etc.
	if typeOrModifier.Kind != lexer.VAR {
		modifier = typeOrModifier.Kind
		dataType = p.advance().Kind
	} else {
		dataType = typeOrModifier.Kind
	}

	identifier := p.expectError(lexer.IDENTIFIER, "Expected identifier after type declaration").Value

	p.expect(lexer.ASSIGNMENT)

	assignedValue := parseExpression(p, ASSIGNMENT)

	p.expect(lexer.SEMICOLON)

	return ast.VarDeclStmt{
		Identifier: identifier,
		Value:      assignedValue,
		Modifier:   modifier,
		Type:       dataType,
	}
}
