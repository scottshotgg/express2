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

func TestAddPrimitive(t *testing.T) {
	// test := "float"
	test := "type myInt = int"

	b, err = getBuilderFromString(test)
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseTypeDeclarationStatement()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)

	fmt.Println(b.AddPrimitive(node.Left.Value.(string), node.Right))

	fmt.Println(b.ScopeTree)
}

func TestAddStructured(t *testing.T) {
	var test = "struct thing = { int i = 7 string s = \"hey\"}"

	b, err = getBuilderFromString(test)
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	node, err = b.ParseStructStatement()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)

	var thing = b.ScopeTree.GetType("thing")
	fmt.Println("thing", thing)
	nodeJSON, _ = json.Marshal(thing)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func TestBuildNodeFromTypeValue(t *testing.T) {
	TestAddStructured(t)

	var node, err = b.BuildNodeFromTypeValue(b.ScopeTree.GetType("thing"))
	if err != nil {
		t.Fatalf("err %+v", err)
	}

	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}
