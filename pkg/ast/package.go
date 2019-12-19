package ast

import (
	"fmt"
	"io"

	"github.com/SimonRichardson/juju-schema-gen/pkg/cursor"
	"github.com/SimonRichardson/juju-schema-gen/pkg/parser"
)

type Package struct {
	cursor.Position
	Name    parser.Expression
	Facades []Facade
}

func (p *Package) Read(expressions []parser.Expression) (int, error) {
	exp, err := matchString(expressions, 0, "package")
	if err != nil {
		return -1, err
	}
	p.Position = exp.Position()

	exp, _, err = pluckString(expressions, 1)
	if err != nil {
		return -1, err
	}
	p.Name = exp

	return 2, nil
}

func (p *Package) Format(out io.Writer) error {
	fmt.Fprintf(out, "package %s\n", tokenStrings(p.Name.Tokens()))
	return nil
}
