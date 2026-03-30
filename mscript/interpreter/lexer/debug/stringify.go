package lex_debug

import (
	"fmt"

	"github.com/robogg133/MoonMS/mscript/interpreter/lexer"
)

func Stringify(tk lexer.Token) string {

	switch tk.Type {
	case lexer.TOKEN_IDENT:
		return fmt.Sprintf("IDENT(%s)", tk.Value)
	case lexer.TOKEN_DOT:
		return "DOT"
	case lexer.TOKEN_LPAREN:
		return "LPAREN"
	case lexer.TOKEN_RPAREN:
		return "RPAREN"
	case lexer.TOKEN_NUMBER:
		return "NUMBER(" + tk.Value + ")"
	case lexer.TOKEN_COMMA:
		return "COMMA"
	case lexer.TOKEN_STRING:
		return fmt.Sprintf("STRING(%s)", tk.Value)
	case lexer.TOKEN_ILLEGAL:
		return fmt.Sprintf("ILLEGAL(%s)", tk.Value)
	case lexer.TOKEN_EQUALS:
		return "EQUALS"
	case lexer.TOKEN_EOF:
		return "EOF"
	case lexer.TOKEN_AFTER:
		return "AFTER"
	case lexer.TOKEN_PLUS:
		return "PLUS"
	case lexer.TOKEN_PLUS_ASSIGN:
		return "PLUS ASSIGN"
	case lexer.TOKEN_INCREMENT:
		return "INCREMENT"
	case lexer.TOKEN_TRUE:
		return "TRUE"
	case lexer.TOKEN_RETURN:
		return "KEYWORD RETURN"
	case lexer.TOKEN_RBRACE:
		return "RBRACE"
	case lexer.TOKEN_LBRACE:
		return "LBRACE"
	case lexer.TOKEN_IF:
		return "KEYWORD IF"
	case lexer.TOKEN_SEMICOLON:
		return "SEMICOLON"
	case lexer.TOKEN_PLUGIN:
		return "KEYWORD PLUGIN"
	case lexer.TOKEN_VERSION:
		return "KEYWORD PLUGIN VERSION"
	case lexer.TOKEN_EXPORT:
		return "KEYWORD EXPORT"
	case lexer.TOKEN_COLON:
		return "COLON"
	case lexer.TOKEN_FN:
		return "KEYWORD FN"
	case lexer.TOKEN_EMITS:
		return "EMITS"
	case lexer.TOKEN_ASSIGN:
		return "ASSIGN"
	case lexer.TOKEN_LBRACKET:
		return "LBRACKET"
	case lexer.TOKEN_RBRACKET:
		return "RBRACKET"
	default:
		return fmt.Sprintf("UNKNOWN FOR DEBUG: %d", tk.Type)
	}

}
