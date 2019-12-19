package lexer

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"text/scanner"

	"github.com/SimonRichardson/juju-schema-gen/pkg/cursor"
	"github.com/SimonRichardson/juju-schema-gen/pkg/errors"
)

// Lexer iterates over a source
type Lexer struct {
	types  map[string]TokenType
	tokens []*Token
	ptr    int
}

// New creates a new lexer from some token types
func New(types map[string]TokenType) *Lexer {
	return &Lexer{
		types:  types,
		tokens: make([]*Token, 0),
		ptr:    0,
	}
}

func (l *Lexer) Tokens() []*Token {
	return l.tokens[:]
}

func (l *Lexer) Write(p []byte) (index int, err error) {
	var line string

	var s scanner.Scanner
	s.Init(bytes.NewBuffer(p))
	s.Filename = "service.api"
	s.Error = func(s *scanner.Scanner, msg string) {
		text := s.TokenText()
		err = errors.CharPositionError{
			Context: fmt.Sprintf("%s%s", line, text),
			Char:    text,
			Position: cursor.Position{
				Line:  s.Position.Line,
				Start: (s.Position.Column - 1),
				End:   (s.Position.Column - 1) + len(text),
			},
		}
	}

	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		if err != nil {
			return
		}

		text := s.TokenText()
		line += fmt.Sprintf("%s ", text)
		// Check to see if the text is a known type
		if tokenType, ok := l.types[text]; ok {
			l.tokens = append(l.tokens, &Token{
				Type:  tokenType,
				Bytes: []byte(text),
				Position: cursor.Position{
					Line:  s.Position.Line,
					Start: (s.Position.Column - 1),
					End:   (s.Position.Column - 1) + len(text),
				},
			})
			goto SKIP
		}

		if _, err := strconv.Atoi(text); err == nil {
			l.tokens = append(l.tokens, &Token{
				Type:  TNumber,
				Bytes: []byte(text),
				Position: cursor.Position{
					Line:  s.Position.Line,
					Start: (s.Position.Column - 1),
					End:   (s.Position.Column - 1) + len(text),
				},
			})
			goto SKIP
		}

		l.tokens = append(l.tokens, &Token{
			Type:  TString,
			Bytes: []byte(text),
			Position: cursor.Position{
				Line:  s.Position.Line,
				Start: (s.Position.Column - 1),
				End:   (s.Position.Column - 1) + len(text),
			},
		})

	SKIP:
		if s.Peek() == '\n' {
			line = ""
		}
	}
	return len(p), nil
}

func (l *Lexer) String() string {
	var buf []string
	for _, v := range l.tokens {
		buf = append(buf, fmt.Sprintf("%d %q", v.Type, string(v.Bytes)))
	}
	return strings.Join(buf, "\n")
}
