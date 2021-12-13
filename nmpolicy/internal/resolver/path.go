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

type applyFn func(map[string]interface{}, string) (interface{}, error)

func applyFuncOnPath(inputState interface{},
	path []ast.Node,
	funcToApply applyFn,
	shouldFilterSlice bool, shouldFilterMap bool) (interface{}, error) {
	if len(path) == 0 {
		return inputState, nil
	}
	originalMap, isMap := inputState.(map[string]interface{})
	if isMap {
		if len(path) == 1 {
			return applyFuncOnLastMapOnPath(path, originalMap, inputState, funcToApply)
		}
		return applyFuncOnMap(path, originalMap, funcToApply, shouldFilterSlice, shouldFilterMap)
	}

	originalSlice, isSlice := inputState.([]interface{})
	if isSlice {
		return applyFuncOnSlice(originalSlice, path, funcToApply, shouldFilterSlice)
	}

	return nil, pathError("invalid type %T for identity step '%v'", inputState, path[0])
}

func applyFuncOnSlice(originalSlice []interface{},
	path []ast.Node,
	funcToApply applyFn,
	shouldFilterSlice bool) (interface{}, error) {
	adjustedSlice := []interface{}{}
	sliceEmptyAfterApply := true
	for _, valueToCheck := range originalSlice {
		valueAfterApply, err := applyFuncOnPath(valueToCheck, path, funcToApply, false, false)
		if err != nil {
			return nil, err
		}
		if valueAfterApply != nil {
			sliceEmptyAfterApply = false
			adjustedSlice = append(adjustedSlice, valueAfterApply)
		} else if !shouldFilterSlice {
			adjustedSlice = append(adjustedSlice, valueToCheck)
		}
	}

	if sliceEmptyAfterApply {
		return nil, nil
	}

	return adjustedSlice, nil
}

func applyFuncOnMap(path []ast.Node,
	originalMap map[string]interface{},
	funcToApply applyFn,
	shouldFilterSlice bool, shouldFilterMap bool) (interface{}, error) {
	currentStep := path[0]
	if currentStep.Identity == nil {
		return nil, pathError("%v has unsupported fromat", currentStep)
	}

	nextPath := path[1:]
	key := *currentStep.Identity

	valueToCheck, ok := originalMap[key]
	if !ok {
		return nil, pathError("cannot find key %s in %v", key, originalMap)
	}

	adjustedValue, err := applyFuncOnPath(valueToCheck, nextPath, funcToApply, shouldFilterSlice, shouldFilterMap)
	if err != nil {
		return nil, err
	}
	if adjustedValue == nil {
		return nil, nil
	}

	adjustedMap := map[string]interface{}{}
	if !shouldFilterMap {
		for k, v := range originalMap {
			adjustedMap[k] = v
		}
	}
	adjustedMap[key] = adjustedValue
	return adjustedMap, nil
}

func applyFuncOnLastMapOnPath(path []ast.Node,
	originalMap map[string]interface{},
	inputState interface{},
	funcToApply applyFn) (interface{}, error) {
	if funcToApply != nil {
		key := *path[0].Identity
		outputState, err := funcToApply(originalMap, key)
		if err != nil {
			return nil, err
		}
		return outputState, nil
	}
	return inputState, nil
}

func (p captureEntryNameAndSteps) walkState(stateToWalk map[string]interface{}) (interface{}, error) {
	var (
		walkedState interface{}
		walkedPath  []string
	)
	walkedState = stateToWalk
	for _, step := range p.steps {
		if step.Identity != nil {
			identityStep := *step.Identity
			walkedPath = append(walkedPath, identityStep)
			walkedStateMap, ok := walkedState.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("failed walking non map state '%+v' with path '%+v'", walkedState, walkedPath)
			}
			walkedState, ok = walkedStateMap[identityStep]
			if !ok {
				return nil, fmt.Errorf("step '%s' from path '%s' not found at map state '%+v'", identityStep, walkedPath, walkedStateMap)
			}
		} else if step.Number != nil {
			numberStep := *step.Number
			walkedPath = append(walkedPath, strconv.Itoa(numberStep))
			walkedStateSlice, ok := walkedState.([]interface{})
			if !ok {
				return nil, fmt.Errorf("failed walking non slice state '%+v' with path '%+v'", walkedState, walkedPath)
			}
			if len(walkedStateSlice) == 0 || numberStep >= len(walkedStateSlice) {
				return nil, fmt.Errorf("step '%d' from path '%s' not found at slice state '%+v'", numberStep, walkedPath, walkedStateSlice)
			}
			walkedState = walkedStateSlice[numberStep]
		}
	}
	return walkedState, nil
}
