package ast

type BlockStmt struct {
	Body []Stmt
}

func (stmt BlockStmt) stmt() {}

// To enforce semicolons really easily
type ExpressionStmt struct {
	Expression Expr
}

func (stmt ExpressionStmt) stmt() {}
