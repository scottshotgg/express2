package symbol_table_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"io/ioutil"

	"github.com/scottshotgg/express2/symbol_table"
)

const (
	astFormatString  = "ast: %+v\n"
	errFormatString  = "err: %+v\n"
	jsonFormatString = "JSON: %s\n"
)

func TestStuff(t *testing.T) {
	testBytes, err := ioutil.ReadFile("test.expr")
	if err != nil {
		t.Fatalf("Could not read file: "+errFormatString, err)
	}

	// keep this string here for injection
	test := string(testBytes)

	b, err := getBuilderFromString(test)
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err := b.BuildAST()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	nodeJSON, _ := json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)

	symbol_table.Stuff(node)
}
