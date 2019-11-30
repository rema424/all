package lexer

import "monkey/token"

type Lexer struct {
	input   string
	currIdx int
	nextIdx int
	char    byte
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readNextChar()
	return l
}

func (l *Lexer) readNextChar() {
	if l.nextIdx < len(l.input) {
		l.char = l.input[l.nextIdx]
	} else {
		l.char = 0
	}
	l.currIdx = l.nextIdx
	l.nextIdx += 1
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	switch l.char {
	case '=':
		tok = newToken(token.ASSIGN, l.char)
	case ';':
		tok = newToken(token.SEMICOLON, l.char)
	case '(':
		tok = newToken(token.LPAREN, l.char)
	case ')':
		tok = newToken(token.RPAREN, l.char)
	case ',':
		tok = newToken(token.COMMA, l.char)
	case '+':
		tok = newToken(token.PLUS, l.char)
	case '{':
		tok = newToken(token.LBRACE, l.char)
	case '}':
		tok = newToken(token.RBRACE, l.char)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	}

	l.readNextChar()
	return tok
}

func newToken(tokenType token.TokenType, char byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(char)}
}
