package lexer

type Token struct {
	Type  TokenType
	Bytes []byte
}

type TokenType int

const (
	TString TokenType = iota
	TNumber
	TRightAngleBracket
	TLeftAngleBracket
	TRightSquareBracket
	TLeftSquareBracket
	TRightCurlyBracket
	TLeftCurlyBracket
	TRightBracket
	TLeftBracket
	TComma
)
