package parser

import (
	"fmt"
	"strings"

	"github.com/SimonRichardson/juju-schema-gen/pkg/lexer"
)

type Parser struct {
	types       map[lexer.TokenType]Parselet
	expressions []Expression
}

func New(types map[lexer.TokenType]Parselet) *Parser {
	return &Parser{
		types: types,
	}
}

func (p *Parser) Read(lex *lexer.Lexer) (int, error) {
	tokens := lex.Tokens()
	for i := 0; i < len(tokens); i++ {
		token := tokens[i]
		if parselet, ok := p.types[token.Type]; ok {
			r := &reader{
				tokens: tokens[i+1:],
			}
			expression, err := parselet.Parse(r, *token)
			if err != nil {
				return i, err
			}
			p.expressions = append(p.expressions, expression)
			i += r.offset
		} else {
			return i, fmt.Errorf("unexpected token type %d at %d", token.Type, i)
		}
	}
	return len(tokens), nil
}

func (p *Parser) String() string {
	var buf []string
	for _, v := range p.expressions {
		buf = append(buf, fmt.Sprintf("%d %v", v.Type(), v.Tokens()))
	}
	return strings.Join(buf, "\n")
}

type reader struct {
	tokens []*lexer.Token
	offset int
}

func (r *reader) Peek(i int) (lexer.Token, error) {
	if i >= len(r.tokens) {
		return lexer.Token{}, fmt.Errorf("not found")
	}
	return *r.tokens[i], nil
}

func (r *reader) Advance(i int) {
	r.offset = i
}
