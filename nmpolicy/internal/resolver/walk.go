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

package resolver

import (
	"fmt"

	"github.com/nmstate/nmpolicy/nmpolicy/internal/ast"
)

func walk(inputState map[string]interface{}, pathSteps ast.VariadicOperator) (interface{}, error) {
	pathVisitorWithWalkOp := pathVisitor{
		stateVisitor: &walkOpVisitor{},
	}

	visitResult, err := pathVisitorWithWalkOp.visitNextStep(path{steps: pathSteps}, inputState)
	if err != nil {
		return nil, fmt.Errorf("failed walking path: %w", err)
	}

	return visitResult, nil
}

type walkOpVisitor struct{}

func (*walkOpVisitor) visitLastMap(mapToAccess map[string]interface{}, accessKey string) (interface{}, error) {
	v, ok := mapToAccess[accessKey]
	if !ok {
		return nil, fmt.Errorf("step not found at map state '%+v'", mapToAccess)
	}
	return v, nil
}

func (*walkOpVisitor) visitLastSlice(sliceToAccess []interface{}, accessIdx int) (interface{}, error) {
	if len(sliceToAccess) <= accessIdx {
		return nil, fmt.Errorf("step not found at slice state '%+v'", sliceToAccess)
	}
	return sliceToAccess[accessIdx], nil
}

func (*walkOpVisitor) visitSliceWithoutIndex(sp stepVisitor, p path, sliceToVisit []interface{}) (interface{}, error) {
	return nil, pathError(p.peekNextStep(), "unexpected non numeric step for slice state '%+v'", sliceToVisit)
}

func (w *walkOpVisitor) visitSliceWithIndex(sv stepVisitor, p path, sliceToVisit []interface{}, index int) (interface{}, error) {
	interfaceToVisit, err := w.visitLastSlice(sliceToVisit, index)
	if err != nil {
		return nil, wrapWithPathError(p.currentStep, err)
	}
	return sv.visitNextStep(p, interfaceToVisit)
}

func (w *walkOpVisitor) visitMapWithIdentity(sv stepVisitor, p path,
	mapToVisit map[string]interface{}, identity string) (interface{}, error) {
	interfaceToVisit, err := w.visitLastMap(mapToVisit, identity)
	if err != nil {
		return nil, wrapWithPathError(p.currentStep, err)
	}
	return sv.visitNextStep(p, interfaceToVisit)
}
