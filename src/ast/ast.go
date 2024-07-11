package ast

import "github.com/FabianRolfMatthiasNoll/Golang-CSharp-Compiler/src/lexer"

// Base interfaces
type Stmt interface {
	stmt()
	GetLine() int
	GetColumn() int
}

type Expr interface {
	expr()
	GetLine() int
	GetColumn() int
}

// Modifiers and Type
type Modifier struct {
	Kind lexer.TokenKind
}

type Type struct {
	Name string
}

// Program
type Program struct {
	Classes []ClassDeclStmt
}

// Statements
type BlockStmt struct {
	Body   []Stmt
	Line   int
	Column int
}

func (stmt BlockStmt) stmt() {}
func (stmt BlockStmt) GetLine() int { return stmt.Line }
func (stmt BlockStmt) GetColumn() int { return stmt.Column }

type ExpressionStmt struct {
	Expression Expr
	Line       int
	Column     int
}

func (stmt ExpressionStmt) stmt() {}
func (stmt ExpressionStmt) GetLine() int { return stmt.Line }
func (stmt ExpressionStmt) GetColumn() int { return stmt.Column }

type VarDeclStmt struct {
	Modifiers  []Modifier
	Identifier string
	Type       Type
	Value      Expr
	Line       int
	Column     int
}

func (stmt VarDeclStmt) stmt() {}
func (stmt VarDeclStmt) GetLine() int { return stmt.Line }
func (stmt VarDeclStmt) GetColumn() int { return stmt.Column }

type ReturnStmt struct {
	Value  Expr
	Line   int
	Column int
}

func (stmt ReturnStmt) stmt() {}
func (stmt ReturnStmt) GetLine() int { return stmt.Line }
func (stmt ReturnStmt) GetColumn() int { return stmt.Column }

// Class-related statements
type ClassDeclStmt struct {
	Modifiers []Modifier
	Name      string
	Body      ClassBody
	Line      int
	Column    int
}

func (stmt ClassDeclStmt) stmt() {}
func (stmt ClassDeclStmt) GetLine() int { return stmt.Line }
func (stmt ClassDeclStmt) GetColumn() int { return stmt.Column }

type ClassBody struct {
	Members []ClassMember
	Line    int
	Column  int
}

func (body ClassBody) GetLine() int { return body.Line }
func (body ClassBody) GetColumn() int { return body.Column }

type ClassMember interface {
	classMember()
	GetLine() int
	GetColumn() int
}

type FieldDeclStmt struct {
	Modifiers  []Modifier
	Type       Type
	Identifier string
	Value      Expr
	Line       int
	Column     int
}

func (stmt FieldDeclStmt) classMember() {}
func (stmt FieldDeclStmt) GetLine() int { return stmt.Line }
func (stmt FieldDeclStmt) GetColumn() int { return stmt.Column }

type MethodDeclStmt struct {
	Modifiers  []Modifier
	ReturnType Type
	Name       string
	Parameters []Parameter
	Body       BlockStmt
	Line       int
	Column     int
}

func (stmt MethodDeclStmt) classMember() {}
func (stmt MethodDeclStmt) GetLine() int { return stmt.Line }
func (stmt MethodDeclStmt) GetColumn() int { return stmt.Column }

type ConstructorDeclStmt struct {
	Modifiers  []Modifier
	Name       string
	Parameters []Parameter
	Body       BlockStmt
	Line       int
	Column     int
}

func (stmt ConstructorDeclStmt) classMember() {}
func (stmt ConstructorDeclStmt) GetLine() int { return stmt.Line }
func (stmt ConstructorDeclStmt) GetColumn() int { return stmt.Column }

type Parameter struct {
	Type       Type
	Identifier string
}

// Expressions
type NumberExpr struct {
	Value  float64
	Line   int
	Column int
}

func (expr NumberExpr) expr() {}
func (expr NumberExpr) GetLine() int { return expr.Line }
func (expr NumberExpr) GetColumn() int { return expr.Column }

type StringExpr struct {
	Value  string
	Line   int
	Column int
}

func (expr StringExpr) expr() {}
func (expr StringExpr) GetLine() int { return expr.Line }
func (expr StringExpr) GetColumn() int { return expr.Column }

type IdentifierExpr struct {
	Name   string
	Line   int
	Column int
}

func (expr IdentifierExpr) expr() {}
func (expr IdentifierExpr) GetLine() int { return expr.Line }
func (expr IdentifierExpr) GetColumn() int { return expr.Column }

type BinaryExpr struct {
	Left     Expr
	Operator lexer.Token
	Right    Expr
	Line     int
	Column   int
}

func (expr BinaryExpr) expr() {}
func (expr BinaryExpr) GetLine() int { return expr.Line }
func (expr BinaryExpr) GetColumn() int { return expr.Column }

type PrefixExpr struct {
	Operator    lexer.Token
	Expression  Expr
	Line        int
	Column      int
}

func (expr PrefixExpr) expr() {}
func (expr PrefixExpr) GetLine() int { return expr.Line }
func (expr PrefixExpr) GetColumn() int { return expr.Column }

type AssignmentExpr struct {
	Assigne  Expr
	Operator lexer.Token
	Value    Expr
	Line     int
	Column   int
}

func (expr AssignmentExpr) expr() {}
func (expr AssignmentExpr) GetLine() int { return expr.Line }
func (expr AssignmentExpr) GetColumn() int { return expr.Column }
