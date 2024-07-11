package parser

import (
	"fmt"

	"github.com/FabianRolfMatthiasNoll/Golang-CSharp-Compiler/src/ast"
	"github.com/FabianRolfMatthiasNoll/Golang-CSharp-Compiler/src/lexer"
)

func parseStatement(p *parser) ast.Stmt {
	stmt_fn, exists := stmtTable[p.currentTokenKind()]
	if exists {
		return stmt_fn(p)
	}

	if p.currentTokenKind() == lexer.RETURN {
		return parseReturnStmt(p)
	}

	expression := parseExpression(p, DEFAULT)
	p.expect(lexer.SEMICOLON)

	return ast.ExpressionStmt{
		Expression: expression,
		Line:       p.currentToken().Line,
		Column:     p.currentToken().Column,
	}
}

func parseReturnStmt(p *parser) ast.Stmt {
	token := p.advance() // consume 'return'

	var expression ast.Expr
	if p.currentTokenKind() != lexer.SEMICOLON {
		expression = parseExpression(p, DEFAULT)
	}

	p.expect(lexer.SEMICOLON)

	return ast.ReturnStmt{
		Value:  expression,
		Line:   token.Line,
		Column: token.Column,
	}
}

func parseVarDeclStmt(p *parser) ast.Stmt {
	modifiers := parseModifiers(p)
	dataType := parseType(p)

	identifier := p.expectError(lexer.IDENTIFIER, "Expected identifier after type declaration").Value

	var assignedValue ast.Expr
	if p.currentTokenKind() == lexer.ASSIGNMENT {
		p.advance() // consume '='
		assignedValue = parseExpression(p, ASSIGNMENT)
	}

	p.expect(lexer.SEMICOLON)

	return ast.VarDeclStmt{
		Modifiers:  modifiers,
		Identifier: identifier,
		Type:       dataType,
		Value:      assignedValue,
		Line:       p.currentToken().Line,
		Column:     p.currentToken().Column,
	}
}

func parseModifiers(p *parser) []ast.Modifier {
	modifiers := []ast.Modifier{}

	for isModifier(p.currentTokenKind()) {
		modifiers = append(modifiers, ast.Modifier{Kind: p.advance().Kind})
	}

	if len(modifiers) == 0 {
		modifiers = append(modifiers, ast.Modifier{Kind: lexer.PUBLIC})
	}

	return modifiers
}

func parseType(p *parser) ast.Type {
	return ast.Type{Name: p.advance().Value}
}

func isModifier(kind lexer.TokenKind) bool {
	switch kind {
	case lexer.PUBLIC, lexer.PRIVATE, lexer.PROTECTED, lexer.STATIC, lexer.FINAL:
		return true
	}
	return false
}

func isType(kind lexer.TokenKind) bool {
	switch kind {
	case lexer.INT, lexer.STRING_TYPE, lexer.FLOAT, lexer.BOOL, lexer.VAR, lexer.CHAR, lexer.DOUBLE, lexer.VOID:
		return true
	}
	return false
}

func parseClassDeclStmt(p *parser) ast.Stmt {
	modifiers := parseModifiers(p)
	p.expect(lexer.CLASS)
	className := p.expectError(lexer.IDENTIFIER, "Expected class name").Value
	p.expect(lexer.OPEN_BRACE)
	members := []ast.ClassMember{}

	for p.currentTokenKind() != lexer.CLOSE_BRACE {
		members = append(members, parseClassMember(p))
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
			Line:       p.currentToken().Line,
			Column:     p.currentToken().Column,
		}
		members = append(members, standardConstructor)
	}

	return ast.ClassDeclStmt{
		Modifiers: modifiers,
		Name:      className,
		Body:      ast.ClassBody{Members: members},
		Line:      p.currentToken().Line,
		Column:    p.currentToken().Column,
	}
}

func parseClassMember(p *parser) ast.ClassMember {
	modifiers := parseModifiers(p)

	if isType(p.currentTokenKind()) {
		typ := parseType(p)
		identifier := p.expectError(lexer.IDENTIFIER, "Expected identifier").Value

		if p.currentTokenKind() == lexer.OPEN_PAREN {
			// It's a method
			return parseMethod(p, modifiers, typ, identifier)
		} else {
			// It's a field
			var value ast.Expr
			if p.currentTokenKind() == lexer.ASSIGNMENT {
				p.advance() // consume '='
				value = parseExpression(p, ASSIGNMENT)
			}

			p.expect(lexer.SEMICOLON)
			return ast.FieldDeclStmt{
				Modifiers:  modifiers,
				Type:       typ,
				Identifier: identifier,
				Value:      value,
				Line:       p.currentToken().Line,
				Column:     p.currentToken().Column,
			}
		}
	}

	if p.currentTokenKind() == lexer.IDENTIFIER {
		return parseConstructor(p, modifiers)
	}

	panic(fmt.Sprintf("Unexpected token in class body: %s at line %d, column %d", lexer.TokenKindString(p.currentTokenKind()), p.currentToken().Line, p.currentToken().Column))
}

func parseMethod(p *parser, modifiers []ast.Modifier, returnType ast.Type, name string) ast.ClassMember {
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
		Line:       p.currentToken().Line,
		Column:     p.currentToken().Column,
	}
}

func parseConstructor(p *parser, modifiers []ast.Modifier) ast.ClassMember {
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
		Line:       p.currentToken().Line,
		Column:     p.currentToken().Column,
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
	startToken := p.expect(lexer.OPEN_BRACE)
	body := []ast.Stmt{}

	for p.currentTokenKind() != lexer.CLOSE_BRACE {
		body = append(body, parseStatement(p))
	}

	p.expect(lexer.CLOSE_BRACE)
	return ast.BlockStmt{
		Body:   body,
		Line:   startToken.Line,
		Column: startToken.Column,
	}
}
