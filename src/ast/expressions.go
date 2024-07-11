package ast

import "github.com/FabianRolfMatthiasNoll/Golang-CSharp-Compiler/src/lexer"

// --------------------------------------------
// LITERAL EXPRESSIONS
// --------------------------------------------

type NumberExpr struct {
	Value float64
}

func (expr NumberExpr) expr() {}

type StringExpr struct {
	Value string
}

func (expr StringExpr) expr() {}

type IdentifierExpr struct {
	Name string
}

func (expr IdentifierExpr) expr() {}

// --------------------------------------------
// COMPLEX EXPRESSIONS
// --------------------------------------------

type BinaryExpr struct {
	Left     Expr
	Operator lexer.Token
	Right    Expr
}

func (expr BinaryExpr) expr() {}

type PrefixExpr struct {
	Operator lexer.Token
	Expression    Expr
}

func (expr PrefixExpr) expr() {}

type AssignmentExpr struct {
	Assigne  Expr
	Operator lexer.Token
	Value    Expr
}

func (expr AssignmentExpr) expr() {}
