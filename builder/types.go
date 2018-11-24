package builder

import "github.com/scottshotgg/express-token"

type (
	opCallbackFn func(n *Node) (*Node, error)

	Index struct {
		Type  string
		Value interface{}
	}

	Node struct {
		Type     string
		Kind     string
		Value    interface{}
		Metadata map[string]interface{}
		Left     *Node
		Right    *Node
	}

	Builder struct {
		Tokens []token.Token
		Index  int
		// [op_tier][op] -> func
		OpFuncMap []map[string]opCallbackFn
	}
)
