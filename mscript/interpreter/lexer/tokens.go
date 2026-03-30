package lexer

type TokenType uint8

const (
	// Special
	TOKEN_ILLEGAL TokenType = iota
	TOKEN_EOF

	// Literals
	TOKEN_IDENT
	TOKEN_NUMBER
	TOKEN_STRING

	// Operators
	TOKEN_ASSIGN       // =
	TOKEN_EQUALS       // ==
	TOKEN_NOT          // !
	TOKEN_NOT_EQ       // !=
	TOKEN_NOT_KW       // not
	TOKEN_GREATER      // >
	TOKEN_LESS         //
	TOKEN_GTE          // >=
	TOKEN_LTE          // <=
	TOKEN_PLUS         // +
	TOKEN_PLUS_ASSIGN  // +=
	TOKEN_INCREMENT    // ++
	TOKEN_MINUS        // -
	TOKEN_MINUS_ASSIGN // -=
	TOKEN_DECREMENT    // --
	TOKEN_ASTERISK     // *
	TOKEN_SLASH        // /
	TOKEN_REST         // %

	// Delimiters
	TOKEN_DOT       // .
	TOKEN_COMMA     // ,
	TOKEN_COLON     // :
	TOKEN_SEMICOLON // ;
	TOKEN_LPAREN    // (
	TOKEN_RPAREN    // )
	TOKEN_LBRACE    // {
	TOKEN_RBRACE    // }
	TOKEN_LBRACKET  // [
	TOKEN_RBRACKET  // ]

	// Keywords
	TOKEN_IF
	TOKEN_ELSE
	TOKEN_ELSE_IF
	TOKEN_RETURN
	TOKEN_FN
	TOKEN_LET
	TOKEN_PLUGIN
	TOKEN_VERSION
	TOKEN_EXPORT
	TOKEN_EMITS
	TOKEN_REQUIRE
	TOKEN_INCLUDE
	TOKEN_COMMAND
	TOKEN_ON
	TOKEN_EMIT
	TOKEN_AFTER
	TOKEN_EVERY
	TOKEN_TRUE
	TOKEN_FALSE
	TOKEN_NULL
	TOKEN_IS
	TOKEN_AS
	TOKEN_NOT_KEYWORD
	TOKEN_CANCEL
	TOKEN_DURATION
)

var keywords = map[string]TokenType{
	"if":      TOKEN_IF,
	"else":    TOKEN_ELSE,
	"else if": TOKEN_ELSE_IF,
	"return":  TOKEN_RETURN,
	"fn":      TOKEN_FN,
	"let":     TOKEN_LET,
	"plugin":  TOKEN_PLUGIN,
	"version": TOKEN_VERSION,
	"export":  TOKEN_EXPORT,
	"emits":   TOKEN_EMITS,
	"require": TOKEN_REQUIRE,
	"include": TOKEN_INCLUDE,
	"command": TOKEN_COMMAND,
	"on":      TOKEN_ON,
	"emit":    TOKEN_EMIT,
	"after":   TOKEN_AFTER,
	"every":   TOKEN_EVERY,
	"true":    TOKEN_TRUE,
	"false":   TOKEN_FALSE,
	"null":    TOKEN_NULL,
	"as":      TOKEN_AS,
	"not":     TOKEN_NOT_KW,
	"cancel":  TOKEN_CANCEL,
}

type Token struct {
	Type  TokenType
	Value string
}
