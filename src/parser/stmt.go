package parser

import (
	"fmt"

	"github.com/FabianRolfMatthiasNoll/Golang-CSharp-Compiler/src/ast"
	"github.com/FabianRolfMatthiasNoll/Golang-CSharp-Compiler/src/lexer"
)

func parseStatement(p *parser) ast.Stmt {
	line, column := p.currentToken().Line, p.currentToken().Column
	stmt_fn, exists := stmtTable[p.currentTokenKind()]
	if exists {
		return stmt_fn(p)
	}

	if p.currentTokenKind() == lexer.RETURN {
		return parseReturnStmt(p)
	}

	if isType(p) {
		return parseVarDeclStmt(p)
	}

	expression := parseExpression(p, DEFAULT)
	if idExpr, ok := expression.(ast.IdentifierExpr); ok && p.currentTokenKind() == lexer.OPEN_PAREN {
		expression = parseMethodCallExpr(p, ast.ThisExpr{Line: idExpr.Line, Column: idExpr.Column}, idExpr.Name)
	}
	p.expect(lexer.SEMICOLON)

	return ast.ExpressionStmt{
		Expression: expression,
		Line:       line,
		Column:     column,
	}
}

func parseReturnStmt(p *parser) ast.Stmt {
	line, column := p.currentToken().Line, p.currentToken().Column
	p.advance() // consume 'return'

	var expression ast.Expr
	if p.currentTokenKind() != lexer.SEMICOLON {
		expression = parseExpression(p, DEFAULT)
	}

	p.expect(lexer.SEMICOLON)

	return ast.ReturnStmt{
		Value:  expression,
		Line:   line,
		Column: column,
	}
}

func parseVarDeclStmt(p *parser) ast.Stmt {
	line, column := p.currentToken().Line, p.currentToken().Column
	modifiers := parseModifiers(p)

	// Check if the current token is a type
	if !isType(p) {
		panic(fmt.Sprintf("Expected type but got %s at line %d, column %d", lexer.TokenKindString(p.currentTokenKind()), p.currentToken().Line, p.currentToken().Column))
	}
	dataType := parseType(p)

	identifier := p.expectError(lexer.IDENTIFIER, "Expected identifier after type declaration").Value

	var assignedValue ast.Expr
	if p.currentTokenKind() == lexer.ASSIGNMENT {
		p.advance() // consume '='
		assignedValue = parseExpression(p, ASSIGNMENT)
	} else {
		assignedValue = assignStandardType(dataType, assignedValue, p)
	}

	p.expect(lexer.SEMICOLON)

	return ast.VarDeclStmt{
		Modifiers:  modifiers,
		Identifier: identifier,
		Type:       dataType,
		Value:      assignedValue,
		Line:       line,
		Column:     column,
	}
}

func parseClassDeclStmt(p *parser) ast.Stmt {
	line, column := p.currentToken().Line, p.currentToken().Column
	modifiers := parseModifiers(p)
	p.expect(lexer.CLASS)
	className := p.expectError(lexer.IDENTIFIER, "Expected class name").Value
	p.expect(lexer.OPEN_BRACE)
	members := []ast.ClassMember{}

	for p.currentTokenKind() != lexer.CLOSE_BRACE {
		members = append(members, parseClassMember(p, className))
	}

	p.expect(lexer.CLOSE_BRACE)

	// Check if a constructor declaration exists
	hasConstructor := false
	for _, member := range members {
		if _, ok := member.(ast.ConstructorDeclStmt); ok {
			hasConstructor = true
			break
		}
	}

	// Add a standard constructor if no constructor declaration exists
	if !hasConstructor {
		standardConstructor := ast.ConstructorDeclStmt{
			Modifiers:  modifiers,
			Name:       className,
			Parameters: []ast.Parameter{},
			Body:       ast.BlockStmt{Body: []ast.Stmt{}},
			Line:       0,
			Column:     0,
		}
		members = append(members, standardConstructor)
	}

	return ast.ClassDeclStmt{
		Modifiers: modifiers,
		Name:      className,
		Body:      ast.ClassBody{Members: members},
		Line:      line,
		Column:    column,
	}
}

func parseClassMember(p *parser, className string) ast.ClassMember {
	line, column := p.currentToken().Line, p.currentToken().Column
	modifiers := parseModifiers(p)

	if isType(p) {
		return parseFieldOrMethod(p, modifiers)
	} else if p.currentTokenKind() == lexer.IDENTIFIER && p.currentToken().Value == className {
		// Possible constructor
		return parseConstructor(p, modifiers)
	}

	panic(fmt.Sprintf("Expected type or constructor but got %s at line %d, column %d", lexer.TokenKindString(p.currentTokenKind()), line, column))
}

func parseFieldOrMethod(p *parser, modifiers []ast.Modifier) ast.ClassMember {
	line, column := p.currentToken().Line, p.currentToken().Column
	dataType := parseType(p)
	identifier := p.expectError(lexer.IDENTIFIER, "Expected identifier").Value

	if p.currentTokenKind() == lexer.OPEN_PAREN {
		// It's a method
		return parseMethod(p, modifiers, dataType, identifier)
	} else {
		// It's a field
		var assignedValue ast.Expr
		if p.currentTokenKind() == lexer.ASSIGNMENT {
			p.advance() // consume '='
			assignedValue = parseExpression(p, ASSIGNMENT)
		} else {
			assignedValue = assignStandardType(dataType, assignedValue, p)
		}

		p.expect(lexer.SEMICOLON)
		return ast.FieldDeclStmt{
			Modifiers:  modifiers,
			Type:       dataType,
			Identifier: identifier,
			Value:      assignedValue,
			Line:       line,
			Column:     column,
		}
	}
}

func parseConstructor(p *parser, modifiers []ast.Modifier) ast.ClassMember {
	line, column := p.currentToken().Line, p.currentToken().Column
	name := p.expectError(lexer.IDENTIFIER, "Expected constructor name").Value
	p.expect(lexer.OPEN_PAREN)
	parameters := parseParameters(p)
	p.expect(lexer.CLOSE_PAREN)
	body := parseBlockStmt(p)

	return ast.ConstructorDeclStmt{
		Modifiers:  modifiers,
		Name:       name,
		Parameters: parameters,
		Body:       body,
		Line:       line,
		Column:     column,
	}
}

func parseMethod(p *parser, modifiers []ast.Modifier, returnType ast.Type, name string) ast.ClassMember {
	line, column := p.currentToken().Line, p.currentToken().Column
	p.expect(lexer.OPEN_PAREN)
	parameters := parseParameters(p)
	p.expect(lexer.CLOSE_PAREN)
	body := parseBlockStmt(p)

	return ast.MethodDeclStmt{
		Modifiers:  modifiers,
		ReturnType: returnType,
		Name:       name,
		Parameters: parameters,
		Body:       body,
		Line:       line,
		Column:     column,
	}
}

func parseParameters(p *parser) []ast.Parameter {
	parameters := []ast.Parameter{}

	for p.currentTokenKind() != lexer.CLOSE_PAREN {
		paramType := parseType(p)
		paramIdentifier := p.expectError(lexer.IDENTIFIER, "Expected parameter name").Value
		parameters = append(parameters, ast.Parameter{
			Type:       paramType,
			Identifier: paramIdentifier,
		})

		if p.currentTokenKind() == lexer.COMMA {
			p.advance()
		}
	}

	return parameters
}

func parseBlockStmt(p *parser) ast.BlockStmt {
	line, column := p.currentToken().Line, p.currentToken().Column
	p.expect(lexer.OPEN_BRACE)
	body := []ast.Stmt{}

	for p.currentTokenKind() != lexer.CLOSE_BRACE {
		body = append(body, parseStatement(p))
	}

	p.expect(lexer.CLOSE_BRACE)
	return ast.BlockStmt{
		Body:   body,
		Line:   line,
		Column: column,
	}
}

func parseWhileStmt(p *parser) ast.Stmt {
	line, column := p.currentToken().Line, p.currentToken().Column
	p.advance()
	p.expect(lexer.OPEN_PAREN)
	condition := parseExpression(p, DEFAULT)
	p.expect(lexer.CLOSE_PAREN)
	body := parseBlockStmt(p)
	return ast.WhileStmt{Condition: condition, Body: body, Line: line, Column: column}
}

func parseForStmt(p *parser) ast.Stmt {
	line, column := p.currentToken().Line, p.currentToken().Column
	p.advance()
	p.expect(lexer.OPEN_PAREN)

	var initializer ast.Stmt
	if p.currentTokenKind() != lexer.SEMICOLON {
		initializer = parseStatement(p)
	} else {
		p.advance()
	}

	var condition ast.Expr
	if p.currentTokenKind() != lexer.SEMICOLON {
		condition = parseExpression(p, DEFAULT)
	}
	p.expect(lexer.SEMICOLON)

	var increment ast.Expr
	if p.currentTokenKind() != lexer.CLOSE_PAREN {
		increment = parseExpression(p, DEFAULT)
	}
	p.expect(lexer.CLOSE_PAREN)

	body := parseBlockStmt(p)

	// convert for loop to while loop
	body.Body = append(body.Body, ast.ExpressionStmt{Expression: increment, Line: increment.GetLine(), Column: increment.GetColumn()})
	while := ast.WhileStmt{Condition: condition, Body: body, Line: line, Column: column}

	return ast.BlockStmt{Body: []ast.Stmt{initializer, while}, Line: line, Column: column}
}

func parseContinueStmt(p *parser) ast.Stmt {
	line, column := p.currentToken().Line, p.currentToken().Column
	p.advance()
	p.expect(lexer.SEMICOLON)
	return ast.ContinueStmt{Line: line, Column: column}
}

func parseBreakStmt(p *parser) ast.Stmt {
	line, column := p.currentToken().Line, p.currentToken().Column
	p.advance()
	p.expect(lexer.SEMICOLON)
	return ast.BreakStmt{Line: line, Column: column}
}

func parseIfStmt(p *parser) ast.Stmt {
	line, column := p.currentToken().Line, p.currentToken().Column
	p.advance()
	p.expect(lexer.OPEN_PAREN)
	condition := parseExpression(p, DEFAULT)
	p.expect(lexer.CLOSE_PAREN)
	then := parseBlockStmt(p)

	var elseStmt ast.Stmt
	if p.currentTokenKind() == lexer.ELSE {
		p.advance()
		elseStmt = parseBlockStmt(p)
	}

	return ast.IfStmt{Condition: condition, Then: then, Else: elseStmt, Line: line, Column: column}
}

func parseSwitchStmt(p *parser) ast.Stmt {
	line, column := p.currentToken().Line, p.currentToken().Column
	p.advance()
	p.expect(lexer.OPEN_PAREN)
	expression := parseExpression(p, DEFAULT)
	p.expect(lexer.CLOSE_PAREN)
	p.expect(lexer.OPEN_BRACE)

	cases := []ast.SwitchCase{}
	var defaultCase ast.Stmt
	for p.currentTokenKind() != lexer.CLOSE_BRACE {
		if p.currentTokenKind() == lexer.CASE {
			cases = append(cases, parseSwitchCase(p))
		} else if p.currentTokenKind() == lexer.DEFAULT {
			defaultCase = parseDefaultCase(p)
		} else {
			panic(fmt.Sprintf("Expected case or default but got %s at line %d, column %d", lexer.TokenKindString(p.currentTokenKind()), p.currentToken().Line, p.currentToken().Column))
		}
	}

	p.expect(lexer.CLOSE_BRACE)

	return ast.SwitchStmt{Expression: expression, Cases: cases, Default: defaultCase, Line: line, Column: column}
}

func parseSwitchCase(p *parser) ast.SwitchCase {
	line, column := p.currentToken().Line, p.currentToken().Column
	p.advance()
	value := parseExpression(p, DEFAULT)
	p.expect(lexer.COLON)
	body := []ast.Stmt{}

	for p.currentTokenKind() != lexer.CASE && p.currentTokenKind() != lexer.DEFAULT && p.currentTokenKind() != lexer.CLOSE_BRACE {
		body = append(body, parseStatement(p))
	}

	bodyBlock := ast.BlockStmt{Body: body, Line: line, Column: column}

	return ast.SwitchCase{Value: value, Body: bodyBlock, Line: line, Column: column}
}

func parseDefaultCase(p *parser) ast.Stmt {
	line, column := p.currentToken().Line, p.currentToken().Column
	p.advance()
	p.expect(lexer.COLON)
	body := []ast.Stmt{}

	for p.currentTokenKind() != lexer.CLOSE_BRACE {
		body = append(body, parseStatement(p))
	}

	bodyBlock := ast.BlockStmt{Body: body, Line: line, Column: column}

	return bodyBlock
}
