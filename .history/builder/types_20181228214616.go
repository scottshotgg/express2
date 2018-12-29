package builder

import (
	"github.com/scottshotgg/express-token"
)

type (
	opCallbackFn func(n *Node) (*Node, error)

	Index struct {
		Type  string      `json:",omitempty"`
		Value interface{} `json:",omitempty"`
	}

	Node struct {
		Type     string                 `json:",omitempty"`
		Kind     string                 `json:",omitempty"`
		Value    interface{}            `json:",omitempty"`
		Metadata map[string]interface{} `json:",omitempty"`
		Left     *Node                  `json:",omitempty"`
		Right    *Node                  `json:",omitempty"`
	}

	TypeValueType int

	// Need to solve for array here but for now its w/e;
	// if they want to do a typedef for like `type 3dCoord string[3]`
	// we would make Composite true, Type "array", Props would be 0, 1, 2, 3, 4, with the base type etc
	TypeValue struct {
		Composite bool
		Type      string
		Props     []*TypeValue
	}

	Builder struct {
		Tokens []token.Token `json:",omitempty"`
		Index  int           `json:",omitempty"`
		// [op_tier][op] -> func
		OpFuncMap []map[string]opCallbackFn `json:",omitempty"`

		ScopeTree *ScopeTree
		TypeMap   *map[string]*TypeValue
	}
)

// Might have to add one in here for objects as well
const (
	_ TypeValueType = iota
	PrimitiveValue
	RepeatedValue
	StruturedValue
)
