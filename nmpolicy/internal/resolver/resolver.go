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

	"sigs.k8s.io/yaml"

	"github.com/nmstate/nmpolicy/nmpolicy/internal/ast"
	"github.com/nmstate/nmpolicy/nmpolicy/types"
)

type Resolver struct {
	currentState   map[string]interface{}
	capturedStates map[string]map[string]interface{}
	captureASTPool map[string]ast.Node
}

func New() *Resolver {
	return &Resolver{
		currentState:   map[string]interface{}{},
		capturedStates: map[string]map[string]interface{}{},
		captureASTPool: map[string]ast.Node{},
	}
}

func (r *Resolver) Resolve(captureASTPool map[string]ast.Node, currentState []byte) (map[string]types.CaptureState, error) {
	r.captureASTPool = captureASTPool
	err := yaml.Unmarshal(currentState, &r.currentState)
	if err != nil {
		return nil, wrapWithResolveError(err)
	}

	capturedStates := map[string]types.CaptureState{}
	for captureEntryName := range r.captureASTPool {
		capturedStateEntry, err := r.resolveCaptureEntryName(captureEntryName)
		if err != nil {
			return nil, wrapWithResolveError(err)
		}
		marshaledCapturedStateEntry, err := yaml.Marshal(capturedStateEntry)
		if err != nil {
			return nil, wrapWithResolveError(err)
		}
		capturedStates[captureEntryName] = types.CaptureState{
			State:    marshaledCapturedStateEntry,
			MetaInfo: types.MetaInfo{},
		}
	}
	return capturedStates, nil
}

func (r *Resolver) resolveCaptureEntryName(captureEntryName string) (map[string]interface{}, error) {
	capturedStateEntry, ok := r.capturedStates[captureEntryName]
	if ok {
		return capturedStateEntry, nil
	}
	captureASTEntry, ok := r.captureASTPool[captureEntryName]
	if !ok {
		return nil, fmt.Errorf("capture entry '%s' not found", captureEntryName)
	}
	capturedStateEntry, err := r.resolveCaptureASTEntry(captureASTEntry)
	if err != nil {
		return nil, err
	}
	r.capturedStates[captureEntryName] = capturedStateEntry
	return capturedStateEntry, nil
}

func (r Resolver) resolveCaptureASTEntry(captureASTEntry ast.Node) (map[string]interface{}, error) {
	if captureASTEntry.EqFilter != nil {
		return r.resolveEqFilter(captureASTEntry.EqFilter)
	}
	return nil, fmt.Errorf("root node has unsupported operation : %v", captureASTEntry)
}

func (r Resolver) resolveEqFilter(operator *ast.TernaryOperator) (map[string]interface{}, error) {
	inputSource, err := r.resolveInputSource((*operator)[0], r.currentState)
	if err != nil {
		return nil, wrapWithEqFilterError(err)
	}

	path, err := r.resolvePath((*operator)[1])
	if err != nil {
		return nil, wrapWithEqFilterError(err)
	}
	filteredValue, err := r.resolveFilteredValue((*operator)[2])
	if err != nil {
		return nil, wrapWithEqFilterError(err)
	}
	filteredState, err := filter(inputSource, path.steps, *filteredValue)
	if err != nil {
		return nil, wrapWithEqFilterError(err)
	}
	return filteredState, nil
}

func (r Resolver) resolveInputSource(inputSourceNode ast.Node, currentState map[string]interface{}) (map[string]interface{}, error) {
	if ast.CurrentStateIdentity().DeepEqual(inputSourceNode.Terminal) {
		return currentState, nil
	}

	return nil, fmt.Errorf("not supported input source %v. Only the current state is supported", inputSourceNode)
}

func (r Resolver) resolveFilteredValue(filteredValueNode ast.Node) (*ast.Node, error) {
	if filteredValueNode.String != nil {
		return &filteredValueNode, nil
	} else if filteredValueNode.Path != nil {
		resolvedCaptureEntryPath, err := r.resolveCaptureEntryPath(filteredValueNode)
		if err != nil {
			return nil, err
		}
		terminal, err := newTerminalFromInterface(resolvedCaptureEntryPath)
		if err != nil {
			return nil, err
		}
		return &ast.Node{
			Terminal: *terminal,
		}, nil
	} else {
		return nil, fmt.Errorf("not supported filtered value. Only string or paths are supported")
	}
}

func (r Resolver) resolveCaptureEntryPath(pathNode ast.Node) (interface{}, error) {
	resolvedPath, err := r.resolvePath(pathNode)
	if err != nil {
		return nil, err
	}
	if resolvedPath.captureEntryName == "" {
		return nil, fmt.Errorf("not supported filtered value path. Only paths with a capture entry reference are supported")
	}
	capturedStateEntry, err := r.resolveCaptureEntryName(resolvedPath.captureEntryName)
	if err != nil {
		return nil, err
	}
	return resolvedPath.walkState(capturedStateEntry)
}

func (r Resolver) resolvePath(pathNode ast.Node) (*captureEntryNameAndSteps, error) {
	if pathNode.Path == nil {
		return nil, fmt.Errorf("invalid path type %T", pathNode)
	} else if len(*pathNode.Path) == 0 {
		return nil, fmt.Errorf("empty path length")
	} else if (*pathNode.Path)[0].Identity == nil {
		return nil, fmt.Errorf("path first step has to be an identity")
	}
	resolvedPath := captureEntryNameAndSteps{
		steps: *pathNode.Path,
	}
	if *resolvedPath.steps[0].Identity == "capture" {
		const captureRefSize = 2
		if len(resolvedPath.steps) < captureRefSize || resolvedPath.steps[1].Identity == nil {
			return nil, fmt.Errorf("path capture ref is missing capture entry name")
		}
		resolvedPath.captureEntryName = *resolvedPath.steps[1].Identity
		if len(resolvedPath.steps) > captureRefSize {
			resolvedPath.steps = resolvedPath.steps[2:len(resolvedPath.steps)]
		}
	}
	return &resolvedPath, nil
}
