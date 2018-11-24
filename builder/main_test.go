package builder_test

import "github.com/scottshotgg/express2/builder"

type TestType int

const (
	_ TestType = iota

	ExpressionTest
	StatementTest

	programASTString = "programAST %+v\n"
	errString        = "err %+v\n"
)

var (
	b          *builder.Builder
	programAST *builder.Node
	err        error

	expressionTestMap = map[string]string{
		"deref":       "*somethingElse",
		"ident":       "thisIsAnIdent",
		"inc":         "i++",
		"condition":   "something < 10 < (7)",
		"array":       "[ \"something\", [8, 8], 9, i ]",
		"intLit":      "7",
		"identIndex":  "something[9][0]",
		"identCall":   "funcYou(6, 7)",
		"blockExpr":   "{ int i = 7 }",
		"identSelect": "some.thing.whatever.yeah",
		"binop":       "9 + 8 * 7",
	}

	statementTestMap = map[string]string{
		"decl":             "int i = 10",
		"ifElse":           "if something { } else if somethingElse { }",
		"funcDef":          "function something(int i, string s) int { return 10.1 }",
		"simpleAssign":     "i = 10",
		"callAssign":       "something = something(5, i, s)",
		"block":            "{ int i = 10 int j = 99 }",
		"import":           "import \"somethingHere.expr\"",
		"include":          "include \"somethingHere.expr\"",
		"stdFor":           "for int i = 1; i < 10; i++ { }",
		"arrayDef":         "int[] i = [ 8, 9, 0 ]",
		"forin":            "for i in [ 7, 8, 9 ] { }",
		"forof":            "for i of [ 7, 8, 9 ] { }",
		"indexAssign":      "something[7] = \"hey its me\"",
		"assignFromIndex":  "something = here[9][0]",
		"typeDef":          "type myInt = int",
		"selectionAssign":  "some.thing.whatever.yeah = 10",
		"assignFromSelect": "somethingNew = some.thing",
		"returnSomething":  "return something[\"here\"].me()",
		"struct":           "struct something = { int i = 10 }",
		"simpleLet":        "let something = 99",
		"package":          "package something",
		"derefAssign":      "*something = 10",
		"binop":            "i = 9 + 8 * 7",
	}

	tests = map[TestType]map[string]string{
		ExpressionTest: expressionTestMap,
		StatementTest:  statementTestMap,
	}
)
