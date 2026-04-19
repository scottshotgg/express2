package compiler

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	ast "github.com/scottshotgg/express-ast"
	lex "github.com/scottshotgg/express-lex"
	"github.com/scottshotgg/express2/builder"
	"github.com/scottshotgg/express2/pkg/logger"
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
	// TODO: c.path is set but the binary always lands at strings.TrimSuffix(inputFile, ".expr").
	// Wire this up properly when the output flag is fully implemented.
	path       string
	Outputs    map[string]string
	OutputData map[string][]byte
	log        logger.Logger
}

func (c *Compiler) SetOutput(o map[string]string) {
	if o != nil {
		c.Outputs = o
	}
}

// New creates a compiler with default flags, base lib and others.
// An optional Logger may be provided; if omitted a no-op logger is used.
func New(output string, log ...logger.Logger) (*Compiler, error) {
	var libpath = os.Getenv("EXPRPATH")

	if libpath == "" {
		return nil, errors.New("`EXPRPATH` is not set; set this to the root of your Express installation")
	}

	var err error
	libpath, err = filepath.Abs(libpath)
	if err != nil {
		return nil, err
	}

	var l logger.Logger = logger.Noop()
	if len(log) > 0 && log[0] != nil {
		l = log[0]
	}

	return &Compiler{
		path:          output,
		OutputData:    map[string][]byte{},
		Outputs:       map[string]string{},
		LibBase:       libpath + "/lib/",
		PipelineTimes: map[string]string{},
		Flags: []string{
			stdCppVersion,
			"-Ofast",
			"-D_GNU_SOURCE",
		},
		log: l,
	}, nil
}

func (c *Compiler) timeTrack(start time.Time, name string) {
	c.PipelineTimes[name] = time.Since(start).String()
}

func (c *Compiler) writeAndFormat(source, output string) (string, error) {
	c.log.Debug("Writing transpiled C++ code to " + output + " ...")

	var (
		start = time.Now()
		err   = ioutil.WriteFile(output, []byte(source), 0644)
	)

	if err != nil {
		return "", err
	}
	c.timeTrack(start, "write")

	c.log.Debug("Formatting C++ code ...")

	start = time.Now()
	outputB, err := exec.Command("clang-format", "-i", output).CombinedOutput()
	if err != nil {
		return "", err
	}
	c.timeTrack(start, "format")

	return string(outputB), nil
}

func (c *Compiler) generateBinary(outputName string) error {
	defer c.timeTrack(time.Now(), "clang")

	c.log.Debug("Using Clang to create binary ...")

	c.Flags = append(c.Flags, outputName+".cpp", "-o", outputName, c.LibBase+"libmill/.libs/libmill.a")

	c.log.Debugf("Using command: `clang++ %s`", strings.Join(c.Flags, " "))
	var clangCmd = exec.Command("clang++", c.Flags...)

	output, err := clangCmd.CombinedOutput()
	if err != nil {
		c.log.Errorf("Clang error:\n%s", string(output))
		return err
	}

	return nil
}

func (c *Compiler) setOutput(name string, output interface{}) error {
	var outputJSON, err = json.Marshal(output)
	if err != nil {
		return err
	}

	c.OutputData[name] = outputJSON

	return nil
}

func (c *Compiler) ProduceOutput(raw string) error {
	var (
		err  error
		data []byte
		ok   bool
	)

	for t, f := range c.Outputs {
		data, ok = c.OutputData[t]
		if !ok {
			continue
		}

		c.log.Debug("writing file", f)
		err = ioutil.WriteFile(f, data, 0666)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Compiler) CompileFile(filename string) error {
	defer func() {
		var times, _ = json.MarshalIndent(c.PipelineTimes, "", "  ")
		c.log.Debug("Pipeline timings:", string(times))
	}()

	var (
		globalStart = time.Now()
		err         = c.compileFile(filename)
	)

	c.PipelineTimes["compile"] = time.Since(globalStart).String()
	if err != nil {
		return err
	}

	var (
		rawPath      = strings.TrimSuffix(filename, ".expr")
		rawPathSplit = strings.Split(rawPath, "/")
	)

	return c.ProduceOutput(rawPathSplit[len(rawPathSplit)-1])
}

func (c *Compiler) compileFile(filename string) error {
	if !strings.HasSuffix(filename, ".expr") {
		return errors.Errorf("File does not have `.expr` suffix: %s", filename)
	}

	c.log.Debug("Reading input file ...")

	var rawFilename = strings.TrimSuffix(filename, ".expr")

	c.log.Debug("Tokenizing source ...")

	start := time.Now()
	l, err := lex.NewFromFile(filename)
	if err != nil {
		return err
	}

	tokens, err := l.Lex()
	if err != nil {
		return err
	}
	c.PipelineTimes["read"] = time.Since(start).String()

	c.setOutput("lex", tokens)

	c.log.Debug("Compressing tokens ...")

	tokens, err = ast.CompressTokens(tokens)
	if err != nil {
		return err
	}
	c.setOutput("compress", tokens)

	c.log.Debug("Building AST ...")

	start = time.Now()
	var b = builder.New(tokens, c.log)
	astNode, err := b.BuildAST()
	c.PipelineTimes["build"] = time.Since(start).String()
	if err != nil {
		return err
	}
	c.setOutput("ast", astNode)

	astJSON, err := json.Marshal(astNode)
	if err != nil {
		return err
	}
	c.log.Debug(string(astJSON))

	c.log.Debug("Running semantic pass ...")
	start = time.Now()
	astNode, err = builder.NewChecker(astNode, builder.NewTypeResolverWithScope(b.ScopeTree)).Execute()
	c.PipelineTimes["semantic"] = time.Since(start).String()
	if err != nil {
		return err
	}

	c.log.Debug("Transpiling to C++ ...")
	start = time.Now()
	var tr = transpiler.New(astNode, b, "main", c.LibBase)

	err = tr.Transpile()
	c.PipelineTimes["transpile"] = time.Since(start).String()
	if err != nil {
		return err
	}

	var cpp = tr.ToCpp()

	result, err := c.writeAndFormat(cpp, rawFilename+".cpp")
	if err != nil {
		c.log.Warnf("Error writing C++ file; this does NOT inherently affect binary generation: %s : %+v", result, err)
	}

	// Only link a binary if there is a main function (library files have none)
	if tr.GenerateMain {
		err = c.generateBinary(rawFilename)
		if err != nil {
			return err
		}
	}

	c.log.Debug("Finished!")

	return nil
}

func (c *Compiler) RunFile(filename string) error {
	var err = c.CompileFile(filename)
	if err != nil {
		return err
	}

	c.log.Debug("Running binary ...")

	var rawFilename = strings.TrimSuffix(filename, ".expr")
	output, err := exec.Command(rawFilename).Output()
	if err != nil {
		return err
	}

	c.log.Debug("Done!")
	c.log.Debug("Output:", string(output))

	return nil
}
