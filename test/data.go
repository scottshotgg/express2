package test

type TestType int

const (
	_ TestType = iota

	ExpressionTest
	StatementTest
)

var (
	expressionTestMap = map[string]string{
		"deref":       "*somethingElse",
		"ident":       "thisIsAnIdent",
		"inc":         "i++",
		"condition":   "something < 10 < (7)",
		"array":       "[ \"something\", [8, 8], 9, i ]",
		"intLitArray": "[ 4, 6, 9, 7 ]",
		"intLit":      "7",
		"identIndex":  "something[9][0]",
		"identCall":   "funcYou(now() + 7)",
		"blockExpr":   "{ int i = 7 }",
		"identSelect": "some.thing.whatever.yeah",
		"binop":       "9 + 8 * 7",
		"equality":    "i == 0",
	}

	statementTestMap = map[string]string{
		"pointerDeclaration": `func main() {
			*int i
			}`,
		"isEqualBool": "bool b = 2 + 2 == 3 + 3",
		"sgroup":      "(int i, string s)",
		"decl":        "int i = 10",
		"ifElse": `  if something {
    int x = 7
  } else if true {
    string y = "1000000" + true
  } else {
    launch something()
  }`,
		"anotherFuncDef": `
		func delayedPrintln() {
			msleep(now() + 1000)
			Println("hi")
		  }
		`,
		"funcDef":       "func something(int i, string s) int { return 10 }",
		"simpleAssign":  "i = 10",
		"callNonAssign": "c.fputs(5, i, s)",
		"callAssign":    "something = something(5, i, s)",
		"block":         "{ int i = 10 int j = 99 }",
		"import":        "import \"../compiler/test/something.expr\"",
		"cimport":       "import c",
		"include":       "include \"somethingHere.expr\"",
		"stdFor":        "for int i = 1; i < 10; i++ { int k = 10 }",
		"arrayDef":      "int[] i = [ 8, 9, 0 ]",
		// "forin":            "for i in is { i = 10 }",
		"forin":            "for i in [ 7, 8, 9 ] { j = 10 }",
		"forof":            "for i of [ 7, 8, 9 ] { int i = 10 }",
		"indexAssign":      "something[7] = \"hey its me\"",
		"assignFromIndex":  "something = here[9][0]",
		"typeDef":          "type myInt = int",
		"selectionAssign":  "some.thing.whatever.yeah = 10",
		"assignFromSelect": "somethingNew = some.thing",
		"returnSomething":  "return something[\"here\"].me()",
		"struct":           "struct something = { int i = 10 string s = \"hey its me\" }",
		"simpleLet":        "let something = 99",
		"package":          "package something",
		"derefAssign":      "*something = 10",
		"binop":            "i = 9 + 8 * 7",
		"object":           "object o = { int a = 6 }",
	}

	Tests = map[TestType]map[string]string{
		ExpressionTest: expressionTestMap,
		StatementTest:  statementTestMap,
	}
)
