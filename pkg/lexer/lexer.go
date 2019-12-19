package lexer

import (
	"fmt"
	"io"
	"strings"

	"github.com/SimonRichardson/juju-schema-gen/pkg/errors"
)

// Lexer iterates over a source
type Lexer struct {
	types  map[byte]TokenType
	tokens []*Token
	ptr    int
}

// New creates a new lexer from some token types
func New(types map[byte]TokenType) *Lexer {
	return &Lexer{
		types:  types,
		tokens: make([]*Token, 0),
		ptr:    0,
	}
}

func (l *Lexer) Tokens() []*Token {
	return l.tokens[:]
}

func (l *Lexer) Write(p []byte) (int, error) {
	var index int
	var b byte
	var line int
	var currentLine []byte
LOOP:
	for index, b = range p {
		currentLine = append(currentLine, b)
		switch {
		case b == ' ' || b == '\n':
			if b == '\n' {
				line++
				currentLine = make([]byte, 0)
			}
			l.ptr = len(l.tokens)
			continue LOOP
		case (b >= '0' && b <= '9'):
			if l.ptr == len(l.tokens) {
				l.tokens = append(l.tokens, &Token{
					Type:  TNumber,
					Bytes: []byte{b},
				})
				l.ptr = len(l.tokens) + 1
			} else {
				token := l.tokens[len(l.tokens)-1]
				token.Bytes = append(token.Bytes, b)
			}
		case (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z'):
			if l.ptr == len(l.tokens) {
				l.tokens = append(l.tokens, &Token{
					Type:  TString,
					Bytes: []byte{b},
				})
				l.ptr = len(l.tokens) + 1
			} else {
				token := l.tokens[len(l.tokens)-1]
				token.Bytes = append(token.Bytes, b)
			}
		default:
			if tokenType, ok := l.types[b]; ok {
				l.tokens = append(l.tokens, &Token{
					Type:  tokenType,
					Bytes: []byte{b},
				})
				l.ptr = len(l.tokens)
				continue LOOP
			}
			break LOOP
		}
	}
	if index == len(p)-1 {
		return index, io.EOF
	}
	return index, errors.CharPositionError{
		Context:       string(currentLine),
		Char:          string(p[index]),
		Line:          line + 1,
		PositionStart: len(currentLine) - 1,
		PositionEnd:   len(currentLine),
	}
}

func (l *Lexer) String() string {
	var buf []string
	for _, v := range l.tokens {
		buf = append(buf, fmt.Sprintf("%d %q", v.Type, string(v.Bytes)))
	}
	return strings.Join(buf, "\n")
}
