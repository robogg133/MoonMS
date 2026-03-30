package lexer

import "unicode"

func lookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return TOKEN_IDENT
}

func isLetter(c rune) bool {
	return unicode.IsLetter(c) || c == '_'
}
