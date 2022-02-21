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

type captureEntryNameAndSteps struct {
	captureEntryName string
	steps            ast.VariadicOperator
}

type path struct {
	steps            []ast.Node
	currentStepIndex int
	currentStep      *ast.Node
}

type stepVisitor interface {
	visitNextStep(path, interface{}) (interface{}, error)
}

type stateVisitor interface {
	visitLastMap(map[string]interface{}, string) (interface{}, error)
	visitLastSlice([]interface{}, int) (interface{}, error)
	visitSliceWithoutIndex(stepVisitor, path, []interface{}) (interface{}, error)
	visitSliceWithIndex(stepVisitor, path, []interface{}, int) (interface{}, error)
	visitMapWithIdentity(stepVisitor, path, map[string]interface{}, string) (interface{}, error)
}

type pathVisitor struct {
	stateVisitor stateVisitor
}

func (v *pathVisitor) visitNextStep(p path, inputState interface{}) (interface{}, error) {
	p.nextStep()
	originalMap, isMap := inputState.(map[string]interface{})
	if isMap {
		if p.hasMoreSteps() {
			return v.visitMap(p, originalMap)
		}
		return v.visitLastMapOnPath(p, originalMap, inputState)
	}

	originalSlice, isSlice := inputState.([]interface{})
	if isSlice {
		if p.hasMoreSteps() || p.currentStep.Number == nil {
			return v.visitSlice(p, originalSlice)
		}
		return v.visitLastSliceOnPath(p, originalSlice, inputState)
	}

	return nil, pathError(p.currentStep, "invalid type %T for identity step '%v'", inputState, *p.currentStep)
}

func (v *pathVisitor) visitSlice(p path, originalSlice []interface{}) (interface{}, error) {
	if p.currentStep.Number == nil {
		p.backStep()
		return v.stateVisitor.visitSliceWithoutIndex(v, p, originalSlice)
	}
	return v.stateVisitor.visitSliceWithIndex(v, p, originalSlice, *p.currentStep.Number)
}

func (v *pathVisitor) visitMap(p path, originalMap map[string]interface{}) (interface{}, error) {
	if p.currentStep.Identity == nil {
		return nil, pathError(p.currentStep, "unexpected non identity step for map state '%+v'", originalMap)
	}
	return v.stateVisitor.visitMapWithIdentity(v, p, originalMap, *p.currentStep.Identity)
}

func (v *pathVisitor) visitLastMapOnPath(p path, originalMap map[string]interface{}, _ interface{}) (interface{}, error) {
	outputState, err := v.stateVisitor.visitLastMap(originalMap, *p.currentStep.Identity)
	if err != nil {
		return nil, wrapWithPathError(p.currentStep, err)
	}
	return outputState, nil
}

func (v *pathVisitor) visitLastSliceOnPath(p path, originalSlice []interface{}, _ interface{}) (interface{}, error) {
	outputState, err := v.stateVisitor.visitLastSlice(originalSlice, *p.currentStep.Number)
	if err != nil {
		return nil, wrapWithPathError(p.currentStep, err)
	}
	return outputState, nil
}

func (p *path) nextStep() {
	if p.currentStep == nil {
		p.currentStepIndex = 0
	} else if p.hasMoreSteps() {
		p.currentStepIndex++
	}
	p.currentStep = &p.steps[p.currentStepIndex]
}

func (p *path) backStep() {
	if p.currentStep == nil {
		p.currentStepIndex = 0
	} else if p.currentStepIndex > 0 {
		p.currentStepIndex--
	}
	p.currentStep = &p.steps[p.currentStepIndex]
}

func (p *path) hasMoreSteps() bool {
	return p.currentStepIndex+1 < len(p.steps)
}

func (p *path) peekNextStep() *ast.Node {
	if !p.hasMoreSteps() {
		return &p.steps[p.currentStepIndex]
	}
	return &p.steps[p.currentStepIndex+1]
}

func (p captureEntryNameAndSteps) walkState(stateToWalk map[string]interface{}) (interface{}, error) {
	return walk(stateToWalk, p.steps)
}
