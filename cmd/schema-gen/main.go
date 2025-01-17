package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/SimonRichardson/juju-schema-gen/pkg/ast"
	"github.com/SimonRichardson/juju-schema-gen/pkg/lexer"
	"github.com/SimonRichardson/juju-schema-gen/pkg/parser"
)

func main() {
	bytes, err := ioutil.ReadFile("examples/service.api")
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	// Convert a series of bytes into tokens

	lex := lexer.New(map[string]lexer.TokenType{
		"[": lexer.TLeftSquareBracket,
		"]": lexer.TRightSquareBracket,
		"{": lexer.TLeftCurlyBracket,
		"}": lexer.TRightCurlyBracket,
		"(": lexer.TLeftBracket,
		")": lexer.TRightBracket,
		"<": lexer.TLeftAngleBracket,
		">": lexer.TRightAngleBracket,
		",": lexer.TComma,
	})
	_, err = lex.Write(bytes)
	fmt.Println(err)
	fmt.Println("----")

	// Parse the lexer to form a series of expressions

	par := parser.New(map[lexer.TokenType]parser.Parselet{
		lexer.TString: &parser.Keyword{},
		lexer.TLeftAngleBracket: &parser.Version{
			Left:  lexer.TLeftAngleBracket,
			Right: lexer.TRightAngleBracket,
		},
		lexer.TLeftSquareBracket: &parser.Type{
			Left:    lexer.TLeftSquareBracket,
			Right:   lexer.TRightSquareBracket,
			Keyword: lexer.TString,
		},
		lexer.TLeftCurlyBracket: &parser.Body{
			Left:  lexer.TLeftCurlyBracket,
			Right: lexer.TRightCurlyBracket,
		},
		lexer.TLeftBracket: &parser.List{
			Left:      lexer.TLeftBracket,
			Right:     lexer.TRightBracket,
			Separator: lexer.TComma,
		},
	})
	_, err = par.Read(lex)
	fmt.Println(err)
	fmt.Println(par.String())
	fmt.Println("----")

	// Form a AST from the expressions

	tree, err := ast.Generate(par.Expressions())
	fmt.Println(err)
	fmt.Println(tree.String())

	// Interpret the AST
}
