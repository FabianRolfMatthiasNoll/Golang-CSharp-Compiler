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

//=========================================================================================================
// Expressions
//=========================================================================================================

func (expr NumberExpr) String() string {
	return fmt.Sprintf("NumberExpr{Value: %f}", expr.Value)
}

func (expr StringExpr) String() string {
	return fmt.Sprintf("StringExpr{Value: %s}", expr.Value)
}

func (expr IdentifierExpr) String() string {
	return fmt.Sprintf("IdentifierExpr{Name: %s}", expr.Name)
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
		indentString(fmt.Sprintf("%s", expr.Assigne), 1), expr.Operator, indentString(fmt.Sprintf("%s", expr.Value), 1))
}

//=========================================================================================================
// Statements
//=========================================================================================================

func (stmt BlockStmt) String() string {
	body := make([]string, len(stmt.Body))
	for i, s := range stmt.Body {
		body[i] = indentString(fmt.Sprintf("%s", s), 1)
	}
	return fmt.Sprintf("BlockStmt{\n  Body: [\n%s\n  ]\n}", strings.Join(body, ",\n"))
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
		strings.Join(modifiers, ", "), stmt.Name, indentString(fmt.Sprintf("%s", stmt.Body), 1))
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
	return fmt.Sprintf("ReturnStmt{\n  Value: %s,\n}", indentString(fmt.Sprintf("%s", stmt.Value), 1))
}