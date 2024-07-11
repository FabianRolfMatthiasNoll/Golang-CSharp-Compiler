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
		return ast.NumberExpr{Value: number, Line: p.currentToken().Line, Column: p.currentToken().Column}
	case lexer.STRING:
		return ast.StringExpr{Value: p.advance().Value, Line: p.currentToken().Line, Column: p.currentToken().Column}
	case lexer.IDENTIFIER:
		token := p.advance()
		var expr ast.Expr = ast.IdentifierExpr{Name: token.Value, Line: token.Line, Column: token.Column}
		if p.currentTokenKind() == lexer.OPEN_PAREN {
			return parseMethodCallExpr(p, ast.ThisExpr{Line: token.Line, Column: token.Column}, token.Value)
		}
		for p.currentTokenKind() == lexer.DOT {
			p.advance() // consume the dot
			member := p.expect(lexer.IDENTIFIER).Value
			if p.currentTokenKind() == lexer.OPEN_PAREN {
				return parseMethodCallExpr(p, expr, member)
			}
			expr = ast.MemberAccessExpr{
				Receiver: expr,
				Member:   member,
				Line:     p.currentToken().Line,
				Column:   p.currentToken().Column,
			}
		}
		return expr
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

func parseMemberAccessOrMethodCall(p *parser, receiver ast.Expr, bp bindingPower) ast.Expr {
	memberName := p.expect(lexer.IDENTIFIER).Value
	if p.currentTokenKind() == lexer.OPEN_PAREN {
		return parseMethodCallExpr(p, receiver, memberName)
	}
	return ast.MemberAccessExpr{
		Receiver: receiver,
		Member:   memberName,
		Line:     p.currentToken().Line,
		Column:   p.currentToken().Column,
	}
}

func parseMethodCallExpr(p *parser, receiver ast.Expr, methodName string) ast.Expr {
	p.expect(lexer.OPEN_PAREN)
	args := parseArguments(p)
	p.expect(lexer.CLOSE_PAREN)

	return ast.MethodCallExpr{
		Receiver:   receiver,
		MethodName: methodName,
		Args:       args,
		Line:       p.currentToken().Line,
		Column:     p.currentToken().Column,
	}
}

func parseArguments(p *parser) []ast.Expr {
	args := []ast.Expr{}

	if p.currentTokenKind() == lexer.CLOSE_PAREN {
		return args
	}

	for {
		args = append(args, parseExpression(p, DEFAULT))

		if p.currentTokenKind() != lexer.COMMA {
			break
		}

		p.advance()
	}

	return args
}

func parseThisExpr(p *parser) ast.Expr {
	token := p.advance()
	var expr ast.Expr = ast.ThisExpr{Line: token.Line, Column: token.Column}
	for p.currentTokenKind() == lexer.DOT {
		p.advance() // consume the dot
		member := p.expect(lexer.IDENTIFIER).Value
		if p.currentTokenKind() == lexer.OPEN_PAREN {
			return parseMethodCallExpr(p, expr, member)
		}
		expr = ast.MemberAccessExpr{
			Receiver: expr,
			Member:   member,
			Line:     p.currentToken().Line,
			Column:   p.currentToken().Column,
		}
	}
	if p.currentTokenKind() == lexer.OPEN_PAREN {
		return parseMethodCallExpr(p, expr, "this")
	}
	return expr
}

func parseBooleanExpr(p *parser) ast.Expr {
	tokenValue := p.advance().Value
	return ast.BoolExpr{Value: tokenValue == "true", Line: p.currentToken().Line, Column: p.currentToken().Column}
}
