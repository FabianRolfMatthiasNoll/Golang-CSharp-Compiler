package parser

import (
	"github.com/FabianRolfMatthiasNoll/Golang-CSharp-Compiler/src/ast"
	"github.com/FabianRolfMatthiasNoll/Golang-CSharp-Compiler/src/lexer"
)

type bindingPower int

const (
	DEFAULT bindingPower = iota
	COMMA
	ASSIGNMENT
	LOGICAL
	RELATIONAL
	ADDITIVE
	MULTIPLICATIVE
	UNARY
	CALL
	MEMBER
	PRIMARY
)

type stmtHandler func(p *parser) ast.Stmt
type nudHandler func(p *parser) ast.Expr
type ledHandler func(p *parser, left ast.Expr, bp bindingPower) ast.Expr

type stmtLookup map[lexer.TokenKind]stmtHandler
type nudLookup map[lexer.TokenKind]nudHandler
type ledLookup map[lexer.TokenKind]ledHandler
type bpLookup map[lexer.TokenKind]bindingPower

var bpTable = bpLookup{}
var nudTable = nudLookup{}
var ledTable = ledLookup{}
var stmtTable = stmtLookup{}

func led(kind lexer.TokenKind, bp bindingPower, led_fn ledHandler) {
	ledTable[kind] = led_fn
	bpTable[kind] = bp
}

func nud(kind lexer.TokenKind, nud_fn nudHandler) {
	nudTable[kind] = nud_fn
}

func stmt(kind lexer.TokenKind, stmt_fn stmtHandler) {
	stmtTable[kind] = stmt_fn
	bpTable[kind] = DEFAULT
}

func createTokenLookups() {

	led(lexer.ASSIGNMENT, ASSIGNMENT, parseAssignmentExpr)
	led(lexer.PLUS_EQUALS, ASSIGNMENT, parseAssignmentExpr)
	led(lexer.MINUS_EQUALS, ASSIGNMENT, parseAssignmentExpr)
	led(lexer.MULTIPLY_EQUALS, ASSIGNMENT, parseAssignmentExpr)
	led(lexer.DIVIDE_EQUALS, ASSIGNMENT, parseAssignmentExpr)
	led(lexer.MODULUS_EQUALS, ASSIGNMENT, parseAssignmentExpr)

	// Logical
	led(lexer.AND, LOGICAL, parseBinaryExpr)
	led(lexer.OR, LOGICAL, parseBinaryExpr)

	// Relational
	led(lexer.EQUALS, RELATIONAL, parseBinaryExpr)
	led(lexer.NOT_EQUALS, RELATIONAL, parseBinaryExpr)
	led(lexer.LESS_THAN, RELATIONAL, parseBinaryExpr)
	led(lexer.LESS_THAN_OR_EQUAL, RELATIONAL, parseBinaryExpr)
	led(lexer.GREATER_THAN, RELATIONAL, parseBinaryExpr)
	led(lexer.GREATER_THAN_OR_EQUAL, RELATIONAL, parseBinaryExpr)

	// Additive & Multiplicative
	led(lexer.PLUS, ADDITIVE, parseBinaryExpr)
	led(lexer.MINUS, ADDITIVE, parseBinaryExpr)

	led(lexer.MULTIPLY, MULTIPLICATIVE, parseBinaryExpr)
	led(lexer.DIVIDE, MULTIPLICATIVE, parseBinaryExpr)
	led(lexer.MODULUS, MULTIPLICATIVE, parseBinaryExpr)

	// Literals & Symbols
	nud(lexer.NUMBER, parsePrimaryExpr)
	nud(lexer.STRING, parsePrimaryExpr)
	nud(lexer.IDENTIFIER, parsePrimaryExpr)
	nud(lexer.NULL, parsePrimaryExpr)

	nud(lexer.MINUS, parsePrefixExpr)
	nud(lexer.OPEN_PAREN, parseGroupedExpr)
	nud(lexer.THIS, parseThisExpr)
	nud(lexer.TRUE, parseBooleanExpr)
	nud(lexer.FALSE, parseBooleanExpr)

	// Handle member access and method calls
	led(lexer.DOT, MEMBER, parseMemberAccessOrMethodCall)

	// Statements
	stmt(lexer.VAR, parseVarDeclStmt)
	stmt(lexer.INT, parseVarDeclStmt)
	stmt(lexer.STRING_TYPE, parseVarDeclStmt)
	stmt(lexer.BOOL, parseVarDeclStmt)
	stmt(lexer.CHAR, parseVarDeclStmt)
	stmt(lexer.DOUBLE, parseVarDeclStmt)
	stmt(lexer.FLOAT, parseVarDeclStmt)

	stmt(lexer.PUBLIC, parseVarDeclStmt)
	stmt(lexer.PRIVATE, parseVarDeclStmt)
	stmt(lexer.PROTECTED, parseVarDeclStmt)
	stmt(lexer.STATIC, parseVarDeclStmt)
	stmt(lexer.FINAL, parseVarDeclStmt)
	stmt(lexer.CLASS, parseClassDeclStmt)
	stmt(lexer.RETURN, parseReturnStmt)
}
