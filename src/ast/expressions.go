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

