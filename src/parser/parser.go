package parser

import (
	"fmt"

	"github.com/FabianRolfMatthiasNoll/Golang-CSharp-Compiler/src/ast"
	"github.com/FabianRolfMatthiasNoll/Golang-CSharp-Compiler/src/lexer"
)

type parser struct {
	tokens []lexer.Token
	pos    int
}

func createParser(tokenstream []lexer.Token) *parser {
	createTokenLookups()
	return &parser{tokenstream, 0}
}

func Parse(tokenstream []lexer.Token) ast.Program {
	classes := make([]ast.ClassDeclStmt, 0)
	p := createParser(tokenstream)

	for p.hasTokensLeft() {
		classStmt := parseClassDeclStmt(p)
		if class, ok := classStmt.(ast.ClassDeclStmt); ok {
			classes = append(classes, class)
		} else {
			panic(fmt.Sprintf("Expected ClassDeclStmt but got %T", classStmt))
		}
	}

	return ast.Program{Classes: classes}
}

// HELPER METHODS

func (p *parser) currentToken() lexer.Token {
	return p.tokens[p.pos]
}

func (p *parser) currentTokenKind() lexer.TokenKind {
	return p.tokens[p.pos].Kind
}

func (p *parser) advance() lexer.Token {
	tk := p.currentToken()
	p.pos++
	return tk
}

func (p *parser) hasTokensLeft() bool {
	return p.pos < len(p.tokens) && p.currentTokenKind() != lexer.EOF
}

func (p *parser) expectError(expectedKind lexer.TokenKind, err any) lexer.Token {
	kind := p.currentTokenKind()
	token := p.currentToken()

	if kind != expectedKind {
		if err == nil {
			panic(fmt.Sprintf("Expected %s but got %s at line %d, column %d", lexer.TokenKindString(expectedKind), lexer.TokenKindString(kind), token.Line, token.Column))
		} else {
			panic(fmt.Sprintf("%v at line %d, column %d", err, token.Line, token.Column))
		}
	}

	return p.advance()
}

func (p *parser) expect(expectedKind lexer.TokenKind) lexer.Token {
	return p.expectError(expectedKind, nil)
}

func (p *parser) nextTokenKind() lexer.TokenKind {
	if p.pos+1 < len(p.tokens) {
		return p.tokens[p.pos+1].Kind
	}
	return lexer.EOF
}
