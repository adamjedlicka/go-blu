package parser

type Token struct {
	tokenType TokenType
	lexeme    string
	line      int
}

func NewToken(tokenType TokenType, lexeme string, line int) Token {
	return Token{
		tokenType: tokenType,
		lexeme:    lexeme,
		line:      line,
	}
}

func (t Token) Type() TokenType {
	return t.tokenType
}

func (t Token) Lexeme() string {
	return t.lexeme
}

func (t Token) Line() int {
	return t.line
}

func (t Token) String() string {
	return "Token{" + string(t.tokenType) + "<" + t.lexeme + ">}"
}
