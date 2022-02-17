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
	"reflect"

	"github.com/nmstate/nmpolicy/nmpolicy/internal/ast"
)

var (
	filterLookupMap = map[string]bool{
		"interfaces": true,
		"routes":     true,
		"running":    true,
		"config":     true,
	}
)

func filter(inputState map[string]interface{}, pathSteps ast.VariadicOperator, expectedValue interface{}) (map[string]interface{}, error) {
	pathVisitorWithEqFilter := pathVisitor{
		lastMapFn:                mapContainsValue(expectedValue),
		visitSliceWithoutIndexFn: filterSliceEntriesWithVisitResult,
		visitMapWithIdentityFn:   filterMapEntryWithVisitResult,
	}

	filtered, err := pathVisitorWithEqFilter.visitNextStep(path{steps: pathSteps}, inputState)

	if err != nil {
		return nil, fmt.Errorf("failed applying operation on the path: %w", err)
	}

	if filtered == nil {
		return nil, nil
	}

	filteredMap, ok := filtered.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("failed converting filtering result to a map")
	}
	return filteredMap, nil
}

func mapContainsValue(expectedValue interface{}) mapEntryVisitFn {
	return func(mapToFilter map[string]interface{}, filterKey string) (interface{}, error) {
		obtainedValue, ok := mapToFilter[filterKey]
		if !ok {
			return nil, nil
		}
		if reflect.TypeOf(obtainedValue) != reflect.TypeOf(expectedValue) {
			return nil, fmt.Errorf(`type missmatch: the value in the path doesn't match the value to filter. `+
				`"%T" != "%T" -> %+v != %+v`, obtainedValue, expectedValue, obtainedValue, expectedValue)
		}
		if obtainedValue == expectedValue {
			return mapToFilter, nil
		}
		return nil, nil
	}
}

func filterMapEntryWithVisitResult(v *pathVisitor, p path, mapToVisit map[string]interface{}, identity string) (interface{}, error) {
	interfaceToVisit, ok := mapToVisit[identity]
	if !ok {
		return nil, nil
	}

	visitResult, err := v.visitNextStep(p, interfaceToVisit)
	if err != nil {
		return nil, err
	}
	if visitResult == nil {
		return nil, nil
	}
	filteredMap := map[string]interface{}{}
	if !shouldFilter(p.currentStep) {
		for k, v := range mapToVisit {
			filteredMap[k] = v
		}
	}
	filteredMap[identity] = visitResult
	return filteredMap, nil
}

func filterSliceEntriesWithVisitResult(v *pathVisitor, p path, sliceToVisit []interface{}) (interface{}, error) {
	filteredSlice := []interface{}{}
	hasVisitResult := false
	for _, interfaceToVisit := range sliceToVisit {
		visitResult, err := v.visitNextStep(p, interfaceToVisit)
		if err != nil {
			return nil, err
		}
		if visitResult != nil {
			hasVisitResult = true
			filteredSlice = append(filteredSlice, visitResult)
		} else if !shouldFilter(p.currentStep) {
			filteredSlice = append(filteredSlice, interfaceToVisit)
		}
	}

	if !hasVisitResult {
		return nil, nil
	}
	return filteredSlice, nil
}

func shouldFilter(currentStep *ast.Node) bool {
	if currentStep == nil {
		return true
	}
	if currentStep.Identity == nil {
		return false
	}
	_, ok := filterLookupMap[*currentStep.Identity]
	return ok
}
