package source

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSource(t *testing.T) {

	tests := []struct {
		str     string
		pos     int
		snippet string
	}{
		{"012345678", 5, `
| 012345678
| .....^`},
		{"012345678", -3, `
| 012345678
| ^`},
		{"012345678", 10, `
| 012345678
| ........^`},
		{"", 10, `
`},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("snippet of '%s' at '%d'", tt.str, tt.pos), func(t *testing.T) {
			snippet := "\n" + New(tt.str).Snippet(tt.pos)
			assert.Equal(t, tt.snippet, snippet)
		})
	}
}
