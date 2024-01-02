package main

import (
	"fmt"
	"strconv"
)

var reservedWords = map[string]TokenType{
	"and":    AndToken,
	"class":  ClassToken,
	"else":   ElseToken,
	"false":  FalseToken,
	"for":    ForToken,
	"fun":    FunToken,
	"if":     IfToken,
	"nil":    NilToken,
	"or":     OrToken,
	"print":  PrintToken,
	"return": ReturnToken,
	"super":  SuperToken,
	"this":   ThisToken,
	"true":   TrueToken,
	"var":    VarToken,
	"while":  WhileToken,
}

// Creates a new scanner.
func NewScanner(source string) *Scanner {
	return &Scanner{
		source:  source,
		start:   0,
		current: 0,
		line:    1,
	}
}

type Scanner struct {
	// Raw source code.
	source string
	// First character in lexeme being scanned.
	start int
	// Character currently being scanned.
	current int
	// What source line `current` is on. Used to assign line location to tokens.
	line int
	// Scanned tokens.
	tokens []Token
}

// Returns an ordered list of Tokens from scanning the source.
func (s *Scanner) ScanTokens() []Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, Token{
		TokenType: EOFToken,
		Lexeme:    "",
		Literal:   nil,
		Line:      s.line,
	})
	return s.tokens
}

func (s *Scanner) scanToken() {
	c := s.advance()
	// Handle the single characters first.
	switch c {
	case '(':
		s.addToken(LeftParenToken)
	case ')':
		s.addToken(RightParenToken)
	case '{':
		s.addToken(LeftBraceToken)
	case '}':
		s.addToken(RightBraceToken)
	case ',':
		s.addToken(CommaToken)
	case '.':
		s.addToken(DotToken)
	case '-':
		s.addToken(MinusToken)
	case '+':
		s.addToken(PlusToken)
	case ';':
		s.addToken(SemicolonToken)
	case '*':
		s.addToken(StarToken)
	case '!':
		if s.match('=') {
			s.addToken(BangEqualToken)
		} else {
			s.addToken(BangToken)
		}
	case '=':
		if s.match('=') {
			s.addToken(EqualEqualToken)
		} else {
			s.addToken(EqualToken)
		}
	case '<':
		if s.match('=') {
			s.addToken(LessEqualToken)
		} else {
			s.addToken(LessToken)
		}
	case '>':
		if s.match('=') {
			s.addToken(GreaterEqualToken)
		} else {
			s.addToken(GreaterToken)
		}
	case '/':
		if s.match('/') {
			// Comment goes until the end of the line.
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(SlashToken)
		}
	case ' ':
	case '\r':
	case '\t':
		break
	case '\n':
		s.line++
	case '"':
		s.scanString()
	default:
		// Handling numbers is a bit tedious so put it in the default case and
		// handle the cases in isDigit(). This has the weird case that
		// print(-123.abs()) will give you -123 since -123 is considered as negation
		// applied to the number 123 and -123.abs() is really -(123.abs()).
		if isDigit(c) {
			println("handling digit ->", s.source[s.start:s.current])
			s.scanNumber()
		} else if isAlpha(c) {
			s.scanIdentifier()
		} else {
			printErr(s.line, fmt.Sprintf("Unexpected character: %c.", c))
		}
	}
}

func (s *Scanner) scanString() {
	for s.peek() != '"' && !s.isAtEnd() {
		// glox supports strings.
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		printErr(s.line, "Unterminated string.")
		return
	}

	// Closing '"'
	s.advance()
	// Trim the surround quotes
	value := string([]rune(s.source)[s.start+1 : s.current-1])
	s.addTokenWithLiteral(StringToken, value)
}

func (s *Scanner) scanNumber() {
	println("scannign number ->", s.source[s.start:s.current])
	println("s.peek(), s.peekNext()", string([]rune{s.peek()}), string([]rune{s.peekNext()}))
	fmt.Printf("'%c' -> '%c'\n", s.peek(), s.peekNext())
	// Consume the rest of the digits.
	for isDigit(s.peek()) {
		s.advance()
	}

	// Check for the fractional part.
	if s.peek() == '.' && isDigit(s.peekNext()) {
		s.advance()
		// Consume the '.'
		for isDigit(s.peek()) {
			s.advance()
		}
	}

	text := string([]rune(s.source)[s.start:s.current])
	n, err := strconv.ParseFloat(text, 64)
	if err != nil {
		printErr(s.line, fmt.Sprintf("Unable to parse string: %v.", err))
	}
	s.addTokenWithLiteral(NumberToken, n)
}

func (s *Scanner) scanIdentifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}

	// TODO: Refactor substring conversion into its own function.
	text := string([]rune(s.source)[s.start:s.current])
	tokenType, ok := reservedWords[text]
	if !ok {
		tokenType = IdentifierToken
	}
	s.addToken(tokenType)
}

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}
	// TODO: Figure out how efficient this []rune cast is.
	currRune := []rune(s.source)[s.current]
	if currRune != expected {
		return false
	}
	s.current++
	return true
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return '\000'
	}
	return []rune(s.source)[s.current]
}

// peekNext does lookahead by 2 characters. It is useful when parsing decimals.
// We don't want to consume a '.' unless we're sure it is followed by a digit.
func (s *Scanner) peekNext() rune {
	if s.current+1 >= len(s.source) {
		return '\000'
	}
	return []rune(s.source)[s.current+1]
}

// isAtEnd checks whether we have consumed all characters in `source`.
func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

// advance consumes and returns the next character.
func (s *Scanner) advance() rune {
	// TODO: ThisToken advances by one rune?
	c := []rune(s.source)[s.current]
	s.current++
	return c
}

func (s *Scanner) addToken(t TokenType) {
	s.addTokenWithLiteral(t, nil)
}

func (s *Scanner) addTokenWithLiteral(t TokenType, literal interface{}) {
	text := []rune(s.source)[s.start:s.current]
	s.tokens = append(s.tokens, Token{
		TokenType: t,
		Lexeme:    string(text),
		Literal:   literal,
		Line:      s.line,
	})
}

func isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func isAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func isAlphaNumeric(c rune) bool {
	return isDigit(c) || isAlpha(c)
}
