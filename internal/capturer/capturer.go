package capturer

import (
	"fmt"

	"github.com/nmstate/nmpolicy/internal/ast"
	"github.com/nmstate/nmpolicy/internal/state"
)

type Capturer struct {
	commandByName       CommandByName
	currentState        state.State
	capturedStateByName CapturedStateByName
}

func New(commandByName CommandByName, capturedStateByName CapturedStateByName) Capturer {
	return Capturer{
		commandByName:       commandByName,
		capturedStateByName: capturedStateByName,
	}
}

func (c *Capturer) Capture(currentState state.State) (CapturedStateByName, error) {
	c.currentState = currentState
	for commandName, commandNode := range c.commandByName {
		_, err := c.capture(commandName, commandNode)
		if err != nil {
			return nil, err
		}
	}
	return c.capturedStateByName, nil
}

func (c *Capturer) capture(name string, command ast.Command) (state.State, error) {
	capturedState, ok := c.capturedStateByName[name]
	if ok {
		return capturedState, nil
	}
	var err error
	capturedState, err = c.resolveCommand(c.currentState, command)
	if err != nil {
		return nil, err
	}
	c.capturedStateByName[name] = capturedState
	return capturedState, nil
}

func (c *Capturer) resolveCommand(s state.State, command ast.Command) (state.State, error) {
	capturedState := state.State{}
	if command.Equal != nil {
		var err error
		capturedState, err = c.filterStateByEquality(s, command.Equal)
		if err != nil {
			return nil, err
		}
	}
	return capturedState, nil

}

func (c *Capturer) filterStateByEquality(s state.State, arguments []ast.Argument) (state.State, error) {
	if len(arguments) != 2 {
		return nil, fmt.Errorf("invalid ast: number of arguments for equality command has to be two")
	}

	lhs := arguments[0]
	rhs := arguments[1]

	if len(lhs.Path) == 0 {
		return nil, fmt.Errorf("invalid ast: zero length at right hand side path argument for equality command")
	}

	// At foo.bar.dar "foo.bar" is the path and "bar" is the field to compare
	// within the slice of structs
	path := lhs.Path[:len(lhs.Path)-1]
	field := *lhs.Path[len(lhs.Path)-1].Identity

	filtered, err := filter(s, path, func(valueToFilter interface{}) (interface{}, error) {
		sliceToFilter, ok := valueToFilter.([]interface{})
		if ok {
			filteredSlice := []interface{}{}
			for _, valueToCheck := range sliceToFilter {
				matches, err := c.matchFieldValue(valueToCheck, field, rhs)
				if err != nil {
					return nil, err
				}
				if matches {
					filteredSlice = append(filteredSlice, valueToCheck)
				}
			}
			return filteredSlice, nil
		}
		return valueToFilter, nil
	})

	if err != nil {
		return nil, err
	}

	if filtered == nil {
		return nil, nil
	}

	return filtered.(state.State), nil
}

func filter(toFilter interface{}, path ast.Path, filterFn func(interface{}) (interface{}, error)) (interface{}, error) {
	if len(path) == 0 {
		if filterFn != nil {
			filtered, err := filterFn(toFilter)
			if err != nil {
				return nil, err
			}
			return filtered, nil
		}
		return toFilter, nil
	} else {
		currentStep := path[0]
		nextPath := path[1:]
		if currentStep.Identity != nil {
			key := *currentStep.Identity
			toFilterTyped, ok := toFilter.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("invalid path: invalid type %T for identity step '%s'", toFilter, key)
			}
			value, ok := toFilterTyped[key]
			if ok {
				value, err := filter(value, nextPath, filterFn)
				if err != nil {
					return nil, err
				}
				return map[string]interface{}{key: value}, nil
			} else {
				return nil, nil
			}
		} else if currentStep.Index != nil {
			key := *currentStep.Index
			toFilterTyped, ok := toFilter.([]interface{})
			if !ok {
				return nil, fmt.Errorf("invalid path: invalid type %T for index step '%d'", toFilter, key)
			}
			if key < len(toFilterTyped) {
				value := toFilterTyped[key]
				value, err := filter(value, nextPath, filterFn)
				if err != nil {
					return nil, err
				}
				return []interface{}{value}, nil
			} else {
				return nil, nil
			}
		} else {
			return toFilter, nil
		}
	}
}

func (c *Capturer) matchFieldValue(valueToCheck interface{}, field string, value ast.Argument) (bool, error) {
	mapToCheck, ok := valueToCheck.(map[string]interface{})
	if !ok {
		return false, fmt.Errorf("invalid equality command: non map type '%T' when accessing '%s'", valueToCheck, field)
	}
	valueToCompare, ok := mapToCheck[field]
	if !ok {
		return false, nil
	}
	if value.String != nil {
		stringToCompare, ok := valueToCompare.(string)
		if !ok {
			return false, nil
		}
		return stringToCompare == *value.String, nil
	}
	if value.Number != nil {
		numberToCompare, ok := valueToCompare.(int)
		if !ok {
			return false, nil
		}
		return numberToCompare == *value.Number, nil
	}
	return false, nil
}
