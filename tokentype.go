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
		return "LeftParenToken"
	case RightParenToken:
		return "RightParenToken"
	case LeftBraceToken:
		return "LeftBraceToken"
	case RightBraceToken:
		return "RightBraceToken"
	case CommaToken:
		return "CommaToken"
	case DotToken:
		return "DotToken"
	case MinusToken:
		return "MinusToken"
	case PlusToken:
		return "PlusToken"
	case SemicolonToken:
		return "SemicolonToken"
	case SlashToken:
		return "SlashToken"
	case StarToken:
		return "StarToken"
	case BangToken:
		return "BangToken"
	case BangEqualToken:
		return "BangEqualToken"
	case EqualToken:
		return "EqualToken"
	case EqualEqualToken:
		return "EqualEqualToken"
	case GreaterToken:
		return "GreaterToken"
	case GreaterEqualToken:
		return "GreaterEqualToken"
	case LessToken:
		return "LessToken"
	case LessEqualToken:
		return "LessEqualToken"
	case IdentifierToken:
		return "IdentifierToken"
	case StringToken:
		return "StringToken"
	case NumberToken:
		return "NumberToken"
	case AndToken:
		return "AndToken"
	case ClassToken:
		return "ClassToken"
	case ElseToken:
		return "ElseToken"
	case FalseToken:
		return "FalseToken"
	case FunToken:
		return "FunToken"
	case ForToken:
		return "ForToken"
	case IfToken:
		return "IfToken"
	case NilToken:
		return "NilToken"
	case OrToken:
		return "OrToken"
	case PrintToken:
		return "PrintToken"
	case ReturnToken:
		return "ReturnToken"
	case SuperToken:
		return "SuperToken"
	case ThisToken:
		return "ThisToken"
	case TrueToken:
		return "TrueToken"
	case VarToken:
		return "VarToken"
	case WhileToken:
		return "WhileToken"
	case EOFToken:
		return "EOFToken"
	default:
		return "UnknownToken"
	}
}
