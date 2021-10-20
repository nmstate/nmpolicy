package resolver

import "fmt"

func filterError(format string, a ...interface{}) error {
	return fmt.Errorf("invalid filter: %v", fmt.Errorf(format, a...))
}

func pathError(format string, a ...interface{}) error {
	return fmt.Errorf("invalid path: %v", fmt.Errorf(format, a...))
}

func wrapWithResolveError(err error) error {
	return fmt.Errorf("resolve error: %v", err)
}
