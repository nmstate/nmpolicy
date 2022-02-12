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

func eqFilter(inputState map[string]interface{}, path ast.VariadicOperator, expectedValue interface{}) (map[string]interface{}, error) {
	return filterWithPathVisitor(inputState, pathVisitor{
		path:              path,
		lastMapFn:         matchMapValue(expectedValue, eqMatcher),
		shouldFilterSlice: true,
		shouldFilterMap:   true,
		sliceVisitor:      notNilEntryFilteringSliceVisitor,
	})
}
func neFilter(inputState map[string]interface{}, path ast.VariadicOperator, expectedValue interface{}) (map[string]interface{}, error) {
	return filterWithPathVisitor(inputState, pathVisitor{
		path:              path,
		lastMapFn:         matchMapValue(expectedValue, neMatcher),
		shouldFilterSlice: true,
		shouldFilterMap:   true,
		sliceVisitor:      nilEntryFilteringSliceVisitor,
	})
}

func filterWithPathVisitor(inputState map[string]interface{}, pathVisitorWithFilter pathVisitor) (map[string]interface{}, error) {
	filtered, err := pathVisitorWithFilter.visitMap(inputState)

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

func matchMapValue(expectedValue interface{}, match func(interface{}, interface{}) bool) mapEntryVisitFn {
	return func(mapToFilter map[string]interface{}, filterKey string) (interface{}, error) {
		obtainedValue, ok := mapToFilter[filterKey]
		if !ok {
			return nil, nil
		}
		if reflect.TypeOf(obtainedValue) != reflect.TypeOf(expectedValue) {
			return nil, fmt.Errorf(`type missmatch: the value in the path doesn't match the value to filter. `+
				`"%T" != "%T" -> %+v != %+v`, obtainedValue, expectedValue, obtainedValue, expectedValue)
		}
		if match(obtainedValue, expectedValue) {
			return mapToFilter, nil
		}
		return nil, nil
	}
}

func eqMatcher(lhs, rhs interface{}) bool {
	return lhs == rhs
}

func neMatcher(lhs, rhs interface{}) bool {
	return lhs != rhs
}

func notNilEntryFilteringSliceVisitor(v pathVisitor, originalSlice []interface{}) (interface{}, error) {
	adjustedSlice := []interface{}{}
	matching := false
	for _, valueToCheck := range originalSlice {
		valueAfterApply, err := disableFiltering(v).visitInterface(valueToCheck)
		if err != nil {
			return nil, err
		}
		if valueAfterApply != nil {
			matching = true
			adjustedSlice = append(adjustedSlice, valueAfterApply)
		} else if !v.shouldFilterSlice {
			adjustedSlice = append(adjustedSlice, valueToCheck)
		}
	}

	if !matching {
		return nil, nil
	}

	return adjustedSlice, nil
}

func nilEntryFilteringSliceVisitor(v pathVisitor, originalSlice []interface{}) (interface{}, error) {
	adjustedSlice := []interface{}{}
	for _, valueToCheck := range originalSlice {
		valueAfterApply, err := disableFiltering(v).visitInterface(valueToCheck)
		if err != nil {
			return nil, err
		}
		if valueAfterApply != nil {
			adjustedSlice = append(adjustedSlice, valueAfterApply)
		} else if !v.shouldFilterSlice {
			return nil, nil
		}
	}

	return adjustedSlice, nil
}

func disableFiltering(v pathVisitor) pathVisitor {
	pathVisitorWithoutFilters := v
	pathVisitorWithoutFilters.shouldFilterSlice = false
	pathVisitorWithoutFilters.shouldFilterMap = false
	return pathVisitorWithoutFilters
}
