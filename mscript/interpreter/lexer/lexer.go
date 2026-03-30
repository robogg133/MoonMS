package lexer

import (
	"strings"
	"unicode"
)

type Lexer struct {
	input    []rune
	position int
	next     int
	now      rune
}

func New(input string) *Lexer {

	l := &Lexer{
		input: []rune(input),
	}

	l.position = -1
	l.readRune()
	return l
}

func (l *Lexer) NextToken() Token {
	l.skipWhitespace()

	var tok Token

	switch l.now {
	case '=':
		if l.peekRune() == '=' {
			ch := l.now
			l.readRune()
			tok = Token{Type: TOKEN_EQUALS, Value: string(ch) + string(l.now)}
		} else {
			tok = Token{Type: TOKEN_ASSIGN, Value: "="}
		}

	case '+':
		if l.peekRune() == '+' {
			l.readRune()
			tok = Token{Type: TOKEN_INCREMENT, Value: "++"}
		} else if l.peekRune() == '=' {
			l.readRune()
			tok = Token{Type: TOKEN_PLUS_ASSIGN, Value: "+="}
		} else {
			tok = Token{Type: TOKEN_PLUS, Value: "+"}
		}

	case '-':
		if l.peekRune() == '-' {
			l.readRune()
			tok = Token{Type: TOKEN_DECREMENT, Value: "--"}
		} else if l.peekRune() == '=' {
			l.readRune()
			tok = Token{Type: TOKEN_MINUS_ASSIGN, Value: "-="}
		} else {
			tok = Token{Type: TOKEN_MINUS, Value: "-"}
		}

	case '*':
		tok = Token{Type: TOKEN_ASTERISK, Value: "*"}
	case '/':
		if l.peekRune() == '/' {
			for l.now != '\n' && l.now != 0 {
				l.readRune()
			}
			return l.NextToken()
		}
		tok = Token{Type: TOKEN_SLASH, Value: "/"}
	case '%':
		tok = Token{Type: TOKEN_REST, Value: "%"}
	case '>':
		if l.peekRune() == '=' {
			l.readRune()
			tok = Token{Type: TOKEN_GTE, Value: ">="}
		} else {
			tok = Token{Type: TOKEN_GREATER, Value: ">"}
		}

	case '<':
		if l.peekRune() == '=' {
			l.readRune()
			tok = Token{Type: TOKEN_LTE, Value: "<="}
		} else {
			tok = Token{Type: TOKEN_LESS, Value: "<"}
		}
	case '!':
		if l.peekRune() == '=' {
			l.readRune()
			tok = Token{Type: TOKEN_NOT_EQ, Value: "!="}
		} else {
			tok = Token{Type: TOKEN_NOT, Value: "!"}
		}
	case '"':
		tok.Type = TOKEN_STRING
		tok.Value = l.readString()
		return tok

	case '.':
		tok = Token{Type: TOKEN_DOT, Value: "."}

	case '(':
		tok = Token{Type: TOKEN_LPAREN, Value: "("}
	case ')':
		tok = Token{Type: TOKEN_RPAREN, Value: ")"}

	case '{':
		tok = Token{Type: TOKEN_LBRACE, Value: "{"}
	case '}':
		tok = Token{Type: TOKEN_RBRACE, Value: "}"}

	case '[':
		tok = Token{Type: TOKEN_LBRACKET, Value: "["}
	case ']':
		tok = Token{Type: TOKEN_RBRACKET, Value: "]"}

	case ',':
		tok = Token{Type: TOKEN_COMMA, Value: ","}

	case 0:
		tok = Token{Type: TOKEN_EOF, Value: "EOF"}

	case ';':
		tok = Token{Type: TOKEN_SEMICOLON, Value: ";"}
	case ':':
		tok = Token{Type: TOKEN_COLON, Value: ":"}

	default:
		if isLetter(l.now) {
			literal := l.readIdentifier()
			tok.Type = lookupIdent(literal)
			tok.Value = literal
			return tok
		} else if unicode.IsDigit(l.now) {

			tok.Value = l.readNumberOrDuration()
			if looksLikeDuration(tok.Value) {
				tok.Type = TOKEN_DURATION
				tok.Value = strings.TrimSuffix(tok.Value, "t")
			} else {
				tok.Type = TOKEN_NUMBER
			}

			return tok
		} else {
			tok = Token{Type: TOKEN_ILLEGAL, Value: string(l.now)}
		}
	}

	l.readRune()
	return tok
}

func (l *Lexer) skipWhitespace() {
	for l.now == ' ' || l.now == '\n' || l.now == '\t' || l.now == '\r' {
		l.readRune()
	}
}

func (l *Lexer) readString() string {
	var builder strings.Builder

	l.readRune()

	for l.now != '"' && l.now != 0 {

		if l.now == '\\' {
			l.readRune()

			switch l.now {
			case '"':
				builder.WriteRune('"')
			case '\\':
				builder.WriteRune('\\')
			case 'n':
				builder.WriteRune('\n')
			case 't':
				builder.WriteRune('\t')
			default:
				builder.WriteRune(l.now)
			}

		} else {
			builder.WriteRune(l.now)
		}

		l.readRune()
	}
	l.readRune()
	return builder.String()
}

func (l *Lexer) readIdentifier() string {
	start := l.position

	if !isLetter(l.input[start]) {
		panic("identifier need to begin with letter")
	}

	for isLetter(l.now) || unicode.IsDigit(l.now) {
		//fmt.Printf("%d: %s\n", l.position, string(l.now))
		l.readRune()
	}

	return string(l.input[start:l.position])
}

func (l *Lexer) readNumberOrDuration() string {
	start := l.position
	hasDot := false

	for unicode.IsDigit(l.now) || (l.now == '.' && !hasDot) {
		if l.now == '.' {
			hasDot = true
		}
		l.readRune()
	}

	if l.now == 't' {
		l.readRune()
	}

	return string(l.input[start:l.position])
}

func looksLikeDuration(s string) bool {
	return strings.HasSuffix(s, "t")
}

func (l *Lexer) readRune() {
	if l.next >= len(l.input) {
		l.now = 0
		l.position = l.next
		return
	}
	l.now = l.input[l.next]
	l.position++
	l.next++
}

func (l *Lexer) peekRune() rune {
	if l.next >= len(l.input) {
		return 0
	}
	return l.input[l.next]
}
