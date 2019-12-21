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

	l.skipWhitespace()

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
	default:
		if isLetter(l.char) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.char) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.char)
		}
	}

	l.readNextChar()
	return tok
}

func newToken(tokenType token.TokenType, char byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(char)}
}

func (l *Lexer) readIdentifier() string {
	start := l.currIdx
	for isLetter(l.char) {
		l.readNextChar()
	}
	return l.input[start:l.currIdx]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (l *Lexer) skipWhitespace() {
	for l.char == ' ' || l.char == '\t' || l.char == '\n' || l.char == '\r' {
		l.readNextChar()
	}
}

func (l *Lexer) readNumber() string {
	start := l.currIdx
	for isDigit(l.char) {
		l.readNextChar()
	}
	return l.input[start:l.currIdx]
}

func isDigit(ch byte) bool {
	return '0' <= ch || ch <= '9'
}
