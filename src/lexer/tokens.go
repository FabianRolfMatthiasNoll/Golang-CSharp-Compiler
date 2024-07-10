package lexer

import "fmt"

type TokenKind int

// Dev Notes: Iota is making this basically an enum
const (
	EOF TokenKind = iota
	NUMBER
	STRING
	IDENTIFIER

	OPEN_BRACKET
	CLOSE_BRACKET
	OPEN_PAREN
	CLOSE_PAREN
	OPEN_BRACE
	CLOSE_BRACE

	ASSIGNMENT            // =
	EQUALS                // ==
	NOT                   // !
	NOT_EQUALS            // !=
	LESS_THAN             // <
	LESS_THAN_OR_EQUAL    // <=
	GREATER_THAN          // >
	GREATER_THAN_OR_EQUAL // >=

	// Operators
	PLUS     // +
	MINUS    // -
	MULTIPLY // *
	DIVIDE   // /
	MODULUS  // %

	DOT       // .
	SEMICOLON // ;
	COMMA     // ,
	COLON     // :

	PLUS_EQUALS     // +=
	MINUS_EQUALS    // -=
	MULTIPLY_EQUALS // *=
	DIVIDE_EQUALS   // /=
	MODULUS_EQUALS  // %=

	// Keywords
	IF
	ELSE
	FOR
	WHILE
	DO
	SWITCH
	CASE
	DEFAULT
	BREAK
	CONTINUE
	RETURN
	TRUE
	FALSE
	NULL

	// Types
	INT
	FLOAT
	STRING_TYPE
	CHAR
	BOOL
	VOID
	STRUCT
	ENUM
	TYPEDEF
	CONST
	STATIC
	EXTERN
	PUBLIC
	PRIVATE
	PROTECTED
	INTERFACE
	CLASS
	IMPLEMENTS
	EXTENDS
	NEW
	DELETE
	THIS
	SUPER
	IMPORT
	PACKAGE
	EXPORT
	FOREACH
)

type Token struct {
	Kind  TokenKind
	Value string
}

func (token Token) isOneOfMany (expectedTokens ...TokenKind) bool {
	for _, expectedToken := range expectedTokens {
		if token.Kind == expectedToken {
			return true
		}
	}
	return false
}

func (token Token) Debug() {
	if token.isOneOfMany(NUMBER, STRING, IDENTIFIER) {
		fmt.Printf("%s (%s)\n", TokenKindString(token.Kind), token.Value)
	} else {
		fmt.Printf("%s ()\n", TokenKindString(token.Kind))
	}
}

func NewToken(kind TokenKind, value string) Token {
	return Token{Kind: kind, Value: value}
}

func TokenKindString(kind TokenKind) string {
	switch kind {
	case EOF:
		return "EOF"
	case NUMBER:
		return "NUMBER"
	case STRING:
		return "STRING"
	case IDENTIFIER:
		return "IDENTIFIER"
	case OPEN_BRACKET:
		return "OPEN_BRACKET"
	case CLOSE_BRACKET:
		return "CLOSE_BRACKET"
	case OPEN_PAREN:
		return "OPEN_PAREN"
	case CLOSE_PAREN:
		return "CLOSE_PAREN"
	case OPEN_BRACE:
		return "OPEN_BRACE"
	case CLOSE_BRACE:
		return "CLOSE_BRACE"
	case ASSIGNMENT:
		return "ASSIGNMENT"
	case EQUALS:
		return "EQUALS"
	case NOT:
		return "NOT"
	case NOT_EQUALS:
		return "NOT_EQUALS"
	case LESS_THAN:
		return "LESS_THAN"
	case LESS_THAN_OR_EQUAL:
		return "LESS_THAN_OR_EQUAL"
	case GREATER_THAN:
		return "GREATER_THAN"
	case GREATER_THAN_OR_EQUAL:
		return "GREATER_THAN_OR_EQUAL"
	case PLUS:
		return "PLUS"
	case MINUS:
		return "MINUS"
	case MULTIPLY:
		return "MULTIPLY"
	case DIVIDE:
		return "DIVIDE"
	case MODULUS:
		return "MODULUS"
	case DOT:
		return "DOT"
	case SEMICOLON:
		return "SEMICOLON"
	case COMMA:
		return "COMMA"
	case COLON:
		return "COLON"
	case PLUS_EQUALS:
		return "PLUS_EQUALS"
	case MINUS_EQUALS:
		return "MINUS_EQUALS"
	case MULTIPLY_EQUALS:
		return "MULTIPLY_EQUALS"
	case DIVIDE_EQUALS:
		return "DIVIDE_EQUALS"
	case MODULUS_EQUALS:
		return "MODULUS_EQUALS"
	case IF:
		return "IF"
	case ELSE:
		return "ELSE"
	case FOR:
		return "FOR"
	case WHILE:
		return "WHILE"
	case DO:
		return "DO"
	case SWITCH:
		return "SWITCH"
	case CASE:
		return "CASE"
	case DEFAULT:
		return "DEFAULT"
	case BREAK:
		return "BREAK"
	case CONTINUE:
		return "CONTINUE"
	case RETURN:
		return "RETURN"
	case TRUE:
		return "TRUE"
	case FALSE:
		return "FALSE"
	case NULL:
		return "NULL"
	case INT:
		return "INT"
	case FLOAT:
		return "FLOAT"
	case STRING_TYPE:
		return "STRING_TYPE"
	case CHAR:
		return "CHAR"
	case BOOL:
		return "BOOL"
	case VOID:
		return "VOID"
	case STRUCT:
		return "STRUCT"
	case ENUM:
		return "ENUM"
	default:
		return "default"
	}
}
