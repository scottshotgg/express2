package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/scottshotgg/express-ast"
	"github.com/scottshotgg/express-lex"
	"github.com/scottshotgg/express2/transpiler"
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

	if something {
		int i = 0
	} else {
		int i = 1
	}
	
	char[] me`
	// There is a problem with adding the extra newline and tab tokens

	transpileTest = `int a = 5`
)

func main() {
	// Lex and tokenize the source code
	tokens, err := lex.New(transpileTest).Lex()
	if err != nil {
		fmt.Println("err", err)
	}

	// Compress certain tokens; i.e: `:` and `=` compress into `:=`
	tokens, err = ast.CompressTokens(tokens)
	if err != nil {
		fmt.Println("err", err)
	}

	for _, t := range tokens {
		fmt.Println(t)
	}

	// Make a builder
	builder := &ast.ASTBuilder{
		Tokens: tokens,
	}

	// Build the AST
	p, err := builder.BuildAST()
	if err != nil {
		fmt.Println("err", err)
		os.Exit(9)
	}

	pJSON, _ := json.Marshal(p)
	fmt.Println()
	fmt.Println(string(pJSON))
	fmt.Println()

	// Transpile the AST into C++
	t, err := transpiler.Transpile(p)
	if err != nil {
		fmt.Println("err", err)
		os.Exit(9)
	}

	// Write the C++ code to a file named `main.cpp`
	err = ioutil.WriteFile("main.cpp", []byte(t), 0755)
	if err != nil {
		fmt.Println("err", err)
		os.Exit(9)
	}

	// Run `clang-format` in-place to format the file for human-readability
	_, err = exec.Command("clang-format", "-i", "main.cpp").CombinedOutput()
	if err != nil {
		fmt.Println("err", err)
		os.Exit(9)
	}

	// Compile the file with Clang to produce a binary
	_, err = exec.Command("clang++", "main.cpp", "-o", "main").CombinedOutput()
	if err != nil {
		fmt.Println("err", err)
		os.Exit(9)
	}
}
