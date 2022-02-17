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
	"github.com/nmstate/nmpolicy/nmpolicy/internal/ast"
)

func replace(inputState map[string]interface{}, pathSteps ast.VariadicOperator, replaceValue interface{}) (map[string]interface{}, error) {
	pathVisitorWithReplace := pathVisitor{
		lastMapFn:                replaceMapFieldValue(replaceValue),
		visitSliceWithoutIndexFn: replaceSliceEntriesWithVisitResult,
		visitMapWithIdentityFn:   replaceMapEntryWithVisitResult,
	}

	replaced, err := pathVisitorWithReplace.visitNextStep(path{steps: pathSteps}, inputState)

	if err != nil {
		return nil, replaceError("failed applying operation on the path: %w", err)
	}

	replacedMap, ok := replaced.(map[string]interface{})
	if !ok {
		return nil, replaceError("failed converting result to a map")
	}
	return replacedMap, nil
}

func replaceMapFieldValue(replaceValue interface{}) mapEntryVisitFn {
	return func(inputMap map[string]interface{}, mapEntryKeyToReplace string) (interface{}, error) {
		modifiedMap := map[string]interface{}{}
		for k, v := range inputMap {
			modifiedMap[k] = v
		}

		modifiedMap[mapEntryKeyToReplace] = replaceValue
		return modifiedMap, nil
	}
}

func replaceMapEntryWithVisitResult(v *pathVisitor, p path, mapToVisit map[string]interface{}, identity string) (interface{}, error) {
	interfaceToVisit, ok := mapToVisit[identity]
	if !ok {
		return nil, nil
	}

	visitResult, err := v.visitNextStep(p, interfaceToVisit)
	if err != nil {
		return nil, err
	}

	replacedMap := map[string]interface{}{}
	for k, v := range mapToVisit {
		replacedMap[k] = v
	}
	replacedMap[identity] = visitResult
	return replacedMap, nil
}

func replaceSliceEntriesWithVisitResult(v *pathVisitor, p path, sliceToVisit []interface{}) (interface{}, error) {
	replacedSlice := make([]interface{}, len(sliceToVisit))
	for i, interfaceToVisit := range sliceToVisit {
		visitResult, err := v.visitNextStep(p, interfaceToVisit)
		if err != nil {
			return nil, err
		}
		replacedSlice[i] = visitResult
	}
	return replacedSlice, nil
}
