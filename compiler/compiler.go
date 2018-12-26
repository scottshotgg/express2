package compiler

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/pkg/errors"
	ast "github.com/scottshotgg/express-ast"
	lex "github.com/scottshotgg/express-lex"
	token "github.com/scottshotgg/express-token"
	"github.com/scottshotgg/express2/builder"
	"github.com/scottshotgg/express2/transpiler"
)

const (
	stdCppVersion = "-std=c++2a"
)

func getTokensFromString(s string) ([]token.Token, error) {
	// Lex and tokenize the source code
	tokens, err := lex.New(s).Lex()
	if err != nil {
		return nil, err
	}

	fmt.Println("\nCompressing tokens ...")

	// Compress certain tokens;
	// i.e: `:` and `=` compress into `:=`
	return ast.CompressTokens(tokens)
}

func getBuilderFromString(test string) (*builder.Builder, error) {
	tokens, err := getTokensFromString(test)
	if err != nil {
		return nil, err
	}

	for _, token := range tokens {
		fmt.Println(token)
	}

	return builder.New(tokens), nil
}

func getASTFromString(test string) (*builder.Node, error) {
	b, err := getBuilderFromString(test)
	if err != nil {
		return nil, err
	}

	return b.BuildAST()
}

func getTranspilerFromString(test, name string) (*transpiler.Transpiler, error) {
	ast, err := getASTFromString(test)
	if err != nil {
		return nil, errors.Errorf("Could not create AST: %+v", err)
	}

	return transpiler.New(ast, name), nil
}

func Run(filename string) error {
	var testBytes, err = ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	tr, err := getTranspilerFromString(string(testBytes), "main")
	if err != nil {
		return err
	}

	cpp, err := tr.Transpile()
	if err != nil {
		return err
	}

	// fmt.Printf("\nC++: %s\n\n", cpp)

	fmt.Println("\nWriting transpilied C++ code to main.cpp ...")

	// Write the C++ code to a file named `main.cpp`
	err = ioutil.WriteFile("test/main.cpp", []byte(cpp), 0644)
	if err != nil {
		fmt.Printf("\nerr %+v\n", err)
		os.Exit(9)
	}

	fmt.Println("\nFormatting C++ code ...")

	// Run `clang-format` in-place to format the file for human-readability
	output, err := exec.Command("clang-format", "-i", "test/main.cpp").CombinedOutput()
	if err != nil {
		fmt.Printf("%s\n%+v\n", output, err)
		os.Exit(9)
	}

	fmt.Println("\nCompiling C++ code ...")

	// Compile the file with Clang to produce a binary
	output, err = exec.Command("clang++", stdCppVersion, "test/main.cpp", "-o", "test/main").CombinedOutput()
	if err != nil {
		fmt.Printf("%s\n%+v\n", output, err)
		os.Exit(9)
	}

	fmt.Println("\nDone!")

	return nil
}
