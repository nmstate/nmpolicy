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

type pathVisitor struct {
	path            []ast.Node
	pathIndex       int
	currentStep     *ast.Node
	lastMapFn       mapEntryVisitFn
	filterLookupMap map[string]bool
}

func (v pathVisitor) visitInterface(inputState interface{}) (interface{}, error) {
	originalMap, isMap := inputState.(map[string]interface{})
	if isMap {
		v.nextStep()
		if v.hasMoreSteps() {
			return v.visitMap(originalMap)
		}
		return v.visitLastMapOnPath(originalMap, inputState)
	}

	originalSlice, isSlice := inputState.([]interface{})
	if isSlice {
		return v.visitSlice(originalSlice)
	}

	return nil, pathError(v.currentStep, "invalid type %T for identity step '%v'", inputState, *v.currentStep)
}

func (v pathVisitor) visitSlice(originalSlice []interface{}) (interface{}, error) {
	return v.visitSliceWithoutIndex(originalSlice)
}

func (v pathVisitor) visitSliceWithoutIndex(originalSlice []interface{}) (interface{}, error) {
	adjustedSlice := []interface{}{}
	sliceEmptyAfterApply := true
	for _, valueToCheck := range originalSlice {
		valueAfterApply, err := v.visitInterface(valueToCheck)
		if err != nil {
			return nil, err
		}
		if valueAfterApply != nil {
			sliceEmptyAfterApply = false
			adjustedSlice = append(adjustedSlice, valueAfterApply)
		} else if !v.shouldFilter() {
			adjustedSlice = append(adjustedSlice, valueToCheck)
		}
	}

	if sliceEmptyAfterApply {
		return nil, nil
	}

	return adjustedSlice, nil
}

func (v pathVisitor) visitMap(originalMap map[string]interface{}) (interface{}, error) {
	if v.currentStep.Identity == nil {
		return nil, pathError(v.currentStep, "%v has unsupported fromat", *v.currentStep)
	}
	return v.visitMapWithIdentity(originalMap, *v.currentStep.Identity)
}

func (v pathVisitor) visitMapWithIdentity(originalMap map[string]interface{}, identity string) (interface{}, error) {
	valueToCheck, ok := originalMap[identity]
	if !ok {
		return nil, nil
	}

	adjustedValue, err := v.visitInterface(valueToCheck)
	if err != nil {
		return nil, err
	}
	if adjustedValue == nil {
		return nil, nil
	}

	adjustedMap := map[string]interface{}{}
	if !v.shouldFilter() {
		for k, v := range originalMap {
			adjustedMap[k] = v
		}
	}
	adjustedMap[identity] = adjustedValue
	return adjustedMap, nil
}

func (v pathVisitor) visitLastMapOnPath(originalMap map[string]interface{}, inputState interface{}) (interface{}, error) {
	if v.lastMapFn != nil {
		key := *v.currentStep.Identity
		outputState, err := v.lastMapFn(originalMap, key)
		if err != nil {
			return nil, err
		}
		return outputState, nil
	}
	return inputState, nil
}

func (v *pathVisitor) nextStep() {
	if v.currentStep == nil {
		v.pathIndex = 0
	} else if v.hasMoreSteps() {
		v.pathIndex++
	}
	v.currentStep = &v.path[v.pathIndex]
}

func (v *pathVisitor) shouldFilter() bool {
	if v.currentStep == nil || v.currentStep.Identity == nil {
		return false
	}
	_, ok := v.filterLookupMap[*v.currentStep.Identity]
	return ok
}

func (v pathVisitor) hasMoreSteps() bool {
	return v.pathIndex+1 < len(v.path)
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
