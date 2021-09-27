package capturer

import (
	"github.com/nmstate/nmpolicy/internal/ast"
	"github.com/nmstate/nmpolicy/internal/state"
)

type CapturedStateByName map[string]state.State
type CommandByName map[string]ast.Command
