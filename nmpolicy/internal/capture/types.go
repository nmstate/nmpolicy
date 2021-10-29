package capture

import (
	"github.com/nmstate/nmpolicy/nmpolicy/internal/ast"
	"github.com/nmstate/nmpolicy/nmpolicy/internal/lexer"
	"github.com/nmstate/nmpolicy/nmpolicy/types"
)

type CapsExpressions = map[types.CaptureID]types.Expression
type CapsState = map[types.CaptureID]types.CaptureState

type Capture struct {
	astPool         AstPooler
	lexerFactory    LexerFactory
	parserFactory   ParserFactory
	resolverFactory ResolverFactory
}

type AstPooler interface {
	Add(id types.CaptureID, ast ast.Node)
	Range() ast.Pool
}

type LexerFactory = func(types.Expression) Lexer

type Lexer interface {
	Lex() ([]lexer.Token, error)
}

type ParserFactory = func([]lexer.Token) Parser

type Parser interface {
	Parse() (ast.Node, error)
}

type ResolverFactory = func(types.NMState, AstPooler) Resolver

type Resolver interface {
	Resolve() (CapsState, error)
}
