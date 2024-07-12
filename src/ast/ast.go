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
	Name   string
	Line   int
	Column int
}

// Program
type Program struct {
	Classes []ClassDeclStmt
}

// ========================================================================================================
// Statements
// ========================================================================================================

type BlockStmt struct {
	Body   []Stmt
	Line   int
	Column int
}

func (stmt BlockStmt) stmt()          {}
func (stmt BlockStmt) GetLine() int   { return stmt.Line }
func (stmt BlockStmt) GetColumn() int { return stmt.Column }

type ExpressionStmt struct {
	Expression Expr
	Line       int
	Column     int
}

func (stmt ExpressionStmt) stmt()          {}
func (stmt ExpressionStmt) GetLine() int   { return stmt.Line }
func (stmt ExpressionStmt) GetColumn() int { return stmt.Column }

type VarDeclStmt struct {
	Modifiers  []Modifier
	Identifier string
	Type       Type
	Value      Expr
	Line       int
	Column     int
}

func (stmt VarDeclStmt) stmt()          {}
func (stmt VarDeclStmt) GetLine() int   { return stmt.Line }
func (stmt VarDeclStmt) GetColumn() int { return stmt.Column }

type ReturnStmt struct {
	Value  Expr
	Line   int
	Column int
}

func (stmt ReturnStmt) stmt()          {}
func (stmt ReturnStmt) GetLine() int   { return stmt.Line }
func (stmt ReturnStmt) GetColumn() int { return stmt.Column }

// Class-related statements
type ClassDeclStmt struct {
	Modifiers []Modifier
	Name      string
	Body      ClassBody
	Line      int
	Column    int
}

func (stmt ClassDeclStmt) stmt()          {}
func (stmt ClassDeclStmt) GetLine() int   { return stmt.Line }
func (stmt ClassDeclStmt) GetColumn() int { return stmt.Column }

type ClassBody struct {
	Members []ClassMember
	Line    int
	Column  int
}

func (body ClassBody) GetLine() int   { return body.Line }
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

func (stmt FieldDeclStmt) classMember()   {}
func (stmt FieldDeclStmt) GetLine() int   { return stmt.Line }
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

func (stmt MethodDeclStmt) classMember()   {}
func (stmt MethodDeclStmt) GetLine() int   { return stmt.Line }
func (stmt MethodDeclStmt) GetColumn() int { return stmt.Column }

type ConstructorDeclStmt struct {
	Modifiers  []Modifier
	Name       string
	Parameters []Parameter
	Body       BlockStmt
	Line       int
	Column     int
}

func (stmt ConstructorDeclStmt) classMember()   {}
func (stmt ConstructorDeclStmt) GetLine() int   { return stmt.Line }
func (stmt ConstructorDeclStmt) GetColumn() int { return stmt.Column }

type Parameter struct {
	Type       Type
	Identifier string
}

// Control flow statements

type WhileStmt struct {
	Condition Expr
	Body      Stmt
	Line      int
	Column    int
}

func (stmt WhileStmt) stmt()          {}
func (stmt WhileStmt) GetLine() int   { return stmt.Line }
func (stmt WhileStmt) GetColumn() int { return stmt.Column }

type IfStmt struct {
	Condition Expr
	Then      Stmt
	Else      Stmt
	Line      int
	Column    int
}

func (stmt IfStmt) stmt()          {}
func (stmt IfStmt) GetLine() int   { return stmt.Line }
func (stmt IfStmt) GetColumn() int { return stmt.Column }

// ========================================================================================================
// Expressions
// ========================================================================================================

type StringExpr struct {
	Value  string
	Line   int
	Column int
}

func (expr StringExpr) expr()          {}
func (expr StringExpr) GetLine() int   { return expr.Line }
func (expr StringExpr) GetColumn() int { return expr.Column }

type IdentifierExpr struct {
	Name   string
	Line   int
	Column int
}

func (expr IdentifierExpr) expr()          {}
func (expr IdentifierExpr) GetLine() int   { return expr.Line }
func (expr IdentifierExpr) GetColumn() int { return expr.Column }

type IntLiteralExpr struct {
	Value  int64
	Line   int
	Column int
}

func (expr IntLiteralExpr) expr()          {}
func (expr IntLiteralExpr) GetLine() int   { return expr.Line }
func (expr IntLiteralExpr) GetColumn() int { return expr.Column }

type BoolLiteralExpr struct {
	Value  bool
	Line   int
	Column int
}

func (expr BoolLiteralExpr) expr()          {}
func (expr BoolLiteralExpr) GetLine() int   { return expr.Line }
func (expr BoolLiteralExpr) GetColumn() int { return expr.Column }

type CharLiteralExpr struct {
	Value  rune
	Line   int
	Column int
}

func (expr CharLiteralExpr) expr()          {}
func (expr CharLiteralExpr) GetLine() int   { return expr.Line }
func (expr CharLiteralExpr) GetColumn() int { return expr.Column }

type ThisExpr struct {
	Line   int
	Column int
}

func (expr ThisExpr) expr()          {}
func (expr ThisExpr) GetLine() int   { return expr.Line }
func (expr ThisExpr) GetColumn() int { return expr.Column }

type NullLiteralExpr struct {
	Line   int
	Column int
}

func (expr NullLiteralExpr) expr()          {}
func (expr NullLiteralExpr) GetLine() int   { return expr.Line }
func (expr NullLiteralExpr) GetColumn() int { return expr.Column }

type BinaryExpr struct {
	Left     Expr
	Operator lexer.Token
	Right    Expr
	Line     int
	Column   int
}

func (expr BinaryExpr) expr()          {}
func (expr BinaryExpr) GetLine() int   { return expr.Line }
func (expr BinaryExpr) GetColumn() int { return expr.Column }

type PrefixExpr struct {
	Operator   lexer.Token
	Expression Expr
	Line       int
	Column     int
}

func (expr PrefixExpr) expr()          {}
func (expr PrefixExpr) GetLine() int   { return expr.Line }
func (expr PrefixExpr) GetColumn() int { return expr.Column }

type AssignmentExpr struct {
	Assignee Expr
	Operator lexer.Token
	Value    Expr
	Line     int
	Column   int
}

func (expr AssignmentExpr) expr()          {}
func (expr AssignmentExpr) GetLine() int   { return expr.Line }
func (expr AssignmentExpr) GetColumn() int { return expr.Column }

type MethodCallExpr struct {
	Receiver   Expr
	MethodName string
	Args       []Expr
	Line       int
	Column     int
}

func (expr MethodCallExpr) expr()          {}
func (expr MethodCallExpr) GetLine() int   { return expr.Line }
func (expr MethodCallExpr) GetColumn() int { return expr.Column }

type MemberAccessExpr struct {
	Receiver Expr
	Member   string
	Line     int
	Column   int
}

func (expr MemberAccessExpr) expr()          {}
func (expr MemberAccessExpr) GetLine() int   { return expr.Line }
func (expr MemberAccessExpr) GetColumn() int { return expr.Column }

type ConstructorCallExpr struct {
	TypeName string
	Args     []Expr
	Line     int
	Column   int
}

func (expr ConstructorCallExpr) expr()          {}
func (expr ConstructorCallExpr) GetLine() int   { return expr.Line }
func (expr ConstructorCallExpr) GetColumn() int { return expr.Column }

type PreIncrementExpr struct {
	Operand Expr
	Line    int
	Column  int
}

func (expr PreIncrementExpr) expr()          {}
func (expr PreIncrementExpr) GetLine() int   { return expr.Line }
func (expr PreIncrementExpr) GetColumn() int { return expr.Column }

type PostIncrementExpr struct {
	Operand Expr
	Line    int
	Column  int
}

func (expr PostIncrementExpr) expr()          {}
func (expr PostIncrementExpr) GetLine() int   { return expr.Line }
func (expr PostIncrementExpr) GetColumn() int { return expr.Column }

type PreDecrementExpr struct {
	Operand Expr
	Line    int
	Column  int
}

func (expr PreDecrementExpr) expr()          {}
func (expr PreDecrementExpr) GetLine() int   { return expr.Line }
func (expr PreDecrementExpr) GetColumn() int { return expr.Column }

type PostDecrementExpr struct {
	Operand Expr
	Line    int
	Column  int
}

func (expr PostDecrementExpr) expr()          {}
func (expr PostDecrementExpr) GetLine() int   { return expr.Line }
func (expr PostDecrementExpr) GetColumn() int { return expr.Column }
