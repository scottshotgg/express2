package builder_test

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/scottshotgg/express2/builder"
)

var teststring = `{"Type":"program","Value":[{"Type":"import","Kind":"c"},{"Type":"function","Kind":"main","Value":{"Type":"block","Value":[{"Type":"enum","Value":{"Type":"ident","Value":"Gender"},"Left":{"Type":"block","Value":[{"Type":"ident","Value":"Male"},{"Type":"ident","Value":"Female"},{"Type":"ident","Value":"Helicopter"}]}},{"Type":"struct","Left":{"Type":"ident","Value":"Person"},"Right":{"Type":"block","Kind":"struct","Value":[{"Type":"decl","Value":{"Type":"ident","Value":"string"},"Left":{"Type":"ident","Value":"Name"},"Right":{"Type":"literal","Kind":"string","Value":""}},{"Type":"decl","Value":{"Type":"ident","Value":"int"},"Left":{"Type":"ident","Value":"Age"},"Right":{"Type":"literal","Kind":"int","Value":0}},{"Type":"decl","Value":{"Type":"ident","Value":"int"},"Left":{"Type":"ident","Value":"Gender"},"Right":{"Type":"ident","Value":"Male"}},{"Type":"decl","Value":{"Type":"ident","Value":"map"},"Left":{"Type":"ident","Value":"Characteristics"},"Right":{"Type":"block","Value":[{"Type":"kv","Left":{"Type":"literal","Kind":"int","Value":444},"Right":{"Type":"literal","Kind":"int","Value":222}}]}}]}},{"Type":"assignment","Left":{"Type":"selection","Left":{"Type":"ident","Value":"c"},"Right":{"Type":"binop","Value":"*","Left":{"Type":"ident","Value":"FILE"},"Right":{"Type":"ident","Value":"f"}}},"Right":{"Type":"selection","Left":{"Type":"ident","Value":"c"},"Right":{"Type":"call","Value":{"Type":"ident","Value":"fopen"},"Metadata":{"args":{"Type":"egroup","Value":[{"Type":"literal","Kind":"string","Value":"something"},{"Type":"literal","Kind":"string","Value":"w+"}]}}}}},{"Type":"decl","Value":{"Type":"ident","Value":"map"},"Left":{"Type":"ident","Value":"chars"},"Right":{"Type":"block","Value":[{"Type":"kv","Left":{"Type":"literal","Kind":"string","Value":"IsProAF"},"Right":{"Type":"literal","Kind":"bool","Value":true}},{"Type":"kv","Left":{"Type":"literal","Kind":"int","Value":69},"Right":{"Type":"literal","Kind":"string","Value":"truth"}}]}},{"Type":"decl","Value":{"Type":"ident","Value":"Person"},"Left":{"Type":"ident","Value":"test"},"Right":{"Type":"block","Value":[{"Type":"assignment","Left":{"Type":"ident","Value":"Name"},"Right":{"Type":"literal","Kind":"string","Value":"scott"}},{"Type":"assignment","Left":{"Type":"ident","Value":"Age"},"Right":{"Type":"literal","Kind":"int","Value":24}},{"Type":"assignment","Left":{"Type":"ident","Value":"Gender"},"Right":{"Type":"ident","Value":"Helicopter"}},{"Type":"assignment","Left":{"Type":"ident","Value":"Characteristics"},"Right":{"Type":"ident","Value":"chars"}}]}},{"Type":"decl","Value":{"Type":"ident","Value":"string"},"Left":{"Type":"ident","Value":"output"},"Right":{"Type":"selection","Left":{"Type":"ident","Value":"test"},"Right":{"Type":"binop","Value":"+","Left":{"Type":"ident","Value":"Name"},"Right":{"Type":"literal","Kind":"string","Value":" is a bawss"}}}},{"Type":"selection","Left":{"Type":"ident","Value":"c"},"Right":{"Type":"call","Value":{"Type":"ident","Value":"fputs"},"Metadata":{"args":{"Type":"egroup","Value":[{"Type":"selection","Left":{"Type":"ident","Value":"output"},"Right":{"Type":"call","Value":{"Type":"ident","Value":"c_str"},"Metadata":{"args":{"Type":"egroup","Value":null}}}},{"Type":"ident","Value":"f"}]}}}},{"Type":"selection","Left":{"Type":"ident","Value":"c"},"Right":{"Type":"call","Value":{"Type":"ident","Value":"fclose"},"Metadata":{"args":{"Type":"egroup","Value":[{"Type":"ident","Value":"f"}]}}}}]},"Metadata":{"args":{"Type":"sgroup","Value":null},"returns":{"Type":"type","Value":"int"}}}]}`

func TestChecker(t *testing.T) {
	var contents, err = ioutil.ReadFile("simple_import.expr")
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	b, err = getBuilderFromString(string(contents))
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	ast, err := b.BuildAST()
	if err != nil {
		t.Errorf(errFormatString, err)
	}

	fmt.Println("\n\nstarting check")
	fmt.Println()

	var ch = builder.NewChecker(ast, builder.NewDummyPass("first_dummy"))

	ch.AddPass(builder.NewTypeResolver())

	ch.Execute()
}
