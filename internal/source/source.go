package source

import (
	"fmt"
	"strings"
)

// Source represents the raw data for a nmpolicy expression
type Source struct {
	src    string
	reader *strings.Reader
}

// New construct Source from a string
func New(src string) *Source {
	return &Source{
		src:    src,
		reader: strings.NewReader(src),
	}
}

// Reader reaturns a strings.Reader to iterate src
func (s Source) Reader() *strings.Reader {
	return strings.NewReader(s.src)
}

// Snippet returns a string containg src and a pointer at pos.
// Example of str "123456" and pos "4":
//
// | 123456
// | ...^
func (s *Source) Snippet(pos int) string {
	if len(s.src) == 0 {
		return ""
	}

	if pos >= len(s.src) {
		pos = len(s.src) - 1
	}

	marker := strings.Builder{}
	for i := 0; i < pos; i++ {
		marker.WriteString(".")
	}
	marker.WriteString("^")
	return fmt.Sprintf("| %s\n| %s", s.src, marker.String())
}
