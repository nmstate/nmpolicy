package lexer

type Lexer interface {
	Lex() (*Token, error)
}
