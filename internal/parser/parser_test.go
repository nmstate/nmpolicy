package parser

import (
	"fmt"
	"testing"

	"github.com/nmstate/nmpolicy/internal/ast"
	"github.com/nmstate/nmpolicy/internal/lexer"
	"github.com/nmstate/nmpolicy/internal/source"
	"github.com/stretchr/testify/assert"
)

type lexStubReturn struct {
	token *lexer.Token
	err   string
}

type whenLexReturns []lexStubReturn

type lexerStub struct {
	idx     int
	returns whenLexReturns
	t       *testing.T
}

func (l *lexerStub) Lex() (*lexer.Token, error) {
	if l.idx >= len(l.returns) {
		return &lexer.Token{Position: len(l.returns) - 1, Type: lexer.EOF, Literal: ""}, nil
	}
	r := l.returns[l.idx]
	if r.token != nil {
		r.token.Position = l.idx
	}

	l.idx++
	var err error
	if r.err != "" {
		err = fmt.Errorf(r.err)
	}
	return r.token, err
}

func (l lexStubReturn) String() string {
	if l.token != nil {
		return fmt.Sprintf("{%s,%s}", l.token.Type, l.token.Literal)
	}
	return fmt.Sprintf("{err: %s}", l.err)
}

func TestParser(t *testing.T) {

	token := func(t lexer.TokenType, l string) lexStubReturn {
		return lexStubReturn{&lexer.Token{Type: t, Literal: l}, ""}
	}

	id := func(literal string) lexStubReturn {
		return token(lexer.IDENTITY, literal)
	}

	num := func(literal string) lexStubReturn {
		return token(lexer.NUMBER, literal)
	}

	str := func(literal string) lexStubReturn {
		return token(lexer.STRING, literal)
	}

	symbol := func(t lexer.TokenType) lexStubReturn {
		return token(t, t.String())
	}

	path := func() lexStubReturn {
		return symbol(lexer.PATH)
	}

	eq := func() lexStubReturn {
		return symbol(lexer.EQUALITY)
	}

	merge := func() lexStubReturn {
		return symbol(lexer.MERGE)
	}

	assign := func() lexStubReturn {
		return symbol(lexer.ASSIGN)
	}

	eof := func() lexStubReturn {
		return symbol(lexer.EOF)
	}

	pipe := func() lexStubReturn {
		return symbol(lexer.PIPE)
	}

	type parserShouldReturn struct {
		ast string
		err string
	}

	var tests = []struct {
		whenLexReturns     whenLexReturns
		parserShouldReturn parserShouldReturn
	}{
		{whenLexReturns{
			{err: "forced failure at test"},
		}, parserShouldReturn{
			ast: "",
			err: "failed lexing next token: forced failure at test",
		}},

		{whenLexReturns{
			eof(),
		}, parserShouldReturn{ast: "", err: `missing command, pos=0
| 0123456789
| ^`}},

		{whenLexReturns{
			id("foo"),
			id("bar"),
			eof(),
		}, parserShouldReturn{ast: "", err: `bad path format: identities has to be separated by a dot, pos=1
| 0123456789
| .^`}},

		{whenLexReturns{
			id("foo"),
			path(),
			eof(),
		}, parserShouldReturn{ast: "", err: `bad path format: dot has to be followed by indentity or index, pos=2
| 0123456789
| ..^`}},

		{whenLexReturns{
			id("foo"),
			path(),
			str("bar"),
			eof(),
		}, parserShouldReturn{ast: "", err: `bad path format: string can only be used with brackets, ej foo["bar"], pos=2
| 0123456789
| ..^`}},

		{whenLexReturns{
			id("foo"),
			path(),
			num("3"),
			path(),
			id("dar"),
			eof(),
		}, parserShouldReturn{
			ast: `
pos: 0
path:
  - pos: 0
    id: foo
  - pos: 2
    idx: 3
  - pos: 4
    id: dar
`,
		}},

		{whenLexReturns{
			id("foo"),
			path(),
			num("3"),
			path(),
			id("dar"),
			eq(),
			str("foo"),
		}, parserShouldReturn{
			ast: `
pos: 0
equal:
- pos: 0
  path:
  - pos: 0
    id: foo
  - pos: 2
    idx: 3
  - pos: 4
    id: dar
- pos: 6
  string: foo
`,
		}},

		{whenLexReturns{
			id("foo"),
			path(),
			num("3"),
			path(),
			id("dar"),
			eq(),
			num("3424"),
		}, parserShouldReturn{
			ast: `
pos: 0
equal:
- pos: 0
  path:
  - pos: 0
    id: foo
  - pos: 2
    idx: 3
  - pos: 4
    id: dar
- pos: 6
  number: 3424
`,
		}},

		{whenLexReturns{
			id("foo"),
			path(),
			num("3"),
			path(),
			id("dar"),
			eq(),
			id("moo"),
			path(),
			id("boo"),
		}, parserShouldReturn{
			ast: `
pos: 0
equal:
- pos: 0
  path:
  - pos: 0
    id: foo
  - pos: 2
    idx: 3
  - pos: 4
    id: dar
- pos: 6
  path:
  - pos: 6
    id: moo
  - pos: 8
    id: boo
`,
		}},

		{whenLexReturns{
			id("foo"),
			path(),
			num("3"),
			path(),
			id("dar"),
			eq(),
		}, parserShouldReturn{ast: "", err: `bad EQUALITY operation format: supported argument missing, pos=5
| 0123456789
| .....^`}},
		{whenLexReturns{
			id("foo"),
			path(),
			num("3"),
			path(),
			id("dar"),
			assign(),
			str("foo"),
		}, parserShouldReturn{
			ast: `
pos: 0
assign:
- pos: 0
  path:
  - pos: 0
    id: foo
  - pos: 2
    idx: 3
  - pos: 4
    id: dar
- pos: 6
  string: foo
`,
		}},

		{whenLexReturns{
			id("foo"),
			path(),
			num("3"),
			path(),
			id("dar"),
			assign(),
			num("3424"),
		}, parserShouldReturn{
			ast: `
pos: 0
assign:
- pos: 0
  path:
  - pos: 0
    id: foo
  - pos: 2
    idx: 3
  - pos: 4
    id: dar
- pos: 6
  number: 3424
`,
		}},

		{whenLexReturns{
			id("foo"),
			path(),
			num("3"),
			path(),
			id("dar"),
			assign(),
			id("moo"),
			path(),
			id("boo"),
		}, parserShouldReturn{
			ast: `
pos: 0
assign:
- pos: 0
  path:
  - pos: 0
    id: foo
  - pos: 2
    idx: 3
  - pos: 4
    id: dar
- pos: 6
  path:
  - pos: 6
    id: moo
  - pos: 8
    id: boo
`,
		}},

		{whenLexReturns{
			id("foo"),
			path(),
			num("3"),
			path(),
			id("dar"),
			assign(),
			eof(),
		}, parserShouldReturn{ast: "", err: `bad ASSIGN operation format: supported argument missing, pos=6
| 0123456789
| ......^`}},

		{whenLexReturns{
			id("foo"),
			path(),
			num("3"),
			path(),
			id("dar"),
			merge(),
			id("moo"),
			path(),
			id("boo"),
		}, parserShouldReturn{
			ast: `
pos: 0
merge:
- pos: 0
  path:
  - pos: 0
    id: foo
  - pos: 2
    idx: 3
  - pos: 4
    id: dar
- pos: 6
  path:
  - pos: 6 
    id: moo
  - pos: 8
    id: boo
`,
		}},

		{whenLexReturns{
			id("foo"),
			path(),
			num("3"),
			path(),
			id("dar"),
			merge(),
		}, parserShouldReturn{ast: "", err: `bad MERGE operation format: supported argument missing, pos=5
| 0123456789
| .....^`}},

		{whenLexReturns{
			id("foo"),
			path(),
			num("3"),
			path(),
			id("dar"),
			merge(),
			str("foo"),
		}, parserShouldReturn{ast: "", err: `bad MERGE operation format: only paths can be merged, pos=6
| 0123456789
| ......^`}},

		{whenLexReturns{
			id("foo"),
			path(),
			num("3"),
			path(),
			id("dar"),
			merge(),
			num("3424"),
		}, parserShouldReturn{ast: "", err: `bad MERGE operation format: only paths can be merged, pos=6
| 0123456789
| ......^`}},
		{whenLexReturns{
			id("foo1"),
			path(),
			id("bar1"),
			pipe(),
			id("foo2"),
			path(),
			id("bar2"),
		}, parserShouldReturn{
			ast: `
pos: 0
path:
- pos: 0
  id: foo1
- pos: 2
  id: bar1
pipe:
  pos: 4
  path: 
  - pos: 4
    id: foo2
  - pos: 6
    id: bar2
`,
		}},
		{whenLexReturns{
			id("foo"),
			path(),
			num("3"),
			path(),
			id("dar"),
			pipe(),
		}, parserShouldReturn{ast: "", err: `missing command, pos=5
| 0123456789
| .....^`}},

		{whenLexReturns{
			pipe(),
			id("foo"),
			path(),
			num("3"),
			path(),
			id("dar"),
		}, parserShouldReturn{ast: "", err: `bad PIPE operation format: missing cmd to pipe from, pos=0
| 0123456789
| ^`}},
	}

	for ti, tt := range tests {
		description := fmt.Sprintf("\nwhenLexReturns: %+v\nparserShouldReturn:\n%+v\n", tt.whenLexReturns, tt.parserShouldReturn)
		t.Run(fmt.Sprintf("%d", ti+1), func(t *testing.T) {
			t.Log(description)
			source := source.New("0123456789")
			parser := New(*source, &lexerStub{t: t, returns: tt.whenLexReturns})
			obtainedAST, obtainedErr := parser.Parse()

			if tt.parserShouldReturn.err != "" {
				assert.EqualError(t, obtainedErr, tt.parserShouldReturn.err)
			} else {
				assert.NoError(t, obtainedErr)
				parserShouldReturnAST, err := ast.FromYAMLString(tt.parserShouldReturn.ast)
				assert.NoError(t, err)
				assert.Equal(t, parserShouldReturnAST, obtainedAST)
			}
		})
	}
}
