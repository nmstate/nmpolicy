package resolver

import "fmt"

func filterError(err error) error {
	return fmt.Errorf("invalid filter: %v", err)
}

func pathError(err error) error {
	return fmt.Errorf("invalid path: %v", err)
}
