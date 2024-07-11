package ast

import "github.com/FabianRolfMatthiasNoll/Golang-CSharp-Compiler/src/lexer"

type BlockStmt struct {
	Body []Stmt
}

func (stmt BlockStmt) stmt() {}

// To enforce semicolons really easily
type ExpressionStmt struct {
	Expression Expr
}

func (stmt ExpressionStmt) stmt() {}

type VarDeclStmt struct {
	Identifier string
	Value      Expr
	Modifier   lexer.TokenKind
	Type       lexer.TokenKind
}

func (stmt VarDeclStmt) stmt() {}
