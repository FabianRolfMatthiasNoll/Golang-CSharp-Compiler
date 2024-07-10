package parser

import (
	"github.com/FabianRolfMatthiasNoll/Golang-CSharp-Compiler/src/ast"
	"github.com/FabianRolfMatthiasNoll/Golang-CSharp-Compiler/src/lexer"
)

type parser struct {
	// TODO: Add a field for errors to keep parsing in the future
	tokens []lexer.Token
	pos    int
}

func createParser(tokenstream []lexer.Token) *parser {
	createTokenLookups()
	return &parser{tokenstream, 0}
}

func Parse(tokenstream []lexer.Token) ast.BlockStmt {
	Body := make([]ast.Stmt, 0)
	p := createParser(tokenstream)

	for p.hasTokensLeft() {
		Body = append(Body, parseStatement(p))
	}

	return ast.BlockStmt{Body: Body}
}

// HELPER METHODS

func (p *parser) currentToken() lexer.Token {
	return p.tokens[p.pos]
}

func (p *parser) currentTokenKind() lexer.TokenKind {
	return p.tokens[p.pos].Kind
}

// func (p *parser) nextToken() lexer.Token {
// 	p.pos++
// 	return p.tokens[p.pos]
// }

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

	if kind != expectedKind {
		if err == nil {
			panic("Expected " + lexer.TokenKindString(expectedKind) + " but got " + lexer.TokenKindString(kind))
		} else {
			panic(err)
		}
	}

	return p.advance()
}

func (p *parser) expect(expectedKind lexer.TokenKind) lexer.Token {
	return p.expectError(expectedKind, nil)
}
