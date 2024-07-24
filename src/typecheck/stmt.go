package typecheck

import (
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

	// The entry for this has to exist at this point
	symbolEntry, _ := tc.env.Lookup("this")

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
		case ast.ExpressionStmt:
			block.Body[i] = tc.CheckExpressionStmt(&stmt)
		case ast.BlockStmt:
			block.Body[i] = tc.CheckBlockStmt(&stmt)
			possibleBlockTypes = append(possibleBlockTypes, block.Body[i].(ast.TypedStmt).Type)
		case ast.IfStmt:
			block.Body[i] = tc.CheckIfStmt(&stmt)
			possibleBlockTypes = append(possibleBlockTypes, block.Body[i].(ast.TypedStmt).Type)
		case ast.WhileStmt:
			block.Body[i] = tc.CheckWhileStmt(&stmt)
			possibleBlockTypes = append(possibleBlockTypes, block.Body[i].(ast.TypedStmt).Type)
		case ast.ReturnStmt:
			block.Body[i] = tc.CheckReturnStmt(&stmt)
			possibleBlockTypes = append(possibleBlockTypes, block.Body[i].(ast.TypedStmt).Type)
		case ast.BreakStmt:
			block.Body[i] = ast.TypedStmt{Stmt: stmt, Type: "void"}
		case ast.ContinueStmt:
			block.Body[i] = ast.TypedStmt{Stmt: stmt, Type: "void"}
		default:
			tc.errorf(stmt.GetLine(), stmt.GetColumn(), "unexpected statement")
		}
	}
	blockType := tc.upperBound(possibleBlockTypes)
	return ast.TypedStmt{Stmt: block, Type: blockType}
}

func (tc *TypeChecker) CheckExpressionStmt(expr *ast.ExpressionStmt) ast.TypedStmt {
	expr.Expression = tc.CheckExpr(expr.Expression)
	return ast.TypedStmt{Stmt: expr, Type: expr.Expression.(ast.TypedExpr).Type}
}

func (tc *TypeChecker) CheckReturnStmt(stmt *ast.ReturnStmt) ast.TypedStmt {

	typ := "void"
	if stmt.Value != nil {
		stmt.Value = tc.CheckExpr(stmt.Value)
		typ = stmt.Value.(ast.TypedExpr).Type
	}

	// thisMethod must exist at this point
	method, _ := tc.env.Lookup("thisMethod")

	if !tc.isTypeCompatible(method.Type, typ) {
		tc.errorf(stmt.GetLine(), stmt.GetColumn(), "type mismatch: expected %s, got %s", method.Type, typ)
	}

	return ast.TypedStmt{Stmt: stmt, Type: typ}
}

func (tc *TypeChecker) CheckWhileStmt(stmt *ast.WhileStmt) ast.TypedStmt {
	stmt.Condition = tc.checkBoolCondition(stmt.Condition)

	if block, ok := stmt.Body.(ast.BlockStmt); ok {
		stmt.Body = tc.CheckBlockStmt(&block)
	} else {
		tc.errorf(stmt.GetLine(), stmt.GetColumn(), "while body should be a block statement")
	}

	return ast.TypedStmt{Stmt: stmt, Type: stmt.Body.(ast.TypedStmt).Type}
}

func (tc *TypeChecker) CheckIfStmt(stmt *ast.IfStmt) ast.TypedStmt {
	stmt.Condition = tc.checkBoolCondition(stmt.Condition)
	var thenType, elseType string

	if thenBlock, ok := stmt.Then.(ast.BlockStmt); ok {
		stmt.Then = tc.CheckBlockStmt(&thenBlock)
		thenType = stmt.Then.(ast.TypedStmt).Type
	} else {
		tc.errorf(stmt.GetLine(), stmt.GetColumn(), "while body should be a block statement")
	}

	if elseBlock, ok := stmt.Else.(ast.BlockStmt); ok {
		stmt.Else = tc.CheckBlockStmt(&elseBlock)
		elseType = stmt.Else.(ast.TypedStmt).Type
	} else {
		tc.errorf(stmt.GetLine(), stmt.GetColumn(), "while body should be a block statement")
	}

	ifType := tc.upperBound([]string{thenType, elseType})

	if thenType == "void" || elseType == "void" {
		ifType = "void"
	}

	return ast.TypedStmt{Stmt: stmt, Type: ifType, Line: stmt.Line, Column: stmt.Column}
}
