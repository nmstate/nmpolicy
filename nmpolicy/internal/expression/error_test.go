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

package expression_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nmstate/nmpolicy/nmpolicy/internal/expression"
)

func TestError(t *testing.T) {
	tests := []struct {
		err         *expression.Error
		expectedFmt string
	}{
		{expression.WrapError(fmt.Errorf("test error"), 33), "test error, pos=33"},
		{expression.WrapError(fmt.Errorf("test error"), 4).Decorate("0123456"), `test error, pos=4
| 0123456
| ....^`},
	}

	for _, tt := range tests {
		t.Run(tt.expectedFmt, func(t *testing.T) {
			assert.Equal(t, tt.expectedFmt, tt.err.Error())
		})
	}
}
