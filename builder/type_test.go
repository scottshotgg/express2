package builder_test

import (
	"encoding/json"
	"fmt"
	"testing"
)

// TODO:
// func TestParseArrayType(t *testing.T) {
// 	test := "int[][5]"

// 	b, err = getBuilderFromString(test)
// 	if err != nil {
// 		t.Errorf(errFormatString, err)
// 	}

// 	node, err = b.ParseArrayType()
// 	if err != nil {
// 		t.Errorf(errFormatString, err)
// 	}

// 		nodeJSON, _ = json.Marshal(node) 	fmt.Printf(jsonFormatString, nodeJSON)
// }

func TestParseType(t *testing.T) {
	// test := "float"
	test := "int[][5]"

	b, err = getBuilderFromString(test)
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseType()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}
