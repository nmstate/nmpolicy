package lexer

type TokenType int

const (
	EOF TokenType = iota
	IDENTITY
	NUMBER
	STRING

	DOT  // .
	PIPE // |

	REPLACE  // :=
	EQFILTER // ==
	MERGE    // +

	TRUE  // true
	FALSE // false
)

var tokens = []string{
	EOF:      "EOF",
	IDENTITY: "IDENTITY",
	NUMBER:   "NUMBER",
	STRING:   "STRING",

	DOT:  "DOT",
	PIPE: "PIPE",

	REPLACE:  "REPLACE",
	EQFILTER: "EQFILTER",
	MERGE:    "MERGE",

	TRUE:  "TRUE",
	FALSE: "FALSE",
}

func (t TokenType) String() string {
	return tokens[t]
}

type Token struct {
	Position int
	Type     TokenType
	Literal  string
}
