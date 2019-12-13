package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/SimonRichardson/juju-schema-gen/pkg/lexer"
	"github.com/SimonRichardson/juju-schema-gen/pkg/parser"
)

func main() {
	bytes, err := ioutil.ReadFile("examples/service.api")
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	lex := lexer.New(map[byte]lexer.TokenType{
		'[': lexer.TLeftSquareBracket,
		']': lexer.TRightSquareBracket,
		'{': lexer.TLeftCurlyBracket,
		'}': lexer.TRightCurlyBracket,
		'(': lexer.TLeftBracket,
		')': lexer.TRightBracket,
		',': lexer.TComma,
	})
	_, err = lex.Write(bytes)
	fmt.Println(err)
	fmt.Println("----")

	par := parser.New(map[lexer.TokenType]parser.Parselet{
		lexer.TString:            &parser.Keyword{},
		lexer.TLeftSquareBracket: &parser.Version{},
	})
	_, err = par.Read(lex)
	fmt.Println(err)
	fmt.Println(par.String())
}
