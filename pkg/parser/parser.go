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

func (p *Parser) Expressions() []Expression {
	return p.expressions
}

func (p *Parser) String() string {
	var buf []string
	for _, v := range p.expressions {
		expressions := describeExpressions(v, "  ")
		buf = append(buf, fmt.Sprintf("%s %v%s", v.Type(), v.Tokens(), expressions))
	}
	return strings.Join(buf, "\n")
}

func describeExpressions(v Expression, s string) string {
	var result string
	if ex, ok := v.(RecursiveExpression); ok {
		result += "\n"
		for _, v := range ex.Expressions() {
			result += fmt.Sprintf("%s- %s %v\n", s, v.Type(), v.Tokens())
			result += describeExpressions(v, s+"  ")
		}
	}
	return result
}

type reader struct {
	types  map[lexer.TokenType]Parselet
	tokens []*lexer.Token
	offset int
}

func (r *reader) Peek(i int) (lexer.Token, error) {
	if i >= len(r.tokens) {
		return lexer.Token{}, fmt.Errorf("unable to peek %d:%d", i, len(r.tokens))
	}
	return *r.tokens[i], nil
}

func (r *reader) AdvanceTo(i int) {
	r.offset = i
}

func (r *reader) Parse() ([]Expression, int, error) {
	expressions := make([]Expression, 0)
	i := r.offset
	for ; i < len(r.tokens); i++ {
		token := r.tokens[i]
		if parselet, ok := r.types[token.Type]; ok {
			sub := &reader{
				types:  r.types,
				tokens: r.tokens[i+1:],
			}
			expression, err := parselet.Parse(sub, *token)
			if err != nil {
				return nil, -1, err
			}
			expressions = append(expressions, expression)
			i += sub.offset
		} else {
			break
		}
	}
	r.offset += i + 1
	return expressions, i - 1, nil
}

func (r *reader) Len() int {
	return len(r.tokens)
}
