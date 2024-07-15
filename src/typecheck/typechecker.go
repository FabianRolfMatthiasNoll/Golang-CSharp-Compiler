package typecheck

import (
	"fmt"

	"github.com/FabianRolfMatthiasNoll/Golang-CSharp-Compiler/src/ast"
)

type TypeChecker struct {
	env *TypeEnvironment
}

func NewTypeChecker() *TypeChecker {
	return &TypeChecker{env: NewTypeEnv(nil)}
}

func (tc *TypeChecker) CheckProgram(prog *ast.Program) ast.Program {
	for i := range prog.Classes {
		tc.CheckClassDeclStmt(&prog.Classes[i])
	}

	return *prog
}

func (tc *TypeChecker) CheckClassDeclStmt(class *ast.ClassDeclStmt) {
	tc.env = NewTypeEnv(tc.env) // New scope for class
	defer func() { tc.env = tc.env.outer }()

	// Register class fields
	for _, member := range class.Body.Members {
		if field, ok := member.(ast.FieldDeclStmt); ok {
			tc.env.Define(field.Identifier, field.Type.Name, true, true, false)
		}
	}

	// Check members
	for i, member := range class.Body.Members {
		var m ast.ClassMember

		switch member := member.(type) {
		case ast.FieldDeclStmt:
			tc.CheckFieldDeclStmt(&member)
			m = member
		// case ast.MethodDeclStmt:
		// 	tc.CheckMethodDeclStmt(&member)
		// 	m = member
		// case ast.ConstructorDeclStmt:
		// 	tc.CheckConstructorDeclStmt(&member)
		// 	m = member
		default:
			//tc.errorf(member.GetLine(), member.GetColumn(), "unexpected class member")
			m = member
		}

		class.Body.Members[i] = m // Update member
	}
}

func (tc *TypeChecker) CheckFieldDeclStmt(field *ast.FieldDeclStmt) {
	if field.Value != nil {
		typedField := tc.CheckExpr(field.Value)
		// TODO: Need Type comparison method to allow nulls for strings and object and allow char to int etc.
		if typedField.Type != field.Type.Name && typedField.Type != "null" {
			tc.errorf(field.Line, field.Column, "type mismatch: expected %s, got %s", field.Type.Name, typedField)
		}
		field.Value = typedField
	}
}

// TODO: Impelement rest of check expr but with some sort of structure to control this monster of code
func (tc *TypeChecker) CheckExpr(expr ast.Expr) ast.TypedExpr {
	switch e := expr.(type) {
	case ast.IntLiteralExpr:
		return ast.TypedExpr{Type: "int", Expr: e}
	case ast.BoolLiteralExpr:
		return ast.TypedExpr{Type: "bool", Expr: e}
	case ast.StringExpr:
		return ast.TypedExpr{Type: "string", Expr: e}
	case ast.IdentifierExpr:
		info, ok := tc.env.Lookup(e.Name)
		if !ok {
			tc.errorf(e.Line, e.Column, "undefined variable: %s", e.Name)
		}
		return ast.TypedExpr{Type: info.Type, Expr: e}
	case ast.NullLiteralExpr:
		return ast.TypedExpr{Type: "null", Expr: e}
	case ast.CharLiteralExpr:
		return ast.TypedExpr{Type: "char", Expr: e}
	case ast.BinaryExpr:
		// TODO: Implement
	case ast.MethodCallExpr:
		// TODO: Implement
		return ast.TypedExpr{}
	case ast.AssignmentExpr:
		assigneeType := tc.CheckExpr(e.Assignee)
		valueType := tc.CheckExpr(e.Value)
		if assigneeType != valueType {
			tc.errorf(e.Line, e.Column, "type mismatch: %s and %s", assigneeType, valueType)
		}
		return assigneeType
	default:
		tc.errorf(expr.GetLine(), expr.GetColumn(), "unexpected expression")
	}
	return ast.TypedExpr{}
}

func (tc *TypeChecker) errorf(line, column int, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	panic(fmt.Sprintf("Error at line %d, column %d: %s", line, column, message))
}
