package ast

import (
	"fmt"
	"io"

	"github.com/SimonRichardson/juju-schema-gen/pkg/cursor"
	"github.com/SimonRichardson/juju-schema-gen/pkg/parser"
)

type Facade struct {
	cursor.Position
	Name    parser.Expression
	Version parser.Expression
	Params  []Data
	Results []Data
	Methods []Methods
}

func (p *Facade) Read(expressions []parser.Expression) (int, error) {
	exp, err := matchString(expressions, 0, "facade")
	if err != nil {
		return -1, err
	}
	p.Position = exp.Position()

	exp, _, err = pluckString(expressions, 1)
	if err != nil {
		return -1, err
	}
	p.Name = exp

	exp, _, err = pluckVersion(expressions, 2)
	if err != nil {
		return -1, err
	}
	p.Version = exp

	return len(expressions) + 1, nil
}

func (p *Facade) Format(out io.Writer) error {
	fmt.Fprintf(out, "facade %s%v\n", tokenStrings(p.Name.Tokens()), tokenStrings(p.Version.Tokens()))
	return nil
}
