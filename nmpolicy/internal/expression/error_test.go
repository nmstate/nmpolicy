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
	err := fmt.Errorf("test error")
	wrappedError := expression.WrapError(err, "0123456", 4)
	decoratedError := expression.DecorateError(err, "0123456", 4)
	expectedErrorMsg := `test error
| 0123456
| ....^`
	assert.EqualError(t, wrappedError, expectedErrorMsg)
	assert.ErrorIs(t, wrappedError, err)
	assert.EqualError(t, decoratedError, expectedErrorMsg)
	assert.NotErrorIs(t, decoratedError, err)
}
