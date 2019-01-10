package compiler

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
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

type Compiler struct {
	Raw     string
	LibBase string
	// Libmill       string
	PipelineTimes map[string]string
	Flags         []string
	path          string
	Outputs       map[string]string
	OutputData    map[string][]byte
}

func (c *Compiler) SetOutput(o map[string]string) {
	if o != nil {
		c.Outputs = o
	}
}

// New creates a compiler with default flags, base lib and others
func New(output string) (*Compiler, error) {
	var libpath = os.Getenv("EXPRPATH")

	if libpath == "" {
		return nil, errors.New("`EXPRPATH` is not set; set this to the root of your Express installation")
	}

	var _, err = filepath.Abs(libpath)
	if err != nil {
		return nil, err
	}

	return &Compiler{
		path:       output,
		OutputData: map[string][]byte{},
		Outputs:    map[string]string{},
		LibBase:    libpath + "/lib/",
		// LibBase:       "/home/scottshotgg/Development/go/src/github.com/scottshotgg/express2/lib/",
		PipelineTimes: map[string]string{},
		Flags: []string{
			stdCppVersion,
			"-Ofast",
			// "-x",
			// "c++",
		},
	}, nil
}

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

	return transpiler.New(ast, b, name, "idk"), nil
}

func (c *Compiler) timeTrack(start time.Time, name string) {
	// fmt.Printf("Function %s took %s\n", name, time.Since(start))
	c.PipelineTimes[name] = time.Since(start).String()
}

func (c *Compiler) writeAndFormat(source, output string) (string, error) {
	fmt.Println("\nWriting transpilied C++ code to " + output + " ...")

	var (
		start = time.Now()
		// Write the C++ code to a file named `main.cpp`
		err = ioutil.WriteFile(output, []byte(source), 0644)
	)

	if err != nil {
		return "", err
	}
	c.timeTrack(start, "write")

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
	c.timeTrack(start, "format")

	return string(outputB), nil
}

func (c *Compiler) generateBinary(source, outputName string) error {
	// Track the time
	defer c.timeTrack(time.Now(), "clang")

	fmt.Println("\nUsing Clang generate create binary ...")

	// TODO: its fine to use the local lib, we should do that, but we need to make an install script that will compile it
	// Compile the file with Clang to produce a binary
	c.Flags = append(c.Flags, outputName+".cpp", "-o", outputName, c.LibBase+"libmill/.libs/libmill.a")

	fmt.Printf("Using command: `clang++ %s`\n", strings.Join(c.Flags, " "))
	// os.Exit(9)
	var clangCmd = exec.Command("clang++", c.Flags...)
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

func (c *Compiler) setOutput(name string, output interface{}) error {
	// For now just disregard the error, also this may bite us in the ass in
	// the future. Might have to make some sort of encoding streamer to
	// save on memory usage vs storing everything
	var outputJSON, err = json.Marshal(output)
	if err != nil {
		return err
	}

	c.OutputData[name] = outputJSON

	return nil
}

func (c *Compiler) ProduceOutput(raw string) error {
	// C++ will not be in this
	var (
		err  error
		data []byte
		ok   bool
	)

	for t, f := range c.Outputs {
		data, ok = c.OutputData[t]
		if !ok {
			// TODO: handle this later
			continue
		}

		fmt.Println("writing file", f)
		err = ioutil.WriteFile(f, data, 0666)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Compiler) CompileFile(filename string) error {
	var (
		globalStart = time.Now()
		err         = c.compileFile(filename)
	)

	c.PipelineTimes["compile"] = time.Since(globalStart).String()

	if err != nil {
		return err
	}

	var times, _ = json.MarshalIndent(c.PipelineTimes, "", "  ")
	fmt.Println("\nPipeline timings:", string(times))

	var (
		rawPath      = strings.Trim(filename, ".expr")
		rawPathSplit = strings.Split(rawPath, "/")
	)

	return c.ProduceOutput(rawPathSplit[len(rawPathSplit)-1])
}

func (c *Compiler) compileFile(filename string) error {
	if !strings.HasSuffix(filename, ".expr") {
		return errors.Errorf("File does not have `.expr` suffix: %s", filename)
	}

	fmt.Println("\nReading input file ...")

	var (
		start       = time.Now()
		source, err = ioutil.ReadFile(filename)
		rawFilename = strings.TrimSuffix(filename, ".expr")
		wg          sync.WaitGroup
	)

	if err != nil {
		return err
	}

	c.PipelineTimes["read"] = time.Since(start).String()

	fmt.Println("\nTokenizing source ...")

	// Lex and tokenize the source code
	tokens, err := lex.New(string(source)).Lex()
	if err != nil {
		return err
	}
	c.setOutput("lex", tokens)

	fmt.Println("\nCompressing tokens ...")

	// Compress certain tokens;
	// i.e: `:` and `=` compress into `:=`
	tokens, err = ast.CompressTokens(tokens)
	if err != nil {
		return err
	}
	c.setOutput("compress", tokens)

	fmt.Println("\nBuilding AST ...")

	// Build the AST
	start = time.Now()
	var b = builder.New(tokens)
	ast, err := b.BuildAST()
	c.PipelineTimes["build"] = time.Since(start).String()
	if err != nil {
		return err
	}
	c.setOutput("ast", ast)

	// Change "main" to something else later

	fmt.Println("\nTranspiling to C++ ...")
	start = time.Now()
	var tr = transpiler.New(ast, b, "main", c.LibBase)
	cpp, err := tr.Transpile()
	c.PipelineTimes["transpile"] = time.Since(start).String()
	if err != nil {
		return err
	}
	// TODO: fix this ... :*(
	// c.OutputData["cpp"] = []byte(cpp)

	wg.Add(1)
	go func() {
		defer wg.Done()
		var result, err = c.writeAndFormat(cpp, rawFilename+".cpp")
		if err != nil {
			fmt.Printf("There was an error writing C++ file; this does NOT inherently effect binary generation: %s : %+v\n", result, err)
		}
	}()

	err = c.generateBinary(cpp, rawFilename)
	if err != nil {
		return err
	}

	// Wait for the write/formatter to finish
	wg.Wait()

	fmt.Println("\nFinished!")

	return nil
}

func (c *Compiler) RunFile(filename string) error {
	var err = c.CompileFile(filename)
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
