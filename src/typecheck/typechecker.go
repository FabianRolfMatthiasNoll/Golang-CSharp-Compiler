package typecheck

import (
	"fmt"

	"github.com/FabianRolfMatthiasNoll/Golang-CSharp-Compiler/src/ast"
)

type TypeChecker struct {
	env     *TypeEnvironment
	classes map[string]ast.ClassDeclStmt
}

func NewTypeChecker() *TypeChecker {
	return &TypeChecker{env: NewTypeEnv(nil)}
}

func (tc *TypeChecker) CheckProgram(prog *ast.Program) ast.Program {
	tc.classes = make(map[string]ast.ClassDeclStmt)
	for i := range prog.Classes {
		tc.classes[prog.Classes[i].Name] = prog.Classes[i]
	}

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
		case ast.MethodDeclStmt:
			tc.CheckMethodDeclStmt(&member)
			m = member
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
		return tc.CheckBinaryExpr(e)
	case ast.MethodCallExpr:
		return tc.CheckMethodCallExpr(e)
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

func (tc *TypeChecker) CheckBinaryExpr(expr ast.Expr) ast.TypedExpr {
	// TODO: Implement
	return ast.TypedExpr{}
}

func (tc *TypeChecker) CheckMethodCallExpr(expr ast.MethodCallExpr) ast.TypedExpr {
	// TODO: Implement
	return ast.TypedExpr{}
}

func (tc *TypeChecker) CheckMethodDeclStmt(method *ast.MethodDeclStmt) {
	tc.env = NewTypeEnv(tc.env) // New scope for method
	defer func() { tc.env = tc.env.outer }()

	// Register parameters
	for _, param := range method.Parameters {
		tc.env.Define(param.Identifier, param.Type.Name, false, false, true)
	}

	tc.env.Define("thisMethod", method.ReturnType.Name, false, false, false)

	// Check and type method body
	if block, ok := method.Body.(ast.BlockStmt); ok {
		method.Body = tc.CheckBlockStmt(&block)
	} else {
		tc.errorf(method.GetLine(), method.GetColumn(), "method body should be a block statement")
	}

	// Check return type
	if !tc.isTypeCompatible(method.ReturnType.Name, method.Body.(ast.TypedStmt).Type) {
		tc.errorf(method.GetLine(), method.GetColumn(), "type mismatch: expected %s, got %s", method.ReturnType.Name, method.Body.(ast.TypedStmt).Type)
	}
}

func (tc *TypeChecker) CheckBlockStmt(block *ast.BlockStmt) ast.TypedStmt {
	tc.env = NewTypeEnv(tc.env)
	defer func() { tc.env = tc.env.outer }()

	for i, stmt := range block.Body {
		switch stmt := stmt.(type) {
		// case ast.VarDeclStmt:
		// 	tc.CheckVarDeclStmt(&stmt)
		// case ast.AssignStmt:
		// 	tc.CheckAssignStmt(&stmt)
		// case ast.ReturnStmt:
		// 	tc.CheckReturnStmt(&stmt)
		// case ast.IfStmt:
		// 	tc.CheckIfStmt(&stmt)
		// case ast.WhileStmt:
		// 	tc.CheckWhileStmt(&stmt)
		// case ast.MethodCallStmt:
		// 	tc.CheckMethodCallStmt(&stmt)
		default:
			tc.errorf(stmt.GetLine(), stmt.GetColumn(), "unexpected statement")
		}
		block.Body[i] = stmt
	}
	// TODO: Check type of block by upper bound of returns etc. Also check if all paths return
	return ast.TypedStmt{Stmt: block, Type: "void"}
}

func (tc *TypeChecker) isTypeCompatible(a, b string) bool {
	if a == "char" && b == "int" {
		return true
	} else if a == "int" && b == "char" {
		return true
	} else if tc.isUserObject(a) && b == "null" {
		return true
	} else if a == "null" && tc.isUserObject(b) {
		return true
	} else if a == b {
		return true
	}
	return false
}

func (tc *TypeChecker) isUserObject(typ string) bool {
	_, ok := tc.classes[typ]
	return ok
}

func (tc *TypeChecker) errorf(line, column int, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	panic(fmt.Sprintf("Error at line %d, column %d: %s", line, column, message))
}
