package lexer

import (
	"sepia/token"
	"sepia/util"
)

// Lexer is a structure that lexes a given input.
type Lexer struct {
	input           string
	position        int
	readingPosition int
	currentChar     byte
}

func (lexer *Lexer) consumeChar() {
	if lexer.readingPosition >= len(lexer.input) {
		lexer.currentChar = 0
	} else {
		lexer.currentChar = lexer.input[lexer.readingPosition]
	}

	lexer.position = lexer.readingPosition
	lexer.readingPosition++

}

// NextToken get the next token.
func (lexer *Lexer) NextToken() token.Token {
	var t token.Token

	lexer.skipWhitespace()

	switch lexer.currentChar {

	// One-character bytes
	case '(':
		t = newToken(token.LPAREN, lexer.currentChar)
	case ')':
		t = newToken(token.RPAREN, lexer.currentChar)
	case ',':
		t = newToken(token.COMMA, lexer.currentChar)
	case '+':
		switch lexer.peekCharacter() {
		case '=':
			character := lexer.currentChar
			lexer.consumeChar()

			t = token.Token{Type: token.PLUSEQ,
				Literal: string(character) + string(lexer.currentChar)}
		case '+':
			character := lexer.currentChar
			lexer.consumeChar()

			t = token.Token{Type: token.INCREMENT,
				Literal: string(character) + string(lexer.currentChar)}
		default:
			t = newToken(token.MINUS, lexer.currentChar)
		}
	case '{':
		t = newToken(token.LBRACE, lexer.currentChar)
	case '}':
		t = newToken(token.RBRACE, lexer.currentChar)
	case '[':
		t = newToken(token.LBRACKET, lexer.currentChar)
	case ']':
		t = newToken(token.RBRACKET, lexer.currentChar)
	case '-':

		switch lexer.peekCharacter() {
		case '>':
			character := lexer.currentChar
			lexer.consumeChar()

			t = token.Token{
				Type:    token.OPENBLOCK,
				Literal: string(character) + string(lexer.currentChar),
			}
		case '=':
			character := lexer.currentChar
			lexer.consumeChar()

			t = token.Token{Type: token.MINUSEQ,
				Literal: string(character) + string(lexer.currentChar)}
		case '-':
			character := lexer.currentChar
			lexer.consumeChar()

			t = token.Token{Type: token.DECREMENT,
				Literal: string(character) + string(lexer.currentChar)}
		default:
			t = newToken(token.MINUS, lexer.currentChar)

		}

	case '/':
		switch lexer.peekCharacter() {
		case '=':
			character := lexer.currentChar
			lexer.consumeChar()

			t = token.Token{Type: token.SLASHEQ,
				Literal: string(character) + string(lexer.currentChar)}
		default:
			t = newToken(token.SLASH, lexer.currentChar)
		}
	case '*':
		switch lexer.peekCharacter() {
		case '=':
			character := lexer.currentChar
			lexer.consumeChar()

			t = token.Token{Type: token.MULEQ,
				Literal: string(character) + string(lexer.currentChar)}
		default:
			t = newToken(token.ASTERISK, lexer.currentChar)
		}
	case '<':
		if lexer.peekCharacter() == '=' {
			character := lexer.currentChar
			lexer.consumeChar()

			t = token.Token{
				Type:    token.LTEQ,
				Literal: string(character) + string(lexer.currentChar),
			}
		} else {
			t = newToken(token.LT, lexer.currentChar)
		}
	case '>':
		if lexer.peekCharacter() == '=' {
			character := lexer.currentChar
			lexer.consumeChar()

			t = token.Token{
				Type:    token.GTEQ,
				Literal: string(character) + string(lexer.currentChar),
			}
		} else {
			t = newToken(token.GTEQ, lexer.currentChar)
		}
	case ';':
		t = newToken(token.SEMICOLON, lexer.currentChar)
	case ':':
		t = newToken(token.COLON, lexer.currentChar)

	case '#':
		for lexer.peekCharacter() != '\n' && lexer.peekCharacter() != 0 {
			lexer.consumeChar()
		}
		lexer.skipWhitespace()
		lexer.consumeChar()
		return lexer.NextToken()

	// EOF
	case 0:
		t.Literal = ""
		t.Type = token.EOF

	// Multiple-character bytes
	case '=':
		if lexer.peekCharacter() == '=' {
			character := lexer.currentChar
			lexer.consumeChar()

			t = token.Token{
				Type:    token.EQ,
				Literal: string(character) + string(lexer.currentChar),
			}
		} else {
			t = newToken(token.ASSIGN, lexer.currentChar)
		}
	case '!':
		if lexer.peekCharacter() == '=' {
			character := lexer.currentChar

			lexer.consumeChar()
			t = token.Token{
				Type:    token.NOT_EQ,
				Literal: string(character) + string(lexer.currentChar),
			}
		} else {
			t = newToken(token.BANG, lexer.currentChar)
		}

	case '|':
		switch lexer.peekCharacter() {
		case '|':
			character := lexer.currentChar
			lexer.consumeChar()

			t = token.Token{
				Type:    token.OR,
				Literal: string(character) + string(lexer.currentChar),
			}
		default:
			t = newToken(token.ILLEGAL, lexer.currentChar)

		}

	case '&':
		switch lexer.peekCharacter() {
		case '&':
			character := lexer.currentChar
			lexer.consumeChar()

			t = token.Token{
				Type:    token.AND,
				Literal: string(character) + string(lexer.currentChar),
			}
		default:
			t = newToken(token.ILLEGAL, lexer.currentChar)

		}

	case '"':
		t.Type = token.STRING
		t.Literal = lexer.readString()

	default:
		if util.IsLetter(lexer.currentChar) {
			t.Literal = lexer.consumeIdentifier()
			t.Type = token.LookupIdent(t.Literal)
			return t
		} else if util.IsDigit(lexer.currentChar) {
			t.Literal = lexer.consumeInteger()
			t.Type = token.INT
			return t
		} else {
			t = newToken(token.ILLEGAL, lexer.currentChar)
		}
	}
	lexer.consumeChar()

	return t
}

func (lexer *Lexer) skipWhitespace() {
	for util.IsWhitespace(lexer.currentChar) {
		lexer.consumeChar()
	}
}

func (lexer *Lexer) consumeInteger() string {
	position := lexer.position
	for util.IsDigit(lexer.currentChar) {
		lexer.consumeChar()
	}

	return lexer.input[position:lexer.position]
}

func (lexer *Lexer) consumeIdentifier() string {
	position := lexer.position
	for util.IsLetter(lexer.currentChar) {
		lexer.consumeChar()
	}

	return lexer.input[position:lexer.position]
}

func (lexer *Lexer) peekCharacter() byte {
	if lexer.readingPosition >= len(lexer.input) {
		return 0 // this will trigger an EOF
	} else {
		return lexer.input[lexer.readingPosition] // Recall that readingPosition is always one character ahead of the pos
	}
}

func newToken(tokenType token.Type, character byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(character)}
}

func (lexer *Lexer) readString() string {
	position := lexer.position + 1
	for {
		lexer.consumeChar()

		if lexer.currentChar == '"' || lexer.currentChar == 0 {
			break
		}
	}

	return lexer.input[position:lexer.position]
}

// New creates a new Lexer and returns a reference to it.
func New(input string) *Lexer {
	lexer := Lexer{input: input}
	lexer.consumeChar()
	return &lexer
}
