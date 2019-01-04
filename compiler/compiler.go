package compiler

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
	"sync"
	"time"

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

var (
	libmill string

	pipelineTimes = map[string]string{}
	compilerFlags = []string{
		stdCppVersion,
		"-Ofast",
		// "-x",
		// "c++",
	}
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
	var tokens, err = getTokensFromString(test)
	if err != nil {
		return nil, err
	}

	// for _, token := range tokens {
	// 	fmt.Println(token)
	// }

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
	b, err := getBuilderFromString(test)
	if err != nil {
		return nil, err
	}

	ast, err := b.BuildAST()
	if err != nil {
		return nil, err
	}

	var astJSON, _ = json.Marshal(ast)

	fmt.Println("AST:", string(astJSON))

	return transpiler.New(ast, b, name), nil
}

func timeTrack(start time.Time, name string) {
	// fmt.Printf("Function %s took %s\n", name, time.Since(start))
	pipelineTimes[name] = time.Since(start).String()
}

func writeAndFormat(source, output string) (string, error) {
	fmt.Println("\nWriting transpilied C++ code to " + output + ".cpp ...")

	var (
		start = time.Now()
		// Write the C++ code to a file named `main.cpp`
		err = ioutil.WriteFile(output, []byte(source), 0644)
	)

	if err != nil {
		return "", err
	}
	timeTrack(start, "write")

	fmt.Println("\nFormatting C++ code ...")

	// TODO: later on format before writing to save the reading
	// Format the file in-place using `clang-format`; mainly for human readability
	start = time.Now()
	// TODO: pump this into clang later so that the errors that come back are formatted
	// for now we'll just return the source
	outputB, err := exec.Command("clang-format", "-i", output).CombinedOutput()
	if err != nil {
		return "", err
	}
	timeTrack(start, "format")

	return string(outputB), nil
}

func generateBinary(source, outputName string) error {
	// Track the time
	defer timeTrack(time.Now(), "clang")

	fmt.Println("\nUsing Clang generate create binary ...")

	// Compile the file with Clang to produce a binary
	compilerFlags = append(compilerFlags, outputName+".cpp", "-o", outputName, libmill)

	fmt.Printf("Using command: `clang++ %s`\n", strings.Join(compilerFlags, " "))
	// os.Exit(9)
	var clangCmd = exec.Command("clang++", compilerFlags...)
	// os.Exit(9)

	// // Grab the stdin of the command
	// var stdin, err = clangCmd.StdinPipe()
	// if err != nil {
	// 	return err
	// }

	// // Copy the bytes to Clang's stdin
	// n, err := copyToPipe(stdin, bytes.NewBufferString(source))
	// if err != nil {
	// 	return err
	// }

	// // Check that the amount copied is the amount we are expecting
	// if n != int64(len(source)) {
	// 	return errors.Errorf("Could not write all (%d) source bytes to clang: %d", len(source), n)
	// }

	// Start Clang to have it waiting
	output, err := clangCmd.CombinedOutput()
	if err != nil {
		fmt.Println("\nClang error:\n" + string(output))

		return err
	}

	return nil
}

// This function is really just to control the defer properly
func copyToPipe(in io.WriteCloser, out io.Reader) (int64, error) {
	// We need to ensure that the pipe is closed so that Clang will know that we are finished
	defer in.Close()

	// Whether we error or not we need to close the pipe
	return io.Copy(in, out)
}

func Compile(filename string) error {
	if !strings.HasSuffix(filename, ".expr") {
		return errors.Errorf("File does not have `.expr` suffix: %s", filename)
	}

	fmt.Println("\nReading input file ...")

	var (
		globalStart    = time.Now()
		start          = time.Now()
		testBytes, err = ioutil.ReadFile(filename)
	)

	if err != nil {
		return err
	}

	pipelineTimes["read"] = time.Since(start).String()

	fmt.Println("\nBuilding AST ...")

	// Build the AST
	start = time.Now()
	tr, err := getTranspilerFromString(string(testBytes), "main")
	if err != nil {
		return err
	}

	pipelineTimes["build"] = time.Since(start).String()

	fmt.Println("\nTranspiling to C++ ...")

	start = time.Now()
	cpp, err := tr.Transpile()
	if err != nil {
		return err
	}

	pipelineTimes["transpile"] = time.Since(start).String()

	if len(tr.Includes["libmill.h"]) > 0 {
		// TODO: fix this later
		libmill = "/usr/local/lib/libmill.a"
	}

	var wg sync.WaitGroup

	var rawFilename = strings.TrimSuffix(filename, ".expr")

	wg.Add(1)
	go func() {
		defer wg.Done()
		result, err := writeAndFormat(cpp, rawFilename+".cpp")
		if err != nil {
			fmt.Printf("There was an error writing C++ file; this does NOT inherently effect binary generation: %s : %+v\n", result, err)
		}
	}()

	err = generateBinary(cpp, rawFilename)
	if err != nil {
		return err
	}

	fmt.Println("\nFinished!")

	// Wait for the write/formatter to finish
	wg.Wait()

	pipelineTimes["compile"] = time.Since(globalStart).String()

	var times, _ = json.MarshalIndent(pipelineTimes, "", "  ")
	log.Println("\nPipeline timings:", string(times))

	return nil
}

func Run(filename string) error {
	var err = Compile(filename)
	if err != nil {
		return err
	}

	fmt.Println("\nRunning binary ...")

	// Run the produced binary
	var rawFilename = strings.TrimSuffix(filename, ".expr")
	output, err := exec.Command(rawFilename).Output()
	if err != nil {
		return err
	}

	fmt.Println("\nDone!")

	fmt.Println("\nOutput:", output)

	return nil
}
