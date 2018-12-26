package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/scottshotgg/express-ast"
	"github.com/scottshotgg/express-lex"
	"github.com/scottshotgg/express2/builder"
)

var (
	simpleTest = `
	function something() {
		stuff := "woah"
		var thing = "yeah"
	}

	ay()

	int something = 3 - 89

	some := 1 + i < 0

	for i := 0, i < 10, i++ {
		something := "here"
	}

	for i in something {
		something := "here"
		return true
	}

	{}

	{
		woah := "random string"
	}

	int i
	if something {
		i = 0
	} else {
		i = 1
	}
	
	char[] me
	`
	// There is a problem with adding the extra newline and tab tokens
)

func main() {
	file := "test/main/test.expr"

	lexer, err := lex.NewFromFile(file)
	if err != nil {
		fmt.Println("Could not find file:", file)
		os.Exit(9)
	}

	fmt.Printf("\nLexing file: %s\n", file)

	// Lex and tokenize the source code
	tokens, err := lexer.Lex()
	if err != nil {
		fmt.Println("err", err)
	}

	fmt.Println("\nCompressing tokens ...")

	// Compress certain tokens;
	// i.e: `:` and `=` compress into `:=`
	tokens, err = ast.CompressTokens(tokens)
	if err != nil {
		fmt.Println("err", err)
	}

	for _, t := range tokens {
		fmt.Println(t)
	}

	programAST, err := builder.New(tokens).BuildAST()
	if err != nil {
		fmt.Printf("err %+v\n", err)
		os.Exit(9)
	}

	pJSON, _ := json.Marshal(programAST)
	fmt.Println()
	fmt.Println(string(pJSON))
	fmt.Println()

	// p2 := *ast.Program{}
	// err = json.Unmarshal(pJSON, p2)
	// if err != nil {
	// 	fmt.Printf("err %+v\n", err)
	// 	os.Exit(9)
	// }

	// fmt.Println("Type checking Express code ...")

	// err = typecheck.TypeCheck(programAST)
	// if err != nil {
	// 	fmt.Printf("\nerr %+v\n", err)
	// 	os.Exit(9)
	// }

	// fmt.Println("Transpiling Express to C++ code ...")

	// // Transpile the AST into C++
	// t, err := transpiler.Transpile(programAST)
	// if err != nil {
	// 	fmt.Printf("\nerr %+v\n", err)
	// 	os.Exit(9)
	// }

	// fmt.Println("\nWriting transpilied C++ code to main.cpp ...")

	// // Write the C++ code to a file named `main.cpp`
	// err = ioutil.WriteFile("main.cpp", []byte(t), 0644)
	// if err != nil {
	// 	fmt.Printf("\nerr %+v\n", err)
	// 	os.Exit(9)
	// }

	// fmt.Println("\nFormatting C++ code ...")

	// // Run `clang-format` in-place to format the file for human-readability
	// output, err := exec.Command("clang-format", "-i", "main.cpp").CombinedOutput()
	// if err != nil {
	// 	fmt.Printf("%s\n%+v\n", output, err)
	// 	os.Exit(9)
	// }

	// fmt.Println("\nCompiling C++ code ...")

	// // Compile the file with Clang to produce a binary
	// output, err = exec.Command("clang++", "main.cpp", "-o", "main").CombinedOutput()
	// if err != nil {
	// 	fmt.Printf("%s\n%+v\n", output, err)
	// 	os.Exit(9)
	// }

	// fmt.Println("\nDone!")
}
