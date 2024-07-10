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
}

func (lex *lexer) advanceN(n int) {
	lex.pos += n
}

func (lex *lexer) push(token Token) {
	lex.Tokens = append(lex.Tokens, token)
}

func (lex *lexer) at() byte {
	return lex.source[lex.pos]
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

			// Pattern must match current position or else its the wrong pattern to use
			// example: 10 + 5 must match the NUMBER pattern first before the PLUS pattern
			if loc != nil && loc[0] == 0 {
				matched = true
				pattern.handler(lex, pattern.regex)
				break
			}
		}

		if !matched {
			panic(fmt.Sprintf("Lexer::Error -> unrecognized token at position %d near %s", lex.pos, lex.remainder()))
		}
	}
	lex.push(NewToken(EOF, "EOF"))
	return lex.Tokens
}

func createLexer(source string) *lexer {
	return &lexer{
		source: source,
		pos:    0,
		Tokens: make([]Token, 0),
		patterns: []regexPattern{
			// Definition of all patterns
			// The Order is really important because the lexer will try to match the first pattern first
			{regexp.MustCompile(`[0-9]+(\.[0-9]+)?`), numberHandler},
			{regexp.MustCompile(`"[^"]*"`), stringHandler},
			{regexp.MustCompile(`\/\/.*`), skipHandler}, // skip comments
			{regexp.MustCompile(`\s+`), skipHandler},  // skip whitespace
			{regexp.MustCompile(`\(`), defaultHandler(OPEN_PAREN, "(")},
			{regexp.MustCompile(`\)`), defaultHandler(CLOSE_PAREN, ")")},
			{regexp.MustCompile(`\{`), defaultHandler(OPEN_BRACE, "{")},
			{regexp.MustCompile(`\}`), defaultHandler(CLOSE_BRACE, "}")},
			{regexp.MustCompile(`\[`), defaultHandler(OPEN_BRACKET, "[")},
			{regexp.MustCompile(`\]`), defaultHandler(CLOSE_BRACKET, "]")},
			{regexp.MustCompile(`\==`), defaultHandler(EQUALS, "==")},
			{regexp.MustCompile(`\=`), defaultHandler(ASSIGNMENT, "=")},
			{regexp.MustCompile(`\!`), defaultHandler(NOT, "!")},
			{regexp.MustCompile(`\+`), defaultHandler(PLUS, "+")},
			{regexp.MustCompile(`\-`), defaultHandler(MINUS, "-")},
			{regexp.MustCompile(`\*`), defaultHandler(MULTIPLY, "*")},
			{regexp.MustCompile(`\/`), defaultHandler(DIVIDE, "/")},
			{regexp.MustCompile(`\!=`), defaultHandler(NOT_EQUALS, "!=")},
			{regexp.MustCompile(`\<\=`), defaultHandler(LESS_THAN_OR_EQUAL, "<=")},
			{regexp.MustCompile(`\>\=`), defaultHandler(GREATER_THAN_OR_EQUAL, ">=")},
			{regexp.MustCompile(`\<`), defaultHandler(LESS_THAN, "<")},
			{regexp.MustCompile(`\>`), defaultHandler(GREATER_THAN, ">")},
			{regexp.MustCompile(`\;`), defaultHandler(SEMICOLON, ";")},
		},
	}
}

func defaultHandler(kind TokenKind, value string) regexHandler {
	return func(lex *lexer, regex *regexp.Regexp) {
		// Advance the lexers position past the end of the match
		lex.advanceN(len(value))
		lex.push(NewToken(kind, value))
	}
}

func numberHandler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindString(lex.remainder())
	lex.push(NewToken(NUMBER, match))
	lex.advanceN(len(match))
}

func skipHandler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindStringIndex(lex.remainder())
	lex.advanceN(match[1])
}

func stringHandler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindString(lex.remainder())
	lex.push(NewToken(STRING, match))
	lex.advanceN(len(match))
}
