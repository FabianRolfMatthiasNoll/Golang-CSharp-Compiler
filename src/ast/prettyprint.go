package ast

import (
	"fmt"
	"strings"

	"github.com/FabianRolfMatthiasNoll/Golang-CSharp-Compiler/src/lexer"
)

func indentString(s string, level int) string {
	indent := strings.Repeat("  ", level)
	return indent + strings.Replace(s, "\n", "\n"+indent, -1)
}

func (prog Program) String() string {
	classes := make([]string, len(prog.Classes))
	for i, c := range prog.Classes {
		classes[i] = indentString(c.String(), 1)
	}
	return fmt.Sprintf("Program{\n  Classes: [\n%s\n  ]\n}", strings.Join(classes, ",\n"))
}

//=========================================================================================================
// Expressions
//=========================================================================================================

func (expr TypedExpr) String() string {
	return fmt.Sprintf("TypedExpr{\n  Type: %s,\n  Expression: %s\n}", expr.Type, indentString(fmt.Sprintf("%s", expr.Expr), 1))
}

func (expr IntLiteralExpr) String() string {
	return fmt.Sprintf("IntLiteralExpr{\n  Value: %d\n}", expr.Value)
}

func (expr StringExpr) String() string {
	return fmt.Sprintf("StringExpr{\n  Value: %s\n}", expr.Value)
}

func (expr IdentifierExpr) String() string {
	return fmt.Sprintf("IdentifierExpr{\n  Name: %s\n}", expr.Name)
}

func (expr FieldVarExpr) String() string {
	return fmt.Sprintf("FieldVarExpr{\n  Name: %s\n}", expr.Name)
}

func (expr LocalVarExpr) String() string {
	return fmt.Sprintf("LocalVarExpr{\n  Name: %s\n}", expr.Name)
}

func (expr ThisExpr) String() string {
	return "ThisExpr{}"
}

func (expr BoolLiteralExpr) String() string {
	return fmt.Sprintf("BoolLiteralExpr{\n  Value: %t\n}", expr.Value)
}

func (expr NullLiteralExpr) String() string {
	return "NullLiteralExpr{}"
}

func (expr CharLiteralExpr) String() string {
	return fmt.Sprintf("CharExpr{\n  Value: %c\n}", expr.Value)
}

func (expr BinaryExpr) String() string {
	return fmt.Sprintf("BinaryExpr{\n  Left: %s,\n  Operator: %s,\n  Right: %s\n}",
		indentString(fmt.Sprintf("%s", expr.Left), 1), expr.Operator, indentString(fmt.Sprintf("%s", expr.Right), 1))
}

func (expr PrefixExpr) String() string {
	return fmt.Sprintf("PrefixExpr{\n  Operator: %s,\n  Expression: %s\n}", expr.Operator, indentString(fmt.Sprintf("%s", expr.Expression), 1))
}

func (expr AssignmentExpr) String() string {
	return fmt.Sprintf("AssignmentExpr{\n  Assignee: %s,\n  Operator: %s,\n  Value: %s\n}",
		indentString(fmt.Sprintf("%s", expr.Assignee), 1), expr.Operator, indentString(fmt.Sprintf("%s", expr.Value), 1))
}

func (expr MethodCallExpr) String() string {
	args := make([]string, len(expr.Args))
	for i, arg := range expr.Args {
		args[i] = indentString(fmt.Sprintf("%s", arg), 2)
	}
	return fmt.Sprintf("MethodCallExpr{\n  Receiver: %s,\n  MethodName: %s,\n  Arguments: [\n%s\n  ]\n}", indentString(fmt.Sprintf("%s", expr.Receiver), 1), expr.MethodName, strings.Join(args, ",\n"))
}

func (expr MemberAccessExpr) String() string {
	return fmt.Sprintf("MemberAccessExpr{\n  Receiver: %s,\n  Member: %s\n}", indentString(fmt.Sprintf("%s", expr.Receiver), 1), expr.Member)
}

func (expr ConstructorCallExpr) String() string {
	args := make([]string, len(expr.Args))
	for i, arg := range expr.Args {
		args[i] = indentString(fmt.Sprintf("%s", arg), 2)
	}
	return fmt.Sprintf("ConstructorCallExpr{\n  ClassName: %s,\n  Arguments: [\n%s\n  ]\n}", expr.TypeName, strings.Join(args, ",\n"))
}

func (expr PreDecrementExpr) String() string {
	return fmt.Sprintf("PreDecrementExpr{\n  Operand: %s\n}", indentString(fmt.Sprintf("%s", expr.Operand), 1))
}

func (expr PostDecrementExpr) String() string {
	return fmt.Sprintf("PostDecrementExpr{\n  Operand: %s\n}", indentString(fmt.Sprintf("%s", expr.Operand), 1))
}

func (expr PreIncrementExpr) String() string {
	return fmt.Sprintf("PreIncrementExpr{\n  Operand: %s\n}", indentString(fmt.Sprintf("%s", expr.Operand), 1))
}

func (expr PostIncrementExpr) String() string {
	return fmt.Sprintf("PostIncrementExpr{\n  Operand: %s\n}", indentString(fmt.Sprintf("%s", expr.Operand), 1))
}

//=========================================================================================================
// Statements
//=========================================================================================================

func (stmt TypedStmt) String() string {
	return fmt.Sprintf("TypedStmt{\n  Type: %s,\n  Statement: %s\n}", stmt.Type, indentString(fmt.Sprintf("%s", stmt.Stmt), 1))
}

func (stmt BlockStmt) String() string {
	body := make([]string, len(stmt.Body))
	for i, s := range stmt.Body {
		body[i] = indentString(fmt.Sprintf("%s", s), 1)
	}
	return fmt.Sprintf("BlockStmt{\n  [\n%s\n  ]\n}", strings.Join(body, ",\n"))
}

func (stmt ExpressionStmt) String() string {
	return fmt.Sprintf("ExpressionStmt{\n  Expression: %s\n}", indentString(fmt.Sprintf("%s", stmt.Expression), 1))
}

func (stmt VarDeclStmt) String() string {
	modifiers := make([]string, len(stmt.Modifiers))
	for i, mod := range stmt.Modifiers {
		modifiers[i] = lexer.TokenKindString(mod.Kind)
	}
	return fmt.Sprintf("VarDeclStmt{\n  Modifiers: [%s],\n  Type: %s,\n  Identifier: %s,\n  Value: %s\n}",
		strings.Join(modifiers, ", "), stmt.Type.Name, stmt.Identifier, indentString(fmt.Sprintf("%s", stmt.Value), 1))
}

func (stmt ClassDeclStmt) String() string {
	modifiers := make([]string, len(stmt.Modifiers))
	for i, mod := range stmt.Modifiers {
		modifiers[i] = lexer.TokenKindString(mod.Kind)
	}
	return fmt.Sprintf("ClassDeclStmt{\n  Modifiers: [%s],\n  Name: %s,\n  Body: %s\n}",
		strings.Join(modifiers, ", "), stmt.Name, indentString(stmt.Body.String(), 1))
}

func (body ClassBody) String() string {
	members := make([]string, len(body.Members))
	for i, m := range body.Members {
		members[i] = indentString(fmt.Sprintf("%s", m), 1)
	}
	return fmt.Sprintf("ClassBody{\n  Members: [\n%s\n  ]\n}", strings.Join(members, ",\n"))
}

func (stmt FieldDeclStmt) String() string {
	modifiers := make([]string, len(stmt.Modifiers))
	for i, mod := range stmt.Modifiers {
		modifiers[i] = lexer.TokenKindString(mod.Kind)
	}
	return fmt.Sprintf("FieldDeclStmt{\n  Modifiers: [%s],\n  Type: %s,\n  Identifier: %s,\n  Value: %s\n}",
		strings.Join(modifiers, ", "), stmt.Type.Name, stmt.Identifier, indentString(fmt.Sprintf("%s", stmt.Value), 1))
}

func (stmt MethodDeclStmt) String() string {
	modifiers := make([]string, len(stmt.Modifiers))
	for i, mod := range stmt.Modifiers {
		modifiers[i] = lexer.TokenKindString(mod.Kind)
	}
	params := make([]string, len(stmt.Parameters))
	for i, p := range stmt.Parameters {
		params[i] = fmt.Sprintf("%s %s", p.Type.Name, p.Identifier)
	}
	return fmt.Sprintf("MethodDeclStmt{\n  Modifiers: [%s],\n  ReturnType: %s,\n  Name: %s,\n  Parameters: [%s],\n  Body: %s\n}",
		strings.Join(modifiers, ", "), stmt.ReturnType.Name, stmt.Name, strings.Join(params, ", "), indentString(fmt.Sprintf("%s", stmt.Body), 1))
}

func (stmt ConstructorDeclStmt) String() string {
	modifiers := make([]string, len(stmt.Modifiers))
	for i, mod := range stmt.Modifiers {
		modifiers[i] = lexer.TokenKindString(mod.Kind)
	}
	params := make([]string, len(stmt.Parameters))
	for i, p := range stmt.Parameters {
		params[i] = fmt.Sprintf("%s %s", p.Type.Name, p.Identifier)
	}
	return fmt.Sprintf("ConstructorDeclStmt{\n  Modifiers: [%s],\n  Name: %s,\n  Parameters: [%s],\n  Body: %s\n}",
		strings.Join(modifiers, ", "), stmt.Name, strings.Join(params, ", "), indentString(fmt.Sprintf("%s", stmt.Body), 1))
}

func (stmt ReturnStmt) String() string {
	return fmt.Sprintf("ReturnStmt{\n  Value: %s\n}", indentString(fmt.Sprintf("%s", stmt.Value), 1))
}

func (stmt IfStmt) String() string {
	return fmt.Sprintf("IfStmt{\n  Condition: %s,\n  Then: %s,\n  Else: %s\n}",
		indentString(fmt.Sprintf("%s", stmt.Condition), 1), indentString(fmt.Sprintf("%s", stmt.Then), 1), indentString(fmt.Sprintf("%s", stmt.Else), 1))
}

func (stmt WhileStmt) String() string {
	return fmt.Sprintf("WhileStmt{\n  Condition: %s,\n  Body: %s\n}",
		indentString(fmt.Sprintf("%s", stmt.Condition), 1), indentString(fmt.Sprintf("%s", stmt.Body), 1))
}

func (stmt ContinueStmt) String() string {
	return "ContinueStmt{}"
}

func (stmt BreakStmt) String() string {
	return "BreakStmt{}"
}

func (stmt SwitchStmt) String() string {
	cases := make([]string, len(stmt.Cases))
	for i, c := range stmt.Cases {
		cases[i] = indentString(fmt.Sprintf("%s", c), 1)
	}
	defaultCase := ""
	if stmt.Default != nil {
		defaultCase = indentString(fmt.Sprintf("%s", stmt.Default), 1)
	}
	return fmt.Sprintf("SwitchStmt{\n  Expression: %s,\n  Cases: [\n%s\n  ],\n  DefaultCase: %s\n}",
		indentString(fmt.Sprintf("%s", stmt.Expression), 1), strings.Join(cases, ",\n"), defaultCase)
}

func (stmt SwitchCase) String() string {
	return fmt.Sprintf("SwitchCase{\n  Value: %s,\n  Body: %s\n}",
		indentString(fmt.Sprintf("%s", stmt.Value), 1), indentString(fmt.Sprintf("%s", stmt.Body), 1))
}
