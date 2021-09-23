package parser

import (
	"fmt"

	"github.com/nmstate/nmpolicy/internal/lexer"
)

func badPathFormatError(message string) error {
	return fmt.Errorf("bad path format: %s", message)
}

func badOperationFormat(opType lexer.TokenType, message string) error {
	return fmt.Errorf("bad %s operation format: %s", opType, message)
}
