package parser

import "unicode"

func isDigit(r rune) bool {
	return unicode.IsDigit(r)
}

func isLetter(r rune) bool {
	return unicode.IsLetter(r)
}

func isAlpha(r rune) bool {
	return isLetter(r) || r == '_'
}
