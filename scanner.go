package main

func newScanner(source string) scanner {
	return scanner{
		source:  source,
		start:   0,
		current: 0,
		line:    1,
	}
}

// TODO: Need to makes sure the source text is an ASCII character.
type scanner struct {
	// Raw source code.
	source string
	// First character in lexeme being scanned.
	start int
	// Character currently being considered.
	current int
	// What source line `current` is on. For producing tokens that know their location.
	line   int
	tokens []token
}

func (s scanner) scanTokens() []token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, token{
		tokenType: EOF,
		lexeme:    "",
		literal:   nil,
		line:      s.line,
	})
	return s.tokens
}

func (s scanner) scanToken() {
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
			// Comment does until the end of the line.
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
		printErr(s.line, "Unexpected character.")
	}
}

func (s scanner) scanString() {
	for s.peek() != '"' && !s.isAtEnd() {
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

func (s scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}
	current := []rune(s.source)[s.current]
	if current != expected {
		return false
	}
	current++
	return true
}

func (s scanner) peek() rune {
	if s.isAtEnd() {
		return '\000'
	}
	return []rune(s.source)[s.current]
}

// Whether we have consumed all characters in `source`.
func (s scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

// advance consumes and returns the next character.
func (s scanner) advance() rune {
	// TODO: This advances by one rune?
	c := []rune(s.source)[s.current]
	s.current++
	return c
}

func (s scanner) addToken(t tokenType) {
	s.addTokenWithLiteral(t, nil)
}

func (s scanner) addTokenWithLiteral(t tokenType, literal interface{}) {
	text := []rune(s.source)[s.start]
	s.tokens = append(s.tokens, token{
		tokenType: t,
		lexeme:    string(text),
		literal:   literal,
		line:      s.line,
	})
}
