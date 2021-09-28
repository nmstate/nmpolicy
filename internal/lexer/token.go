package lexer

type TokenType int

const (
	EOF TokenType = iota
	IDENTITY
	NUMBER
	STRING

	PATH     // .
	PIPE     // |
	LBRACKET // [
	RBRACKET // ]

	beginOperatorSection
	ASSIGN   // =
	EQUALITY // ==
	MERGE    // +
	endOperatorSection
)

var tokens = []string{
	EOF:      "EOF",
	IDENTITY: "IDENTITY",
	NUMBER:   "NUMBER",
	STRING:   "STRING",

	PATH:     "PATH",
	PIPE:     "PIPE",
	LBRACKET: "LBRACKET",
	RBRACKET: "RBRACKET",

	ASSIGN:   "ASSIGN",
	EQUALITY: "EQUALITY",
	MERGE:    "MERGE",
}

func (t TokenType) String() string {
	return tokens[t]
}

type Token struct {
	Position int
	Type     TokenType
	Literal  string
}

func (t TokenType) IsOperator() bool {
	return t > beginOperatorSection && t < endOperatorSection
}
