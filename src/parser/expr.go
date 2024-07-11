package parser

import (
	"fmt"
	"strconv"

	"github.com/FabianRolfMatthiasNoll/Golang-CSharp-Compiler/src/ast"
	"github.com/FabianRolfMatthiasNoll/Golang-CSharp-Compiler/src/lexer"
)

func parseExpression(p *parser, bp bindingPower) ast.Expr {
	// Lookup if a function exists for the current token kind
	tokenKind := p.currentTokenKind()
	nud_fn, exists := nudTable[tokenKind]

	if !exists {
		panic(fmt.Sprintf("NUD Handler expected for token %s\n", lexer.TokenKindString(tokenKind)))
	}

	left := nud_fn(p)

	// While we have a LED and current BP is < BP of current token
	// continue parsing to the right side
	for bpTable[p.currentTokenKind()] > bp {
		tokenKind = p.currentTokenKind()
		led_fn, exists := ledTable[tokenKind]

		if !exists {
			panic(fmt.Sprintf("LED Handler expected for token %s\n", lexer.TokenKindString(tokenKind)))
		}

		left = led_fn(p, left, bpTable[tokenKind])
	}

	return left
}

func parsePrimaryExpr(p *parser) ast.Expr {
	switch p.currentTokenKind() {
	case lexer.NUMBER:
		number, _ := strconv.ParseFloat(p.advance().Value, 64)
		return ast.NumberExpr{Value: number}
	case lexer.STRING:
		return ast.StringExpr{Value: p.advance().Value}
	case lexer.IDENTIFIER:
		return ast.IdentifierExpr{Name: p.advance().Value}
	default:
		panic(fmt.Sprintf("Cannot create primary expression from token %s\n", lexer.TokenKindString(p.currentTokenKind())))
	}
}

func parseBinaryExpr(p *parser, left ast.Expr, bp bindingPower) ast.Expr {
	operatorToken := p.advance()
	right := parseExpression(p, bp)

	return ast.BinaryExpr{Left: left, Operator: operatorToken, Right: right}
}

func parsePrefixExpr(p *parser) ast.Expr {
	operatorToken := p.advance()
	expression := parseExpression(p, DEFAULT)

	return ast.PrefixExpr{Operator: operatorToken, Expression: expression}
}

func parseAssignmentExpr(p *parser, left ast.Expr, bp bindingPower) ast.Expr {
	operatorToken := p.advance()
	value := parseExpression(p, bp)
	return ast.AssignmentExpr{Assigne: left, Operator: operatorToken, Value: value}
}

func parseGroupedExpr(p *parser) ast.Expr {
	p.advance() // Consume left parenthesis
	expr := parseExpression(p, DEFAULT)
	p.expectError(lexer.CLOSE_PAREN, "Expected closing parenthesis")

	return expr
}
