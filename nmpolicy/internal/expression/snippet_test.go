/*
 * Copyright 2001 NMPolicy Authors.
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

package expression

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSource(t *testing.T) {
	tests := []struct {
		expression string
		pos        int
		snippet    string
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
		t.Run(fmt.Sprintf("snippet of '%s' at '%d'", tt.expression, tt.pos), func(t *testing.T) {
			snippet := "\n" + snippet(tt.expression, tt.pos)
			assert.Equal(t, tt.snippet, snippet)
		})
	}
}
