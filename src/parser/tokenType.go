package parser

type TokenType uint16

const (
	// Single-character tokens
	At TokenType = iota
	Caret
	Colon
	Comma
	Dot
	LeftBrace
	LeftBracket
	LeftParen
	Minus
	Percent
	Plus
	RightBrace
	RightBracket
	RightParen
	Semicolon
	Slash
	Star

	// One or two character tokens
	Bang
	BangEqual
	Equal
	EqualEqual
	Greater
	GreaterEqual
	Less
	LessEqual

	// Literals
	Identifier
	Number
	String

	// Keywords
	And
	Assert
	Break
	Class
	Echo
	Else
	False
	Fn
	For
	Foreign
	If
	Import
	Nil
	Or
	Return
	True
	Var
	While

	Eof
	Newline

	Error
)
