package parser

import (
	"fmt"

	"github.com/FabianRolfMatthiasNoll/Golang-CSharp-Compiler/src/ast"
	"github.com/FabianRolfMatthiasNoll/Golang-CSharp-Compiler/src/lexer"
)

var builtInTypes = map[string]bool{
	"int":    true,
	"string": true,
	"float":  true,
	"bool":   true,
	"var":    true,
	"char":   true,
	"double": true,
	"void":   true,
}

func isType(p *parser) bool {
	token := p.currentToken()
	// Assuming built-in types are all lowercase and user-defined types start with an uppercase letter
	if (token.Kind == lexer.IDENTIFIER && token.Value[0] >= 'A' && token.Value[0] <= 'Z') || builtInTypes[token.Value] {
		if p.nextTokenKind() != lexer.OPEN_PAREN {
			return true
		}
	}
	return false
}

func isModifier(kind lexer.TokenKind) bool {
	switch kind {
	case lexer.PUBLIC, lexer.PRIVATE, lexer.PROTECTED, lexer.STATIC, lexer.FINAL:
		return true
	}
	return false
}

func parseModifiers(p *parser) []ast.Modifier {
	modifiers := []ast.Modifier{}

	for isModifier(p.currentTokenKind()) {
		modifiers = append(modifiers, ast.Modifier{Kind: p.advance().Kind})
	}

	if len(modifiers) == 0 {
		modifiers = append(modifiers, ast.Modifier{Kind: lexer.PRIVATE})
	}

	return modifiers
}

func parseType(p *parser) ast.Type {
	token := p.advance()
	return ast.Type{Name: token.Value, Line: token.Line, Column: token.Column}
}

func assignStandardType(dataType ast.Type, assignedValue ast.Expr, p *parser) ast.Expr {
	switch dataType.Name {
	case "int":
		assignedValue = ast.IntLiteralExpr{Value: 0, Line: p.currentToken().Line, Column: p.currentToken().Column}
	case "bool":
		assignedValue = ast.BoolLiteralExpr{Value: false, Line: p.currentToken().Line, Column: p.currentToken().Column}
	case "string":
		assignedValue = ast.NullLiteralExpr{Line: p.currentToken().Line, Column: p.currentToken().Column}
	case "char":
		assignedValue = ast.CharLiteralExpr{Value: 0, Line: p.currentToken().Line, Column: p.currentToken().Column}
	case "void":
		panic(fmt.Sprintf("Cannot assign value to variable of type void at line %d, column %d", p.currentToken().Line, p.currentToken().Column))
	case "var":
		panic(fmt.Sprintf("Cannot use var without assigning a value at line %d, column %d", p.currentToken().Line, p.currentToken().Column))
	default:
		assignedValue = ast.NullLiteralExpr{Line: p.currentToken().Line, Column: p.currentToken().Column}
	}
	return assignedValue
}
