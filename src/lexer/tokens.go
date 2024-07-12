package lexer

import "fmt"

type TokenKind int

// Dev Notes: Iota is making this basically an enum
const (
	EOF TokenKind = iota
	NUMBER
	STRING
	IDENTIFIER
	OPEN_BRACKET          // [
	CLOSE_BRACKET         // ]
	OPEN_PAREN            // (
	CLOSE_PAREN           // )
	OPEN_BRACE            // {
	CLOSE_BRACE           // }
	ASSIGNMENT            // =
	EQUALS                // ==
	NOT                   // !
	NOT_EQUALS            // !=
	LESS_THAN             // <
	LESS_THAN_OR_EQUAL    // <=
	GREATER_THAN          // >
	GREATER_THAN_OR_EQUAL // >=
	PLUS                  // +
	MINUS                 // -
	MULTIPLY              // *
	DIVIDE                // /
	MODULUS               // %
	DOT                   // .
	SEMICOLON             // ;
	COMMA                 // ,
	COLON                 // :
	PLUS_EQUALS           // +=
	MINUS_EQUALS          // -=
	MULTIPLY_EQUALS       // *=
	DIVIDE_EQUALS         // /=
	MODULUS_EQUALS        // %=
	AND                   // &&
	OR                    // ||
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
	NEW
	THIS
	BASE
	IMPORT
	NAMESPACE
	USING
	CLASS
	STRUCT
	INTERFACE
	ENUM
	PUBLIC
	FINAL
	PRIVATE
	PROTECTED
	INTERNAL
	INCREMENT
	DECREMENT
	STATIC
	CONST
	VOID
	VAR
	BOOL
	CHAR
	INT
	FLOAT
	DOUBLE
	STRING_TYPE
)

var keywords = map[string]TokenKind{
	"if":        IF,
	"else":      ELSE,
	"for":       FOR,
	"while":     WHILE,
	"do":        DO,
	"switch":    SWITCH,
	"case":      CASE,
	"default":   DEFAULT,
	"break":     BREAK,
	"continue":  CONTINUE,
	"return":    RETURN,
	"true":      TRUE,
	"false":     FALSE,
	"final":     FINAL,
	"null":      NULL,
	"new":       NEW,
	"this":      THIS,
	"base":      BASE,
	"import":    IMPORT,
	"namespace": NAMESPACE,
	"using":     USING,
	"class":     CLASS,
	"struct":    STRUCT,
	"interface": INTERFACE,
	"enum":      ENUM,
	"public":    PUBLIC,
	"private":   PRIVATE,
	"protected": PROTECTED,
	"internal":  INTERNAL,
	"static":    STATIC,
	"const":     CONST,
	"void":      VOID,
	"var":       VAR,
	"bool":      BOOL,
	"char":      CHAR,
	"int":       INT,
	"float":     FLOAT,
	"double":    DOUBLE,
	"string":    STRING_TYPE,
}

type Token struct {
	Kind   TokenKind
	Value  string
	Line   int
	Column int
}

func (token Token) String() string {
	if token.isOneOfMany(NUMBER, STRING, IDENTIFIER) {
		return fmt.Sprintf("%s (%s) at %d:%d", TokenKindString(token.Kind), token.Value, token.Line, token.Column)
	}
	return fmt.Sprintf("%s at %d:%d", TokenKindString(token.Kind), token.Line, token.Column)
}

func (token Token) isOneOfMany(expectedTokens ...TokenKind) bool {
	for _, expectedToken := range expectedTokens {
		if token.Kind == expectedToken {
			return true
		}
	}
	return false
}

func (token Token) Debug() {
	if token.isOneOfMany(NUMBER, STRING, IDENTIFIER) {
		fmt.Printf("%s (%s) at %d:%d\n", TokenKindString(token.Kind), token.Value, token.Line, token.Column)
	} else {
		fmt.Printf("%s at %d:%d\n", TokenKindString(token.Kind), token.Line, token.Column)
	}
}

func NewToken(kind TokenKind, value string, line int, column int) Token {
	return Token{Kind: kind, Value: value, Line: line, Column: column}
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
	case NEW:
		return "NEW"
	case THIS:
		return "THIS"
	case BASE:
		return "BASE"
	case IMPORT:
		return "IMPORT"
	case NAMESPACE:
		return "NAMESPACE"
	case USING:
		return "USING"
	case CLASS:
		return "CLASS"
	case STRUCT:
		return "STRUCT"
	case INTERFACE:
		return "INTERFACE"
	case ENUM:
		return "ENUM"
	case PUBLIC:
		return "PUBLIC"
	case PRIVATE:
		return "PRIVATE"
	case PROTECTED:
		return "PROTECTED"
	case INTERNAL:
		return "INTERNAL"
	case STATIC:
		return "STATIC"
	case CONST:
		return "CONST"
	case VOID:
		return "VOID"
	case VAR:
		return "VAR"
	case BOOL:
		return "BOOL"
	case CHAR:
		return "CHAR"
	case INT:
		return "INT"
	case FLOAT:
		return "FLOAT"
	case DOUBLE:
		return "DOUBLE"
	case STRING_TYPE:
		return "STRING_TYPE"
	case AND:
		return "AND"
	case OR:
		return "OR"
	case FINAL:
		return "FINAL"
	case INCREMENT:
		return "INCREMENT"
	case DECREMENT:
		return "DECREMENT"
	default:
		return "default"
	}
}
