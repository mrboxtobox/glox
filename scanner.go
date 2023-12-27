package main

import (
	"fmt"
	"strconv"
)

var reservedWords = map[string]TokenType{
	"and":    And,
	"class":  Class,
	"else":   Else,
	"false":  False,
	"for":    For,
	"fun":    Fun,
	"if":     If,
	"nil":    Nil,
	"or":     Or,
	"print":  Print,
	"return": Return,
	"super":  Super,
	"this":   This,
	"true":   True,
	"var":    Var,
	"while":  While,
}

func newScanner(source string) *Scanner {
	return &Scanner{
		source:  source,
		start:   0,
		current: 0,
		line:    1,
	}
}

// TODO: Need to makes sure the source text is an ASCII character.
type Scanner struct {
	// Raw source code.
	source string
	// First character in lexeme being scanned.
	start int
	// Character currently being considered.
	current int
	// What source line `current` is on. For producing tokens that know their location.
	line   int
	tokens []Token
}

func (s Scanner) scanTokens() []Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, Token{
		TokenType: EOF,
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
		s.addToken(LeftParen)
	case ')':
		s.addToken(RightParen)
	case '{':
		s.addToken(LeftBrace)
	case '}':
		s.addToken(RightBrace)
	case ',':
		s.addToken(Comma)
	case '.':
		s.addToken(Dot)
	case '-':
		s.addToken(Minus)
	case '+':
		s.addToken(Plus)
	case ';':
		s.addToken(Semicolon)
	case '*':
		s.addToken(Star)
	case '!':
		if s.match('=') {
			s.addToken(BangEqual)
		} else {
			s.addToken(Bang)
		}
	case '=':
		if s.match('=') {
			s.addToken(EqualEqual)
		} else {
			s.addToken(Equal)
		}
	case '<':
		if s.match('=') {
			s.addToken(LessEqual)
		} else {
			s.addToken(Less)
		}
	case '>':
		if s.match('=') {
			s.addToken(GreaterEqual)
		} else {
			s.addToken(Greater)
		}
	case '/':
		if s.match('/') {
			// Comment goes until the end of the line.
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(Slash)
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
	s.addTokenWithLiteral(String, value)
}

func (s *Scanner) scanNumber() {
	// Consume the rest of the digits.
	for isDigit(s.peek()) {
		s.advance()
	}

	// Check for the fractional part.
	if s.peek() == '.' && isDigit(s.peekNext()) {
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
	s.addTokenWithLiteral(Number, n)
}

func (s *Scanner) scanIdentifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}

	// TODO: Make substring conversation its own function.
	text := string([]rune(s.source)[s.start:s.current])
	TokenType, ok := reservedWords[text]
	if !ok {
		TokenType = Identifier
	}
	s.addToken(TokenType)
}

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}
	// TODO: How efficient is this []rune cast?
	current := []rune(s.source)[s.current]
	if current != expected {
		return false
	}
	current++
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
	// TODO: This advances by one rune?
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
