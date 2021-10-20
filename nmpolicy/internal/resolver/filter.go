package resolver

import (
	"fmt"

	"github.com/nmstate/nmpolicy/nmpolicy/internal/ast"
)

func filter(inputState map[string]interface{}, path ast.VariadicOperator, expectedNode ast.Node) (map[string]interface{}, error) {
	filtered, err := applyFuncOnPath(inputState, path, expectedNode, mapContainsValue, true)

	if err != nil {
		return nil, filterError("error applying operation on the path : %v", err)
	}

	if filtered == nil {
		return nil, nil
	}

	filteredMap, ok := filtered.(map[string]interface{})
	if !ok {
		return nil, filterError("error converting filtering result to a map")
	}
	return filteredMap, nil
}

func isEqual(obtainedValue interface{}, desiredValue ast.Node) (bool, error) {
	if desiredValue.String != nil {
		stringToCompare, ok := obtainedValue.(string)
		if !ok {
			return false, fmt.Errorf("the value %v of type %T not supported,"+
				"curretly only string values are supported", obtainedValue, obtainedValue)
		}
		return stringToCompare == *desiredValue.String, nil
	}

	return false, fmt.Errorf("the desired value %v is not supported. Curretly only string values are supported", desiredValue)
}

func mapContainsValue(mapToFilter map[string]interface{}, filterKey string, expectedNode ast.Node) (interface{}, error) {
	obtainedValue, ok := mapToFilter[filterKey]
	if !ok {
		return nil, filterError("cannot find key %s in %v", filterKey, obtainedValue)
	}
	valueIsEqual, err := isEqual(obtainedValue, expectedNode)
	if err != nil {
		return nil, filterError("error comparing the expected and obtained values : %v", err)
	}
	if valueIsEqual {
		return mapToFilter, nil
	}
	return nil, nil
}
