package builder_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
)

var (
	testBytes []byte
)

func TestProgram(t *testing.T) {
	testBytes, err = ioutil.ReadFile("test.expr")
	if err != nil {
		t.Fatalf("Could not read file: "+errFormatString, err)
	}

	// keep this string here for injection
	test := string(testBytes)

	b, err = getBuilderFromString(test)
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.BuildAST()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}
