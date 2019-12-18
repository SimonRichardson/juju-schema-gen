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
				types:  p.types,
				tokens: tokens[i+1:],
			}
			expression, err := parselet.Parse(r, *token)
			if err != nil {
				return i, err
			}
			p.expressions = append(p.expressions, expression)
			i += r.offset
		} else {
			return i, fmt.Errorf("unexpected token type %q at %d", token.Type.String(), i)
		}
	}
	return len(tokens), nil
}

func (p *Parser) String() string {
	var buf []string
	for _, v := range p.expressions {
		var expressions string
		if ex, ok := v.(RecursiveExpression); ok {
			expressions = fmt.Sprintf(" %s", ex.Expressions())
		}
		buf = append(buf, fmt.Sprintf("%d %v%s", v.Type(), v.Tokens(), expressions))
	}
	return strings.Join(buf, "\n")
}

type reader struct {
	types  map[lexer.TokenType]Parselet
	tokens []*lexer.Token
	offset int
}

func (r *reader) Peek(i int) (lexer.Token, error) {
	if i >= len(r.tokens) {
		return lexer.Token{}, fmt.Errorf("not found")
	}
	return *r.tokens[i], nil
}

func (r *reader) AdvanceTo(i int) {
	r.offset = i
}

func (r *reader) Parse() ([]Expression, error) {
	expressions := make([]Expression, 0)
	for i := r.offset; i < len(r.tokens); i++ {
		token := r.tokens[i]
		if parselet, ok := r.types[token.Type]; ok {
			r := &reader{
				types:  r.types,
				tokens: r.tokens[i+1:],
			}
			expression, err := parselet.Parse(r, *token)
			if err != nil {
				return nil, err
			}
			expressions = append(expressions, expression)
			i += r.offset
		} else {
			return nil, fmt.Errorf("unexpected token type %q at %d", token.Type.String(), i)
		}
	}
	return expressions, nil
}
