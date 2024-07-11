package parser

import (
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
		modifiers = append(modifiers, ast.Modifier{Kind: lexer.PUBLIC})
	}

	return modifiers
}

func parseType(p *parser) ast.Type {
	token := p.advance()
	return ast.Type{Name: token.Value, Line: token.Line, Column: token.Column}
}
