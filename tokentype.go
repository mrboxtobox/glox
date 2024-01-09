package main

type TokenType int

const (
	// Single-character tokens.
	LeftParenToken TokenType = iota
	RightParenToken
	LeftBraceToken
	RightBraceToken
	CommaToken
	DotToken
	MinusToken
	PlusToken
	SemicolonToken
	SlashToken
	StarToken

	// One or two character tokens.
	BangToken
	BangEqualToken
	EqualToken
	EqualEqualToken
	GreaterToken
	GreaterEqualToken
	LessToken
	LessEqualToken

	// Literals.
	IdentifierToken
	StringToken
	NumberToken

	// Keywords
	AndToken
	ClassToken
	ElseToken
	FalseToken
	FunToken
	ForToken
	IfToken
	NilToken
	OrToken
	PrintToken
	ReturnToken
	SuperToken
	ThisToken
	TrueToken
	VarToken
	WhileToken

	EOFToken
)

func (t TokenType) String() string {
	switch t {
	case LeftParenToken:
		return "LeftParen"
	case RightParenToken:
		return "RightParen"
	case LeftBraceToken:
		return "LeftBrace"
	case RightBraceToken:
		return "RightBrace"
	case CommaToken:
		return "Comma"
	case DotToken:
		return "Dot"
	case MinusToken:
		return "Minus"
	case PlusToken:
		return "Plus"
	case SemicolonToken:
		return "Semicolon"
	case SlashToken:
		return "Slash"
	case StarToken:
		return "Star"
	case BangToken:
		return "Bang"
	case BangEqualToken:
		return "BangEqual"
	case EqualToken:
		return "Equal"
	case EqualEqualToken:
		return "EqualEqual"
	case GreaterToken:
		return "Greater"
	case GreaterEqualToken:
		return "GreaterEqual"
	case LessToken:
		return "Less"
	case LessEqualToken:
		return "LessEqual"
	case IdentifierToken:
		return "Identifier"
	case StringToken:
		return "String"
	case NumberToken:
		return "Number"
	case AndToken:
		return "And"
	case ClassToken:
		return "Class"
	case ElseToken:
		return "Else"
	case FalseToken:
		return "False"
	case FunToken:
		return "Fun"
	case ForToken:
		return "For"
	case IfToken:
		return "If"
	case NilToken:
		return "Nil"
	case OrToken:
		return "Or"
	case PrintToken:
		return "Print"
	case ReturnToken:
		return "Return"
	case SuperToken:
		return "Super"
	case ThisToken:
		return "This"
	case TrueToken:
		return "True"
	case VarToken:
		return "Var"
	case WhileToken:
		return "While"
	case EOFToken:
		return "EOF"
	default:
		return "UnknownToken"
	}
}
