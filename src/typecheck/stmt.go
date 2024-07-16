package typecheck

import (
	"fmt"

	"github.com/FabianRolfMatthiasNoll/Golang-CSharp-Compiler/src/ast"
)

func (tc *TypeChecker) CheckClassDeclStmt(class *ast.ClassDeclStmt) {
	tc.env = NewTypeEnv(tc.env)              // Create new scope
	defer func() { tc.env = tc.env.outer }() // Pop scope after checking class

	// Register class fields
	for _, member := range class.Body.Members {
		if field, ok := member.(ast.FieldDeclStmt); ok {
			tc.env.Define(field.Identifier, field.Type.Name, true, true, false)
		}
	}

	tc.env.Define("this", class.Name, true, false, false)

	// Check members
	for i, member := range class.Body.Members {
		var updatedMember ast.ClassMember

		switch member := member.(type) {
		case ast.FieldDeclStmt:
			tc.CheckFieldDeclStmt(&member)
			updatedMember = member
		case ast.MethodDeclStmt:
			tc.CheckMethodDeclStmt(&member)
			updatedMember = member
		case ast.ConstructorDeclStmt:
			tc.CheckConstructorDeclStmt(&member)
			updatedMember = member
		default:
			//tc.errorf(member.GetLine(), member.GetColumn(), "unexpected class member")
			updatedMember = member
		}

		class.Body.Members[i] = updatedMember
	}
}

func (tc *TypeChecker) CheckFieldDeclStmt(field *ast.FieldDeclStmt) {
	typedExpression := tc.CheckExpr(field.Value)

	if !tc.isTypeCompatible(field.Type.Name, typedExpression.Type) {
		tc.errorf(field.Line, field.Column, "type mismatch: expected %s, got %s", field.Type.Name, typedExpression.Type)
	}

	field.Value = typedExpression
}

func (tc *TypeChecker) CheckMethodDeclStmt(method *ast.MethodDeclStmt) {
	tc.env = NewTypeEnv(tc.env)
	defer func() { tc.env = tc.env.outer }()

	// The entry for this has to exist at this point
	symbolEntry, _ := tc.env.Lookup("this")

	if symbolEntry.Type == method.Name {
		tc.errorf(method.GetLine(), method.GetColumn(), "method name can't be the same as the class name")
	}

	for _, param := range method.Parameters {
		tc.env.Define(param.Identifier, param.Type.Name, false, false, true)
	}

	// Register thisMethod so that returns deep into the method body can be checked against the return type
	// TODO: Check later on if this is really needed.
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

func (tc *TypeChecker) CheckConstructorDeclStmt(constructor *ast.ConstructorDeclStmt) {
	tc.env = NewTypeEnv(tc.env)
	defer func() { tc.env = tc.env.outer }()
	fmt.Println(constructor.Name)
	// The entry for this has to exist at this point
	symbolEntry, _ := tc.env.Lookup("this")
	fmt.Println(symbolEntry.Type)
	if symbolEntry.Type != constructor.Name {
		tc.errorf(constructor.GetLine(), constructor.GetColumn(), "constructor name must be the same as the class name")
	}

	for _, param := range constructor.Parameters {
		tc.env.Define(param.Identifier, param.Type.Name, false, false, true)
	}

	// Check and type constructor body
	if block, ok := constructor.Body.(ast.BlockStmt); ok {
		constructor.Body = tc.CheckBlockStmt(&block)
	} else {
		tc.errorf(constructor.GetLine(), constructor.GetColumn(), "constructor body should be a block statement")
	}

	// Check return type
	if constructor.Body.(ast.TypedStmt).Type != "void" {
		tc.errorf(constructor.GetLine(), constructor.GetColumn(), "constructor can not have a return")
	}
}

func (tc *TypeChecker) CheckBlockStmt(block *ast.BlockStmt) ast.TypedStmt {
	tc.env = NewTypeEnv(tc.env)
	defer func() { tc.env = tc.env.outer }()

	possibleBlockTypes := []string{}

	for i, stmt := range block.Body {
		switch stmt := stmt.(type) {
		// TODO: Check what expression statements are and document it
		//case ast.ExpressionStmt:
		//block.Body[i] = tc.CheckExpressionStmt(&stmt)
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
	blockType := tc.upperBound(possibleBlockTypes)
	return ast.TypedStmt{Stmt: block, Type: blockType}
}
