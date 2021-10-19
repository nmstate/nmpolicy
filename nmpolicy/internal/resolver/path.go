package resolver

import (
	"fmt"

	"github.com/nmstate/nmpolicy/nmpolicy/internal/ast"
)

func applyFuncOnPath(inputState interface{},
	path []ast.Node,
	expectedNode ast.Node,
	funcToApply func(map[string]interface{}, string, ast.Node) (interface{}, error),
	shouldFilterSlice bool) (interface{}, error) {
	if len(path) == 0 {
		return inputState, nil
	}
	originalMap, isMap := inputState.(map[string]interface{})
	if isMap {
		if len(path) == 1 {
			return applyFuncOnLastMapOnPath(path, originalMap, expectedNode, inputState, funcToApply)
		}
		return applyFuncOnMap(path, originalMap, expectedNode, funcToApply, shouldFilterSlice)
	}

	originalSlice, isSlice := inputState.([]interface{})
	if isSlice {
		return applyFuncOnSlice(originalSlice, path, expectedNode, funcToApply, shouldFilterSlice)
	}

	return nil, pathError(fmt.Errorf("invalid type %T for identity step '%v'", inputState, path[0]))
}

func applyFuncOnSlice(originalSlice []interface{},
	path []ast.Node,
	expectedNode ast.Node,
	funcToApply func(map[string]interface{}, string, ast.Node) (interface{}, error),
	shouldFilterSlice bool) (interface{}, error) {
	filteredSlice := []interface{}{}
	for _, valueToCheck := range originalSlice {
		value, err := applyFuncOnPath(valueToCheck, path, expectedNode, funcToApply, false)
		if err != nil {
			return nil, err
		}
		if value != nil {
			filteredSlice = append(filteredSlice, valueToCheck)
		}
	}

	if len(filteredSlice) == 0 {
		return nil, nil
	}

	if shouldFilterSlice {
		return filteredSlice, nil
	}
	return originalSlice, nil
}

func applyFuncOnMap(path []ast.Node,
	originalMap map[string]interface{},
	expectedNode ast.Node,
	funcToApply func(map[string]interface{}, string, ast.Node) (interface{}, error),
	shouldFilterSlice bool) (interface{}, error) {
	currentStep := path[0]
	if currentStep.Identity == nil {
		return nil, pathError(fmt.Errorf("%v has unsupported fromat", currentStep))
	}

	nextPath := path[1:]
	key := *currentStep.Identity

	value, ok := originalMap[key]
	if !ok {
		return nil, pathError(fmt.Errorf("cannot find key %s in %v", key, originalMap))
	}

	adjuctedValue, err := applyFuncOnPath(value, nextPath, expectedNode, funcToApply, shouldFilterSlice)
	if err != nil {
		return nil, err
	}
	if adjuctedValue != nil {
		return map[string]interface{}{key: adjuctedValue}, nil
	}
	return nil, nil
}

func applyFuncOnLastMapOnPath(path []ast.Node,
	originalMap map[string]interface{},
	expectedNode ast.Node,
	inputState interface{},
	funcToApply func(map[string]interface{}, string, ast.Node) (interface{}, error)) (interface{}, error) {
	if funcToApply != nil {
		key := *path[0].Identity
		outputState, err := funcToApply(originalMap, key, expectedNode)
		if err != nil {
			return nil, err
		}
		return outputState, nil
	}
	return inputState, nil
}
