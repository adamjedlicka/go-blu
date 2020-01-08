package parser

type Parser struct {
	source   string
	from     int
	at       int
	lineFrom int
	lineTo   int

	previous Token
	current  Token
}

func NewParser(source string) Parser {
	return Parser{
		source:   source,
		from:     0,
		at:       0,
		lineFrom: 1,
		lineTo:   1,
	}
}

func (p *Parser) NextToken() Token {
	p.skipWhitespace()

	if p.isAtEnd() {
		return p.eof()
	}

	p.from = p.at
	p.lineFrom = p.lineTo

	r := p.advance()
	if isAlpha(r) {
		return p.identifier()
	}

	if isDigit(r) {
		return p.number()
	}

	switch r {
	case '(':
		return p.makeToken(LeftParen)
	case ')':
		return p.makeToken(RightParen)
	case '[':
		return p.makeToken(LeftBracket)
	case ']':
		return p.makeToken(RightBracket)
	case '{':
		return p.makeToken(LeftBrace)
	case '}':
		return p.makeToken(RightBrace)
	case ';':
		return p.makeToken(Semicolon)
	case '@':
		return p.makeToken(At)
	case '^':
		return p.makeToken(Caret)
	case ':':
		return p.makeToken(Colon)
	case ',':
		return p.makeToken(Comma)
	case '.':
		return p.makeToken(Dot)
	case '-':
		return p.makeToken(Minus)
	case '+':
		return p.makeToken(Plus)
	case '%':
		return p.makeToken(Percent)
	case '/':
		return p.makeToken(Slash)
	case '*':
		return p.makeToken(Star)
	case '!':
		if p.match('=') {
			return p.makeToken(BangEqual)
		}
		return p.makeToken(Bang)
	case '=':
		if p.match('=') {
			return p.makeToken(EqualEqual)
		}
		return p.makeToken(Equal)
	case '>':
		if p.match('=') {
			return p.makeToken(GreaterEqual)
		}
		return p.makeToken(Greater)
	case '<':
		if p.match('=') {
			return p.makeToken(LessEqual)
		}
		return p.makeToken(Less)
	case '"':
		return p.string()
	case '\n':
		return p.newline()
	}

	return p.error("Unexpected character.")
}

func (p *Parser) GetTokens() []Token {
	tokens := make([]Token, 0)

	for true {
		token := p.NextToken()

		tokens = append(tokens, token)

		if token.tokenType == Eof {
			break
		}
	}

	return tokens
}

func (p *Parser) Previous() Token {
	return p.previous
}

func (p *Parser) SetPrevious(previous Token) {
	p.previous = previous
}

func (p *Parser) Current() Token {
	return p.current
}

func (p *Parser) SetCurrent(token Token) {
	p.current = token
}

func (p *Parser) number() Token {
	for isDigit(p.peek()) {
		p.advance()
	}

	// Look for a fractional part
	if p.peek() == '.' && isDigit(p.peekNext()) {
		// Consume the "."
		p.advance()

		for isDigit(p.peek()) {
			p.advance()
		}
	}

	return p.makeToken(Number)
}

func (p *Parser) string() Token {
	for p.peek() != '"' && !p.isAtEnd() {
		if p.peek() == '\n' {
			p.newline()
		}

		p.advance()
	}

	if p.isAtEnd() {
		return p.error("Unterminated string.")
	}

	// The closing "
	p.advance()

	return p.makeToken(String)
}

func (p *Parser) identifier() Token {
	for isAlpha(p.peek()) || isDigit(p.peek()) {
		p.advance()
	}

	if tokenType, ok := keywords[p.source[p.from:p.at]]; ok {
		return p.makeToken(tokenType)
	}

	return p.makeToken(Identifier)
}

func (p *Parser) makeToken(tokenType TokenType) Token {
	switch tokenType {
	case Newline:
		return NewToken(tokenType, "<Newline>", p.lineFrom)
	case Eof:
		return NewToken(tokenType, "<Eof>", p.lineFrom)
	default:
		return NewToken(tokenType, p.source[p.from:p.at], p.lineFrom)
	}
}

func (p *Parser) eof() Token {
	return p.makeToken(Eof)
}

func (p *Parser) newline() Token {
	p.lineTo++

	return p.makeToken(Newline)
}

func (p *Parser) error(message string) Token {
	return NewToken(Error, message, p.lineFrom)
}

func (p *Parser) advance() rune {
	if p.isAtEnd() {
		return 0
	}

	p.at++

	return rune(p.source[p.at-1])
}

func (p *Parser) peek() rune {
	if p.isAtEnd() {
		return 0
	}

	return rune(p.source[p.at])
}

func (p *Parser) peekNext() rune {
	if p.isAtEnd() || p.isBeforeTheEnd() {
		return 0
	}

	return rune(p.source[p.at+1])
}

func (p *Parser) match(expected rune) bool {
	if p.isAtEnd() {
		return false
	}

	if p.peek() != expected {
		return false
	}

	p.advance()
	return true
}

func (p *Parser) isAtEnd() bool {
	return len(p.source) == p.at
}

func (p *Parser) isBeforeTheEnd() bool {
	return len(p.source) == p.at+1
}

func (p *Parser) skipWhitespace() {
	for true {
		r := p.peek()

		switch r {
		case ' ', '\r', '\t':
			p.advance()
		case '/':
			if p.peekNext() == '/' {
				// A comment goes until the end of the line
				for p.peek() != '\n' && !p.isAtEnd() {
					p.advance()
				}
			} else {
				return
			}
		default:
			return
		}
	}
}
