package ast_test

import (
	"testing"

	assert "github.com/stretchr/testify/require"

	"github.com/nmstate/nmpolicy/nmpolicy/internal/ast"
)

func TestTeminalDeepEqual(t *testing.T) {
	literalA := "literalA"
	literalB := "literalB"

	newPtr := func(literal string) *string {
		return &literal
	}

	tests := []struct {
		lhs      ast.Terminal
		rhs      ast.Terminal
		expected bool
	}{
		{ast.Terminal{Identity: &literalA}, ast.Terminal{Identity: &literalA}, true},
		{ast.Terminal{String: &literalA}, ast.Terminal{String: &literalA}, true},
		{
			ast.Terminal{String: &literalA, Identity: &literalB},
			ast.Terminal{String: &literalA, Identity: &literalB},
			true,
		},
		{
			ast.Terminal{String: nil},
			ast.Terminal{String: nil},
			true,
		},
		{
			ast.Terminal{String: nil},
			ast.Terminal{String: &literalA},
			false,
		},
		{
			ast.Terminal{String: &literalA},
			ast.Terminal{String: &literalB},
			false,
		},
		{
			ast.Terminal{String: &literalA},
			ast.Terminal{String: newPtr(literalA)},
			true,
		},
		{
			ast.Terminal{String: &literalA},
			ast.Terminal{String: newPtr(literalB)},
			false,
		},
		{
			ast.Terminal{Identity: nil},
			ast.Terminal{Identity: &literalA},
			false,
		},
		{
			ast.Terminal{Identity: &literalA},
			ast.Terminal{Identity: &literalB},
			false,
		},
		{
			ast.Terminal{Identity: &literalA},
			ast.Terminal{Identity: newPtr(literalA)},
			true,
		},
		{
			ast.Terminal{Identity: &literalA},
			ast.Terminal{Identity: newPtr(literalB)},
			false,
		},
		{
			ast.Terminal{Identity: &literalA, String: &literalA},
			ast.Terminal{Identity: &literalB, String: newPtr(literalA)},
			false,
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, test.lhs.DeepEqual(test.rhs))
	}
}
