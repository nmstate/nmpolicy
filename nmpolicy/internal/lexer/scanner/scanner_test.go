/*
 * Copyright 2021 NMPolicy Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at:
 *
 *	  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package scanner_test

import (
	"strings"
	"testing"

	assert "github.com/stretchr/testify/require"

	"github.com/nmstate/nmpolicy/nmpolicy/internal/lexer/scanner"
)

func TestReader(t *testing.T) {
	var tests = []struct {
		str string
	}{
		{""},
		{"foo bar dar"},
		{"    "},
	}
	for _, tt := range tests {
		t.Run(tt.str, func(t *testing.T) {
			scn := scanner.New(strings.NewReader(tt.str))
			runes := []rune(tt.str)
			for {
				err := scn.Next()
				if scn.Rune() == scanner.EOF {
					if len(runes) > 0 {
						assert.Equal(t, len(runes)-1, scn.Position())
					}
					err = scn.Next()
					assert.NoError(t, err)
					assert.Equal(t, scanner.EOF, scn.Rune())
					break
				}
				if len(runes) > 0 {
					assert.NoError(t, err)

					t.Logf("Calling Prev go back to previous rune and position, r: \"%s\", p: %d", string(scn.Rune()), scn.Position())
					p := scn.Position()
					err = scn.Prev()
					assert.NoError(t, err)
					if p > 0 {
						assert.Equal(t, p-1, scn.Position())
						assert.Equal(t, string(runes[p-1]), string(scn.Rune()))
					}

					t.Log("Calling Prev twice fail")
					err = scn.Prev()
					assert.Error(t, err)

					t.Log("Going back to next")
					err = scn.Next()
					assert.NoError(t, err)
					assert.Equal(t, string(scn.Rune()), string(runes[scn.Position()]))
				}
			}
		})
	}
}
