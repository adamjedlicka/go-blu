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

func (t Token) GetType() TokenType {
	return t.tokenType
}

func (t Token) GetLexeme() string {
	return t.lexeme
}

func (t Token) GetLine() int {
	return t.line
}

func (t Token) String() string {
	return "Token{" + string(t.tokenType) + "<" + t.lexeme + ">}"
}
