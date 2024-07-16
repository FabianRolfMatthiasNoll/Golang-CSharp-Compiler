package lexer

import (
	"fmt"
	"regexp"
)

type regexHandler func(lex *lexer, regex *regexp.Regexp)

type regexPattern struct {
	regex   *regexp.Regexp
	handler regexHandler
}

type lexer struct {
	patterns []regexPattern
	Tokens   []Token
	source   string
	pos      int
	line     int
	column   int
}

func (lex *lexer) advanceN(n int) {
	for i := 0; i < n; i++ {
		if lex.source[lex.pos] == '\n' {
			lex.line++
			lex.column = 0
		} else {
			lex.column++
		}
		lex.pos++
	}
}

func (lex *lexer) push(token Token) {
	lex.Tokens = append(lex.Tokens, token)
}

func (lex *lexer) remainder() string {
	return lex.source[lex.pos:]
}

func (lex *lexer) at_eof() bool {
	return lex.pos >= len(lex.source)
}

func Tokenize(source string) []Token {
	lex := createLexer(source)
	for !lex.at_eof() {
		// iterate while there are tokens left
		matched := false
		for _, pattern := range lex.patterns {
			loc := pattern.regex.FindStringIndex(lex.remainder())

			// Pattern must match current position or else it's the wrong pattern to use
			// example: 10 + 5 must match the NUMBER pattern first before the PLUS pattern
			if loc != nil && loc[0] == 0 {
				matched = true
				pattern.handler(lex, pattern.regex)
				break
			}
		}

		if !matched {
			panic(fmt.Sprintf("Lexer::Error -> unrecognized token at line %d, column %d near '%s'", lex.line, lex.column, lex.remainder()))
		}
	}
	lex.push(NewToken(EOF, "EOF", lex.line, lex.column))
	return lex.Tokens
}

func createLexer(source string) *lexer {
	return &lexer{
		source: source,
		pos:    0,
		line:   1,
		column: 1,
		Tokens: make([]Token, 0),
		patterns: []regexPattern{
			// Definition of all patterns
			// The Order is really important because the lexer will try to match the first pattern first
			{regexp.MustCompile(`^[0-9]+(\.[0-9]+)?`), numberHandler},
			{regexp.MustCompile(`^"[^"]*"`), stringHandler},
			{regexp.MustCompile(`^'[^']'`), charHandler},
			{regexp.MustCompile(`^\/\/.*`), skipHandler}, // skip comments
			{regexp.MustCompile(`^\s+`), skipHandler},    // skip whitespace
			{regexp.MustCompile(`^\(`), defaultHandler(OPEN_PAREN, "(")},
			{regexp.MustCompile(`^\)`), defaultHandler(CLOSE_PAREN, ")")},
			{regexp.MustCompile(`^\{`), defaultHandler(OPEN_BRACE, "{")},
			{regexp.MustCompile(`^\}`), defaultHandler(CLOSE_BRACE, "}")},
			{regexp.MustCompile(`^\[`), defaultHandler(OPEN_BRACKET, "[")},
			{regexp.MustCompile(`^\]`), defaultHandler(CLOSE_BRACKET, "]")},
			{regexp.MustCompile(`^\==`), defaultHandler(EQUALS, "==")},
			{regexp.MustCompile(`^\=`), defaultHandler(ASSIGNMENT, "=")},
			{regexp.MustCompile(`^\!=`), defaultHandler(NOT_EQUALS, "!=")},
			{regexp.MustCompile(`^\!`), defaultHandler(NOT, "!")},
			{regexp.MustCompile(`^\<=`), defaultHandler(LESS_THAN_OR_EQUAL, "<=")},
			{regexp.MustCompile(`^\>=`), defaultHandler(GREATER_THAN_OR_EQUAL, ">=")},
			{regexp.MustCompile(`^\<`), defaultHandler(LESS_THAN, "<")},
			{regexp.MustCompile(`^\>`), defaultHandler(GREATER_THAN, ">")},
			{regexp.MustCompile(`^\+\=`), defaultHandler(PLUS_EQUALS, "+=")},
			{regexp.MustCompile(`^\-\=`), defaultHandler(MINUS_EQUALS, "-=")},
			{regexp.MustCompile(`^\*\=`), defaultHandler(MULTIPLY_EQUALS, "*=")},
			{regexp.MustCompile(`^\/\=`), defaultHandler(DIVIDE_EQUALS, "/=")},
			{regexp.MustCompile(`^\%\=`), defaultHandler(MODULUS_EQUALS, "%=")},
			{regexp.MustCompile(`^\+\+`), defaultHandler(INCREMENT, "++")},
			{regexp.MustCompile(`^\-\-`), defaultHandler(DECREMENT, "--")},
			{regexp.MustCompile(`^\+`), defaultHandler(PLUS, "+")},
			{regexp.MustCompile(`^\-`), defaultHandler(MINUS, "-")},
			{regexp.MustCompile(`^\*`), defaultHandler(MULTIPLY, "*")},
			{regexp.MustCompile(`^\/`), defaultHandler(DIVIDE, "/")},
			{regexp.MustCompile(`^\%`), defaultHandler(MODULUS, "%")},
			{regexp.MustCompile(`^\.`), defaultHandler(DOT, ".")},
			{regexp.MustCompile(`^\;`), defaultHandler(SEMICOLON, ";")},
			{regexp.MustCompile(`^\:`), defaultHandler(COLON, ":")},
			{regexp.MustCompile(`^\,`), defaultHandler(COMMA, ",")},
			{regexp.MustCompile(`^\&\&`), defaultHandler(AND, "&&")},
			{regexp.MustCompile(`^\|\|`), defaultHandler(OR, "||")},
			{regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*`), identifierHandler},
			// Keywords will be matched as identifiers and converted in the handler
		},
	}
}

func defaultHandler(kind TokenKind, value string) regexHandler {
	return func(lex *lexer, regex *regexp.Regexp) {
		// Advance the lexers position past the end of the match
		lex.push(NewToken(kind, value, lex.line, lex.column))
		lex.advanceN(len(value))
	}
}

func numberHandler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindString(lex.remainder())
	lex.push(NewToken(INTLITERAL, match, lex.line, lex.column))
	lex.advanceN(len(match))
}

func skipHandler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindStringIndex(lex.remainder())
	lex.advanceN(match[1])
}

func stringHandler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindString(lex.remainder())
	match = match[1 : len(match)-1] // remove quotes
	lex.push(NewToken(STRINGLITERAL, match, lex.line, lex.column))
	lex.advanceN(len(match) + 2) // +2 for the quotes
}

func charHandler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindString(lex.remainder())
	match = match[1 : len(match)-1] // remove quotes
	lex.push(NewToken(CHARLITERAL, match, lex.line, lex.column))
	lex.advanceN(len(match) + 2) // +2 for the quotes
}

func identifierHandler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindString(lex.remainder())
	if kind, exists := keywords[match]; exists {
		lex.push(NewToken(kind, match, lex.line, lex.column))
	} else {
		lex.push(NewToken(IDENTIFIER, match, lex.line, lex.column))
	}
	lex.advanceN(len(match))
}
