package lexer

import "fmt"

type Token struct {
	Type  TokenType
	Bytes []byte
}

func (t Token) String() string {
	return fmt.Sprintf("%s: %s", t.Type.String(), string(t.Bytes))
}

type TokenType int

func (t TokenType) Token() (Token, error) {
	switch t {
	case TLeftAngleBracket:
		return Token{Type: t, Bytes: []byte("<")}, nil
	case TRightAngleBracket:
		return Token{Type: t, Bytes: []byte(">")}, nil
	case TLeftSquareBracket:
		return Token{Type: t, Bytes: []byte("[")}, nil
	case TRightSquareBracket:
		return Token{Type: t, Bytes: []byte("]")}, nil
	case TLeftCurlyBracket:
		return Token{Type: t, Bytes: []byte("{")}, nil
	case TRightCurlyBracket:
		return Token{Type: t, Bytes: []byte("}")}, nil
	case TLeftBracket:
		return Token{Type: t, Bytes: []byte("(")}, nil
	case TRightBracket:
		return Token{Type: t, Bytes: []byte(")")}, nil
	case TComma:
		return Token{Type: t, Bytes: []byte(",")}, nil
	default:
		return Token{}, fmt.Errorf("unexpected token type %q", t.String())
	}
}

func (t TokenType) String() string {
	switch t {
	case TString:
		return "string"
	case TNumber:
		return "0-9"
	case TLeftAngleBracket:
		return "<"
	case TRightAngleBracket:
		return ">"
	case TLeftSquareBracket:
		return "["
	case TRightSquareBracket:
		return "]"
	case TLeftCurlyBracket:
		return "{"
	case TRightCurlyBracket:
		return "}"
	case TLeftBracket:
		return "("
	case TRightBracket:
		return ")"
	case TComma:
		return ","
	default:
		return "UNKNOWN"
	}
}

const (
	TString TokenType = iota
	TNumber
	TLeftAngleBracket
	TRightAngleBracket
	TLeftSquareBracket
	TRightSquareBracket
	TLeftCurlyBracket
	TRightCurlyBracket
	TLeftBracket
	TRightBracket
	TComma
)
