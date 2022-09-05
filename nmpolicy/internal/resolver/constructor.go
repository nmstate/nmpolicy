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

import "fmt"

func constructState(p path, inputState, prevStep interface{}) (interface{}, error) {
	var err error
	identity, isString := prevStep.(*string)
	if isString {
		originalMap, isMap := inputState.(map[string]interface{})
		if isMap {
			if p.hasMoreSteps() {
				originalMap[*identity] = map[string]interface{}{}
				originalMap[*identity], err = constructState(p.nextStepByRef(), originalMap[*identity], p.currentStep.Identity)
				if err != nil {
					return nil, fmt.Errorf("failed to construct missing path in inputState '%+v', error: %s", inputState, err.Error())
				}
				return originalMap, nil
			} else {
				var i interface{}
				originalMap[*identity] = i
				return originalMap, nil
			}
		} else {
			return nil, fmt.Errorf("only a map structure can be constructed, got %T", inputState)
		}
	} else {
		return nil, fmt.Errorf("constructState only supports string map keys, got '%+v", identity)
	}
}
