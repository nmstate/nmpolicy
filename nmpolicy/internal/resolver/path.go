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
	"strconv"

	"github.com/nmstate/nmpolicy/nmpolicy/internal/ast"
)

type captureEntryNameAndSteps struct {
	captureEntryName string
	steps            ast.VariadicOperator
}

type mapEntryVisitFn func(map[string]interface{}, string) (interface{}, error)
type sliceEntryVisitFn func([]interface{}, int) (interface{}, error)

type pathVisitor struct {
	path                     []ast.Node
	currentStepIndex         int
	currentStep              *ast.Node
	lastMapFn                mapEntryVisitFn
	lastSliceFn              sliceEntryVisitFn
	visitSliceWithoutIndexFn func(pathVisitor, []interface{}) (interface{}, error)
	visitSliceWithIndexFn    func(pathVisitor, []interface{}, int) (interface{}, error)
	visitMapWithIdentityFn   func(pathVisitor, map[string]interface{}, string) (interface{}, error)
}

func (v pathVisitor) visitNextStep(inputState interface{}) (interface{}, error) {
	v.nextStep()
	originalMap, isMap := inputState.(map[string]interface{})
	if isMap {
		if v.hasMoreSteps() {
			return v.visitMap(originalMap)
		}
		return v.visitLastMapOnPath(originalMap, inputState)
	}

	originalSlice, isSlice := inputState.([]interface{})
	if isSlice {
		if v.hasMoreSteps() || v.currentStep.Number == nil {
			return v.visitSlice(originalSlice)
		}
		return v.visitLastSliceOnPath(originalSlice, inputState)
	}

	return nil, pathError(v.currentStep, "invalid type %T for identity step '%v'", inputState, *v.currentStep)
}

func (v pathVisitor) visitSlice(originalSlice []interface{}) (interface{}, error) {
	if v.currentStep.Number == nil {
		v.backStep()
		return v.visitSliceWithoutIndexFn(v, originalSlice)
	}
	return v.visitSliceWithIndexFn(v, originalSlice, *v.currentStep.Number)
}

func (v pathVisitor) visitMap(originalMap map[string]interface{}) (interface{}, error) {
	if v.currentStep.Identity == nil {
		return nil, pathError(v.currentStep, "unexpected non identity step for map state '%+v'", originalMap)
	}
	return v.visitMapWithIdentityFn(v, originalMap, *v.currentStep.Identity)
}

func (v pathVisitor) visitLastMapOnPath(originalMap map[string]interface{}, inputState interface{}) (interface{}, error) {
	if v.lastMapFn != nil {
		outputState, err := v.lastMapFn(originalMap, *v.currentStep.Identity)
		if err != nil {
			return nil, wrapWithPathError(v.currentStep, err)
		}
		return outputState, nil
	}
	return inputState, nil
}

func (v *pathVisitor) visitLastSliceOnPath(originalSlice []interface{}, inputState interface{}) (interface{}, error) {
	if v.lastSliceFn != nil {
		outputState, err := v.lastSliceFn(originalSlice, *v.currentStep.Number)
		if err != nil {
			return nil, wrapWithPathError(v.currentStep, err)
		}
		return outputState, nil
	}
	return inputState, nil
}

func (v *pathVisitor) nextStep() {
	if v.currentStep == nil {
		v.currentStepIndex = 0
	} else if v.hasMoreSteps() {
		v.currentStepIndex++
	}
	v.currentStep = &v.path[v.currentStepIndex]
}

func (v *pathVisitor) peekNextStep() *ast.Node {
	if !v.hasMoreSteps() {
		return v.currentStep
	}
	return &v.path[v.currentStepIndex+1]
}

func (v *pathVisitor) backStep() {
	if v.currentStep == nil {
		v.currentStepIndex = 0
	} else if v.currentStepIndex > 0 {
		v.currentStepIndex--
	}
	v.currentStep = &v.path[v.currentStepIndex]
}

func (v pathVisitor) hasMoreSteps() bool {
	return v.currentStepIndex+1 < len(v.path)
}

func (p captureEntryNameAndSteps) walkState(stateToWalk map[string]interface{}) (interface{}, error) {
	var (
		walkedState interface{}
		walkedPath  []string
	)
	walkedState = stateToWalk
	for _, step := range p.steps {
		node := step
		if step.Identity != nil {
			identityStep := *step.Identity
			walkedPath = append(walkedPath, identityStep)
			walkedStateMap, ok := walkedState.(map[string]interface{})
			if !ok {
				return nil, wrapWithPathError(&node, fmt.Errorf("failed walking non map state '%+v' with path '%+v'", walkedState, walkedPath))
			}
			walkedState, ok = walkedStateMap[identityStep]
			if !ok {
				return nil, wrapWithPathError(&node,
					fmt.Errorf("step '%s' from path '%s' not found at map state '%+v'", identityStep, walkedPath, walkedStateMap))
			}
		} else if step.Number != nil {
			numberStep := *step.Number
			walkedPath = append(walkedPath, strconv.Itoa(numberStep))
			walkedStateSlice, ok := walkedState.([]interface{})
			if !ok {
				return nil, wrapWithPathError(&node, fmt.Errorf("failed walking non slice state '%+v' with path '%+v'", walkedState, walkedPath))
			}
			if len(walkedStateSlice) == 0 || numberStep >= len(walkedStateSlice) {
				return nil, wrapWithPathError(&node,
					fmt.Errorf("step '%d' from path '%s' not found at slice state '%+v'", numberStep, walkedPath, walkedStateSlice))
			}
			walkedState = walkedStateSlice[numberStep]
		}
	}
	return walkedState, nil
}
