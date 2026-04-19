package transpiler

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/pkg/errors"
	token "github.com/scottshotgg/express-token"
	"github.com/scottshotgg/express2/builder"
	"github.com/scottshotgg/express2/pkg/logger"
	"github.com/scottshotgg/express2/tree_flattener"
)

type WithPriority struct {
	Priority int
	Value    string
}

type Transpiler struct {
	LibBase         string
	Name            string
	Builder         *builder.Builder
	AST             *builder.Node
	ASTCloneJSON    []byte
	Extra           []string
	Functions       map[string]string
	Imports         map[string]string
	Includes        map[string]string
	Packages        map[string]string
	Types           map[string]string
	Structs         map[string]WithPriority
	Methods         map[string][]string // receiver → list of transpiled method strings
	CurrentReceiver string      // set while transpiling a method body; protected by funcMu
	funcMu          sync.Mutex // serializes TranspileFunctionStatement across concurrent workers
	GenerateMain    bool
	workerErr       error
	workerErrOnce   sync.Once
	includesMu      sync.Mutex // protects concurrent writes to Includes
	Wg              *sync.WaitGroup // importWorker
	Wg1             *sync.WaitGroup // functionWorker
	Wg2             *sync.WaitGroup // packageWorker, typeWorker, structWorker
	FuncChan        chan *builder.Node
	TypeChan        chan *builder.Node
	StructChan      chan *builder.Node
	PackageChan     chan *builder.Node
	ImportChan      chan *builder.Node
	log             logger.Logger
}


func (t *Transpiler) functionWorker() {
	defer t.Wg1.Done()

	var (
		function     *string
		err          error
		functionName string
	)

	for f := range t.FuncChan {
		functionName = f.Kind

		// Hold funcMu for the entire transpilation of this function to prevent
		// data races on CurrentReceiver and inParamContext with packageWorker.
		t.funcMu.Lock()

		// Method declaration: func Receiver.Method() — transpile into t.Methods
		if receiver, ok := f.Metadata["receiver"].(string); ok && receiver != "" {
			t.CurrentReceiver = receiver
			function, err = t.TranspileFunctionStatement(f)
			t.CurrentReceiver = ""
			t.funcMu.Unlock()
			if err != nil {
				t.setWorkerErr(errors.Wrapf(err, "transpiling method %q.%q", receiver, functionName))
				continue
			}
			t.Methods[receiver] = append(t.Methods[receiver], *function)
			continue
		}

		if t.Functions[functionName] != "" {
			t.funcMu.Unlock()
			t.setWorkerErr(errors.Errorf("function already declared: %s", functionName))
			continue
		}

		function, err = t.TranspileFunctionStatement(f)
		t.funcMu.Unlock()
		if err != nil {
			t.setWorkerErr(errors.Wrapf(err, "transpiling function %q", functionName))
			continue
		}

		t.Functions[functionName] = *function
	}
}

func New(ast *builder.Node, b *builder.Builder, name, libBase string, log ...logger.Logger) *Transpiler {
	var l logger.Logger = logger.Noop()
	if len(log) > 0 && log[0] != nil {
		l = log[0]
	}

	var t = Transpiler{
		LibBase:     libBase,
		Name:        name,
		AST:         ast,
		Builder:     b,
		Functions:   map[string]string{},
		Imports:     map[string]string{},
		Packages:    map[string]string{},
		Includes:    map[string]string{},
		Structs:     map[string]WithPriority{},
		Types:       map[string]string{},
		Methods:     map[string][]string{},
		FuncChan:    make(chan *builder.Node, 100),
		TypeChan:    make(chan *builder.Node, 100),
		StructChan:  make(chan *builder.Node, 100),
		PackageChan: make(chan *builder.Node, 100),
		ImportChan:  make(chan *builder.Node, 100),
		Wg:          &sync.WaitGroup{},
		Wg1:         &sync.WaitGroup{},
		Wg2:         &sync.WaitGroup{},
		log:         l,
	}

	// go appendWorker(&wg)

	t.ASTCloneJSON, _ = json.Marshal(ast)

	return &t
}

func (t *Transpiler) setWorkerErr(err error) {
	t.workerErrOnce.Do(func() { t.workerErr = err })
}

/*
	Transpile needs to work like this:
	- recurse through each statement
	- if the statement contains ANY block, then flatten on the node
*/

// TODO: rewrite this
func (t *Transpiler) Transpile() error {
	// Extract the nodes
	var (
		// flattenedImports []*builder.Node
		nodes = t.AST.Value.([]*builder.Node)
		// stringP *string
		err error

		// cpp string
	)

	// Spin off workers for each type of statement

	t.Wg1.Add(1)
	go t.functionWorker()

	t.Wg2.Add(1)
	go t.typeWorker()

	t.Wg2.Add(1)
	go t.structWorker()

	t.Wg2.Add(1)
	go t.packageWorker()

	t.Wg.Add(1)
	go t.importWorker()

	// Flatten the tree — include nodes from the flattener are registered directly.
	includes, err := tree_flattener.New().Flatten(t.AST)
	if err != nil {
		return err
	}

	for i := range includes {
		s, serr := t.TranspileIncludeStatement(includes[i])
		if serr != nil {
			return errors.Wrapf(serr, "transpiling include from tree_flattener")
		}
		t.includesMu.Lock()
		t.Includes[includes[i].Left.Value.(string)] = *s
		t.includesMu.Unlock()
	}

	for i := range nodes {
		blob, _ := json.Marshal(nodes[i])
		t.log.Debug("blob:", string(blob))
		// TODO: Switch on the statement type to figure out how to process it
		// TODO: Flatten anything with a scope

		// TODO: need to put the function into the function chan here?

		switch nodes[i].Type {
		case "function":
			t.FuncChan <- nodes[i]

		case "struct":
			t.StructChan <- nodes[i]

		case "typedef":
			t.TypeChan <- nodes[i]

		case "import":
			t.ImportChan <- nodes[i]

		case "include":
			s, serr := t.TranspileIncludeStatement(nodes[i])
			if serr != nil {
				return errors.Wrapf(serr, "transpiling include node")
			}
			t.includesMu.Lock()
			t.Includes[nodes[i].Left.Value.(string)] = *s
			t.includesMu.Unlock()

		case "decl":
			// Just transpile the statement for now
			stringP, err := t.TranspileStatement(nodes[i])
			if err != nil {
				return errors.Wrapf(err, "transpiling %s node", nodes[i].Type)
			}

			t.Extra = append(t.Extra, *stringP)

		case "use":
			// Just transpile the statement for now
			stringP, err := t.TranspileStatement(nodes[i])
			if err != nil {
				return errors.Wrapf(err, "transpiling %s node", nodes[i].Type)
			}

			t.Extra = append(t.Extra, *stringP)

		case "map":
			// Just transpile the statement for now
			stringP, err := t.TranspileStatement(nodes[i])
			if err != nil {
				return errors.Wrapf(err, "transpiling %s node", nodes[i].Type)
			}

			t.Extra = append(t.Extra, *stringP)

		case "package":
			t.PackageChan <- nodes[i]

		default:
			return errors.Errorf("Node was not categorized properly: %+v\n", nodes[i])
		}

	}

	// Close producer channels (package/type/struct workers may send to FuncChan).
	close(t.PackageChan)
	close(t.TypeChan)
	close(t.StructChan)

	// Wait for Wg2 workers (package/type/struct) to finish before closing FuncChan,
	// since packageWorker may call TranspileStatement which sends to FuncChan.
	t.Wg2.Wait()
	close(t.FuncChan)

	// Wait for functionWorker to drain FuncChan.
	t.Wg1.Wait()

	close(t.ImportChan)

	// Wait for the import worker to finish
	t.Wg.Wait()

	if t.Functions["main"] != "" {
		t.GenerateMain = true
	}

	if t.workerErr != nil {
		return t.workerErr
	}

	return nil
}

func (t *Transpiler) typeWorker() {
	defer t.Wg2.Done()

	var (
		stringP *string
		err     error
	)

	for node := range t.TypeChan {
		stringP, err = t.TranspileStatement(node)
		if err != nil {
			t.setWorkerErr(errors.Wrapf(err, "transpiling type %q", node.Left.Value.(string)))
			continue
		}

		t.Types[node.Left.Value.(string)] = *stringP
	}
}

func (t *Transpiler) structWorker() {
	defer t.Wg2.Done()

	var (
		stringP *string
		err     error
	)

	var i int
	for node := range t.StructChan {
		stringP, err = t.TranspileStatement(node)
		if err != nil {
			t.setWorkerErr(errors.Wrapf(err, "transpiling struct %q", node.Left.Value.(string)))
			continue
		}

		t.Structs[node.Left.Value.(string)] = WithPriority{
			Priority: i,
			Value:    *stringP,
		}

		i++
	}
}

func (t *Transpiler) packageWorker() {
	defer t.Wg2.Done()

	var (
		stringP *string
		err     error
	)

	var i int
	for node := range t.PackageChan {
		stringP, err = t.TranspileStatement(node)
		if err != nil {
			t.setWorkerErr(errors.Wrapf(err, "transpiling package %q", node.Left.Value.(string)))
			continue
		}

		t.Packages[node.Left.Value.(string)] = *stringP

		i++
	}
}

func (t *Transpiler) importWorker() {
	defer t.Wg.Done()

	var (
		importStringP *string
		// Why does this shadow ...
		// Is the gofunc "capturing" variables that aren't passed?
		ierr error
	)

	for node := range t.ImportChan {
		// Might want to make this go through the entire pipeline ...
		importStringP, ierr = t.TranspileImportStatement(node)
		if ierr != nil {
			t.setWorkerErr(errors.Wrapf(ierr, "transpiling import %q", node.Left.Value.(string)))
			continue
		}

		if importStringP == nil {
			continue
		}

		// TODO: should really check the deref on all of these, but the usage/running
		// is pretty predictable right now
		t.Imports[node.Left.Value.(string)] = *importStringP
	}
}

func (t *Transpiler) ToCpp() string {
	// Put the main functions before the other cpp code; I don't think
	// there should be anything in cpp, but w/e

	// TODO: need to printout // Global stuff and print out the global statements

	var output []string

	output = append(output, "// Includes:")
	if len(t.Imports) > 0 {
		// output = append(output, strings.Join(t.Includes, "\n")+"\n")
		var importString string
		for _, t := range t.Imports {
			importString += t + "\n"
		}
		output = append(output, importString)
	} else {
		output = append(output, "// none\n")
	}

	output = append(output, "// Imports:")
	if len(t.Includes) > 0 {
		// output = append(output, strings.Join(t.Imports, "\n")+"\n")
		var (
			includeString string
			libmill       string
		)
		for _, t := range t.Includes {
			if strings.Contains(t, "libmill") {
				libmill = t
				continue
			}

			includeString += t + "\n"
		}

		output = append(output, includeString)

		if len(libmill) > 0 {
			output = append(output, libmill)
		}

	} else {
		output = append(output, "// none\n")
	}

	output = append(output, "// Namespaces:")
	if len(t.Packages) > 0 {
		// output = append(output, strings.Join(t.Includes, "\n")+"\n")
		var packageString string
		for _, t := range t.Packages {
			packageString += t + "\n"
		}
		output = append(output, packageString)
	} else {
		output = append(output, "// none\n")
	}

	// Save the main function to separate it from the rest of functions since it has a
	// special purpose.

	return strings.Join(append(output, []string{
		t.generateTypes(),
		strings.Join(t.Extra, "\n"),
		t.generateFunctions(),
	}...), "\n")
}

func (t *Transpiler) generateTypes() string {
	var typesString = "\n\n// Types:\n"
	for _, t := range t.Types {
		typesString += t + "\n"
	}

	if len(typesString) == len("\n\n// Types:\n") {
		typesString += "// none\n"
	}

	var structs = make([]string, len(t.Structs))
	for name, s := range t.Structs {
		value := s.Value
		// Inject methods into struct definition (C++ requires methods inside struct)
		if methods, ok := t.Methods[name]; ok && len(methods) > 0 {
			methodsStr := "\n" + strings.Join(methods, "\n") + "\n"
			// Insert before the closing };
			if idx := strings.LastIndex(value, "};"); idx >= 0 {
				value = value[:idx] + methodsStr + "};"
			}
		}
		structs[s.Priority] = value
	}

	var structsString = "\n\n// Structs:\n"
	if len(structs) == len("\n\n// Structs:\n") {
		structsString += "// none\n"
	} else {
		structsString += strings.Join(structs, "\n")
	}

	return typesString + structsString
}

func (t *Transpiler) generateFunctions() string {
	var (
		prototypes     []string
		functionString string
	)

	t.log.Debug("mainFunc", t.Functions["main"])
	var mainFunc = t.Functions["main"]
	delete(t.Functions, "main")

	// Put the functions at the top of the file before the main function
	for _, f := range t.Functions {
		// TODO: just hack this in here for now to make the function prototypes
		prototypes = append(prototypes, strings.Split(f, "{")[0]+";")
		functionString += "\n" + f + "\n"
	}

	var fullString = "\n// Prototypes:\n"
	if len(prototypes) == 0 {
		fullString += "// none"
	} else {
		fullString += strings.Join(prototypes, "\n")
	}

	fullString += "\n\n// Functions:"
	if len(functionString) == 0 {
		fullString += "// none"
	} else {
		fullString += functionString
	}

	return fullString +
		fmt.Sprintf("\n// Main:\n// generated: %v\n%s", t.GenerateMain, mainFunc)
}

func (t *Transpiler) TranspileTypeDeclaration(n *builder.Node) (*string, error) {
	// Format should be:
	// `type` [ident] `=` [type]
	// Left is the ident
	// Right is the type

	// TODO: gonna have to do something to actually enable this type in the parser/compiler

	if n.Type != "typedef" {
		return nil, errors.New("Node is not a typedef")
	}

	var nString = "typedef "

	var cpp, err = t.TranspileType(n.Right)
	if err != nil {
		return nil, err
	}

	nString += *cpp + " "

	// This will allow technically allow idents to be made from general expressions; not sure if we should keep this or not
	// Might have to change it to TranspileIdent
	cpp, err = t.TranspileExpression(n.Left)
	if err != nil {
		return nil, err
	}

	nString += *cpp + ";"

	return &nString, nil
}

func (t *Transpiler) TranspileObjectStatement(n *builder.Node) (*string, error) {
	/*
		This should transpile to:
		object something = {} : class something {}
		Type is class?
		Left is the ident
		Right is the value
	*/

	if n.Type != "object" {
		t.log.Debugf("n %+v", n)
		t.log.Debugf("n %+v", n.Left)
		t.log.Debugf("n %+v", n.Right)
		return nil, errors.New("Node is not a object")
	}

	// Transpile the ident which will become a usable type
	var vString, err = t.TranspileExpression(n.Left)
	if err != nil {
		return nil, err
	}

	// Could just have it add `object` here but this will show us changes
	var nString = "class " + *vString

	// Transpile the block for the value
	vString, err = t.TranspileBlockStatement(n.Right)
	if err != nil {
		return nil, err
	}

	nString += *vString + ";"

	return &nString, nil
}

func (t *Transpiler) TranspileStructDeclaration(n *builder.Node) (*string, error) {
	/*
		This should transpile to:
		struct something = {} : struct something {}
		Type is struct
		Left is the ident
		Right is the value
	*/

	if n.Type != "struct" {
		return nil, errors.New("Node is not a struct")
	}

	// Transpile the ident which will become a usable type
	var vString, err = t.TranspileExpression(n.Left)
	if err != nil {
		return nil, err
	}

	// Could just have it add `struct` here but this will show us changes
	var nString = n.Type + " " + *vString

	// Transpile the block for the value
	vString, err = t.TranspileBlockStatement(n.Right)
	if err != nil {
		return nil, err
	}

	nString += *vString + ";"

	return &nString, nil
}

func (t *Transpiler) TranspileIncludeStatement(n *builder.Node) (*string, error) {
	if n.Type != "include" {
		return nil, errors.New("Node is not an include")
	}

	lhs, err := t.TranspileExpression(n.Left)
	if err != nil {
		return nil, err
	}

	// Deal with the user defined shit in the semantic stage
	if n.Kind == "path" {
		abs, err := filepath.Abs(*lhs)
		if err != nil {
			return nil, err
		}

		*lhs = "#include \"" + abs + "\""
	} else {
		*lhs = "#include<" + *lhs + ">"
	}

	return lhs, nil
}

// requireStdInclude registers a system header include, e.g. requireStdInclude("string")
// emits #include<string>. Idempotent and concurrent-safe.
func (t *Transpiler) requireStdInclude(header string) {
	t.includesMu.Lock()
	t.Includes[header] = "#include<" + header + ">"
	t.includesMu.Unlock()
}

// requirePathInclude registers a path-based include, e.g. requirePathInclude("/path/to/var.cpp")
// emits #include "/abs/path/to/var.cpp". Idempotent and concurrent-safe.
func (t *Transpiler) requirePathInclude(path string) {
	abs, _ := filepath.Abs(path)
	t.includesMu.Lock()
	t.Includes[path] = `#include "` + abs + `"`
	t.includesMu.Unlock()
}

// requireType maps an Express type name to its C++ equivalent and registers any
// required includes. Returns the C++ type string.
func (t *Transpiler) requireType(exprType string) string {
	switch exprType {
	case "string":
		t.requireStdInclude("string")
		return "std::string"
	case "int", "float", "bool", "char":
		return exprType
	case "var":
		t.requirePathInclude(t.LibBase + "var.cpp")
		return "var"
	default:
		return exprType // user-defined struct types pass through unchanged
	}
}

func (t *Transpiler) TranspileUseStatement(n *builder.Node) (*string, error) {
	if n.Type != "use" {
		return nil, errors.New("Node is not a use")
	}

	/*
		Left is the "imported" file/package
		Right is the new ident
	*/

	lhs, err := t.TranspileExpression(n.Left)
	if err != nil {
		return nil, err
	}

	rhs, err := t.TranspileExpression(n.Right)
	if err != nil {
		return nil, err
	}

	// Ignore the `unused` error for now, we'll fix it later
	_ = lhs
	_ = rhs

	// Imports should not have angled brackets
	// This is really more of a _semantic_ or even _parser_ thing to go grab the code
	// Or to link the object as a shared resource into the binary
	// *lhs = "#include " + *lhs

	return nil, errors.New("`use` statements are currently not available; their implementation is currently waiting on the semantic stage and more improvements to the parser")
}

// TODO: this should transpile the program and take all statements
// and put them into the file under a namespace that is the package name
func (t *Transpiler) TranspileImportStatement(n *builder.Node) (*string, error) {
	if n.Type != "import" {
		return nil, errors.New("Node is not an import")
	}

	// If it is the import for libc — register the standard C headers.
	if n.Kind == "c" {
		for _, hdr := range []string{"stdio.h", "stdlib.h", "string.h", "unistd.h", "libgen.h", "math.h"} {
			t.requireStdInclude(hdr)
		}
		return nil, nil
	}

	// Make a new namespace in the file
	// Transpile that AST into the namespace

	// this should be done in the semantic stage
	switch n.Left.Type {
	case "ident":
		// Need to add quotes and make sure that the library exists

	case "literal":
		// Check the literal obvi brah

		// var tr = New(n.Right, nil, "main", t.LibBase)
		// fmt.Println("tr", tr)
		// cpp, err := tr.Transpile()
		// fmt.Println("cpp, err", cpp, err)
	}

	var tt = New(n.Right, t.Builder, n.Left.Value.(string), t.LibBase)

	err := tt.Transpile()
	if err != nil {
		return nil, err
	}

	for k, v := range tt.Includes {
		t.Includes[k] = v
	}

	tt.Includes = nil

	packageName := n.Left.Value.(string)
	// If the imported file had a `package` declaration it already produced
	// a namespace __name { ... } block via TranspilePackageStatement.
	// In that case, don't wrap again — just use the inner ToCpp() directly.
	if _, hasNamespace := tt.Packages[packageName]; hasNamespace {
		t.Packages[packageName] = tt.ToCpp()
	} else {
		t.Packages[packageName] = fmt.Sprintf("namespace __%s {\n %s\n }\n", packageName, tt.ToCpp())
	}
	// t.IncludeChan <- &builder.Node{
	// 	Type: "include",
	// 	Left: &builder.Node{
	// 		Type:  "literal",
	// 		Value: fmt.Sprintf("/home/scottshotgg/Development/go/src/github.com/scottshotgg/express2/%s.expr.cpp", n.Left.Value.(string)),
	// 	},
	// }

	// return &empty, ioutil.WriteFile(n.Left.Value.(string)+".expr.cpp", []byte(tt.ToCpp()), 0777)
	return &empty, nil
}

var empty = ""

func (t *Transpiler) TranspileIncrementExpression(n *builder.Node) (*string, error) {
	if n.Type != "inc" {
		return nil, errors.New("Node is not an inc")
	}

	lhs, err := t.TranspileExpression(n.Left)
	if err != nil {
		return nil, err
	}

	// Put parenthesis around it
	*lhs = "(" + *lhs + ")++"

	return lhs, nil
}

func (t *Transpiler) TranspileDecrementExpression(n *builder.Node) (*string, error) {
	if n.Type != "dec" {
		return nil, errors.New("Node is not a dec")
	}

	lhs, err := t.TranspileExpression(n.Left)
	if err != nil {
		return nil, err
	}

	// Put parenthesis around it
	*lhs = "(" + *lhs + ")--"

	return lhs, nil
}

func (t *Transpiler) TranspileIndexExpression(n *builder.Node) (*string, error) {
	/*
		Left is an expression
		Right is an expression
	*/

	if n.Type != "index" {
		return nil, errors.New("Node is not an index")
	}

	lhs, err := t.TranspileExpression(n.Left)
	if err != nil {
		return nil, err
	}

	rhs, err := t.TranspileExpression(n.Right)
	if err != nil {
		return nil, err
	}

	var nString = *lhs + "[" + *rhs + "]"

	return &nString, nil
}

// The type checker should produce arrow function ones as well
func (t *Transpiler) TranspileSelectExpression(n *builder.Node) (*string, error) {
	/*
		Left is an expression
		Right is an expression
	*/

	if n.Type != "selection" {
		return nil, errors.New("Node is not an index")
	}

	lhs, err := t.TranspileExpression(n.Left)
	if err != nil {
		return nil, err
	}

	rhs, err := t.TranspileExpression(n.Right)
	if err != nil {
		return nil, err
	}

	// c namespace: strip the c. prefix entirely.
	// After #include<stdio.h> etc., all C symbols are in global scope in C++.
	// c.fopen(...) → fopen(...)   c.SEEK_SET → SEEK_SET
	if n.Left.Type == "ident" && n.Left.Value.(string) == "c" {
		return rhs, nil
	}

	// Inside a method body: Receiver.field → field (C++ member access).
	// Left may be "ident" or "type" (since the receiver name is a registered type).
	if t.CurrentReceiver != "" {
		if name, ok := n.Left.Value.(string); ok && name == t.CurrentReceiver {
			return rhs, nil
		}
	}

	// Imported package namespace: file.X → __file::X
	if n.Left.Type == "ident" {
		if name, ok := n.Left.Value.(string); ok {
			if _, isPackage := t.Packages[name]; isPackage {
				result := "__" + name + "::" + *rhs
				return &result, nil
			}
		}
	}

	var selector = "."
	switch n.Left.Type {
	case "package":
		selector = "::"
	case "ident":
		// Look up the declared type of this variable in the scope tree.
		// If it was declared as a pointer type, use -> instead of .
		if decl := t.Builder.ScopeTree.Get(n.Left.Value.(string)); decl != nil {
			if typeNode, ok := decl.Value.(*builder.Node); ok {
				if typeNode.Kind == "pointer" {
					selector = "->"
				}
			}
		}
	}

	var nString = *lhs + selector + *rhs

	return &nString, nil
}

// The type checker should produce arrow function ones as well
func (t *Transpiler) TranspileDerefExpression(n *builder.Node) (*string, error) {
	if n.Type != "deref" {
		return nil, errors.New("Node is not an deref")
	}

	var nString = "*"

	// Left is the ident; right is nothing
	lhs, err := t.TranspileExpression(n.Left)
	if err != nil {
		return nil, err
	}

	if n.Kind == "type" {
		nString = *lhs + nString
	} else {
		nString += *lhs
	}

	return &nString, nil
}

func (t *Transpiler) TranspileNotExpression(n *builder.Node) (*string, error) {
	if n.Type != "not" {
		return nil, errors.New("Node is not a not")
	}

	operand, err := t.TranspileExpression(n.Left)
	if err != nil {
		return nil, err
	}

	result := "!" + *operand
	return &result, nil
}

// The type checker should produce arrow function ones as well
func (t *Transpiler) TranspileRefExpression(n *builder.Node) (*string, error) {
	if n.Type != "ref" {
		return nil, errors.New("Node is not an ref")
	}

	// Left is the ident; right is nothing
	lhs, err := t.TranspileExpression(n.Left)
	if err != nil {
		return nil, err
	}

	var nString = "&" + *lhs

	return &nString, nil
}

func (t *Transpiler) TranspileExpression(n *builder.Node) (*string, error) {
	switch n.Type {

	case "literal":
		return t.TranspileLiteralExpression(n)

	case "ident":
		return t.TranspileIdentExpression(n)

	case "comp":
		return t.TranspileConditionExpression(n)

	case "binop":
		return t.TranspileBinOpExpression(n)

	case "array":
		return t.TranspileArrayExpression(n)

	case "call":
		return t.TranspileCallExpression(n)

	case "index":
		return t.TranspileIndexExpression(n)

	case "block":
		return t.TranspileBlockExpression(n)

	case "selection":
		return t.TranspileSelectExpression(n)

	case "deref":
		return t.TranspileDerefExpression(n)

	case "ref":
		return t.TranspileRefExpression(n)

	case "not":
		return t.TranspileNotExpression(n)

	case "inc":
		return t.TranspileIncrementExpression(n)

	case "dec":
		return t.TranspileDecrementExpression(n)

	case "type":
		return t.TranspileType(n)

	case "package":
		return t.TranspilePackageExpression(n)
	}

	return nil, errors.Errorf("Not implemented expression: %+v", n)
}

func (t *Transpiler) TranspilePackageExpression(n *builder.Node) (*string, error) {
	var value = "__" + n.Value.(string)
	return &value, nil
}

// transpileMapTypeStr converts a type node (or nested map type node) to a C++ type string.
func (t *Transpiler) transpileMapTypeStr(n *builder.Node) (string, error) {
	if n.Kind == "map" {
		kn, ok1 := n.Metadata["key_node"].(*builder.Node)
		vn, ok2 := n.Metadata["value_node"].(*builder.Node)
		if !ok1 || !ok2 {
			return "", errors.New("transpileMapTypeStr: map node missing key_node/value_node")
		}
		k, err := t.transpileMapTypeStr(kn)
		if err != nil {
			return "", err
		}
		v, err := t.transpileMapTypeStr(vn)
		if err != nil {
			return "", err
		}
		t.requireStdInclude("map")
		return fmt.Sprintf("std::map<%s, %s>", k, v), nil
	}
	return t.requireType(n.Kind), nil
}

func (t *Transpiler) TranspileMapStatement(n *builder.Node) (*string, error) {
	if n.Type != "map" {
		return nil, errors.New("Node is not a map")
	}

	// Transpile the ident
	vString, err := t.TranspileExpression(n.Left)
	if err != nil {
		return nil, err
	}

	// Determine K/V types — use explicit annotation if present, else default.
	// requireType registers necessary includes automatically.
	keyType := t.requireType("string") // "std::string", registers <string>
	valueType := t.requireType("var")  // "var", registers var.cpp
	if n.Metadata != nil {
		if kn, ok := n.Metadata["key_node"].(*builder.Node); ok {
			var err2 error
			keyType, err2 = t.transpileMapTypeStr(kn)
			if err2 != nil {
				return nil, err2
			}
			valueType, err2 = t.transpileMapTypeStr(n.Metadata["value_node"].(*builder.Node))
			if err2 != nil {
				return nil, err2
			}
		}
	}

	// std::map is always required for map statements
	t.requireStdInclude("map")

	// Handle zero-init (no body)
	if n.Right == nil {
		nString := fmt.Sprintf("std::map<%s, %s> %s;", keyType, valueType, *vString)
		return &nString, nil
	}

	var nString = fmt.Sprintf("std::map<%s, %s> %s = ", keyType, valueType, *vString)

	// Transpile the block for the value
	vString, err = t.TranspileMapBlockStatement(n.Right)
	if err != nil {
		return nil, err
	}

	nString += *vString + ";"

	return &nString, nil
}

func (t *Transpiler) TranspileLaunchStatement(n *builder.Node) (*string, error) {
	if n.Type != "launch" {
		return nil, errors.New("Node is not a launch node")
	}

	blob, _ := json.Marshal(n)
	t.log.Debug("launch:", string(blob))

	blob, _ = json.Marshal(n)
	t.log.Debug("launch.left:", string(blob))

	// Transpile the ident
	var vString, err = t.TranspileStatement(n.Left)
	if err != nil {
		return nil, err
	}

	// Include libmill for coroutines
	t.requirePathInclude(t.LibBase + "libmill/libmill.h")

	// This has a lambda in it since you can launch any statement ...
	var nString = "go([=](...){" + *vString + "}());"

	return &nString, nil
}

func (t *Transpiler) TranspileEnumBlockStatement(n *builder.Node) (*string, error) {
	if n.Type != "block" {
		return nil, errors.New("Node is not a block")
	}

	var (
		nString = ""
		vString *string
		err     error
	)

	for _, stmt := range n.Value.([]*builder.Node) {
		if stmt.Type != "assignment" && stmt.Type != "ident" {
			return nil, errors.Errorf("All statements in an enum have to be assignment or ident: %+v\n", stmt)
		}

		vString, err = t.TranspileStatement(stmt)
		if err != nil {
			return nil, err
		}

		if stmt.Type == "assignment" {
			*vString = (*vString)[:len(*vString)-1]
		}

		nString += *vString + ","
	}

	nString = "{" + nString + "}"

	return &nString, nil
}

func (t *Transpiler) TranspileEnumStatement(n *builder.Node) (*string, error) {
	if n.Type != "enum" {
		return nil, errors.New("Node is not a map")
	}

	var enum, err = t.TranspileEnumBlockStatement(n.Left)
	if err != nil {
		return nil, err
	}

	var nString = "enum " + *enum + ";"

	return &nString, err
}

func (t *Transpiler) TranspileKeyValueStatement(n *builder.Node) (*string, error) {
	var left, err = t.TranspileExpression(n.Left)
	if err != nil {
		return nil, err
	}

	right, err := t.TranspileExpression(n.Right)
	if err != nil {
		return nil, err
	}

	var nString = "{" + *left + "," + *right + "}"

	return &nString, nil
}

func (t *Transpiler) TranspileDeferStatement(n *builder.Node) (*string, error) {
	var stmt, err = t.TranspileStatement(n.Left)
	if err != nil {
		return nil, err
	}

	// TODO: we need to wipe the defer stacks unless they are explicitly used
	//		[=] - value
	//		[&] - reference

	// TODO: only onReturn is supported for now
	var nString = "onReturn.deferStack.push([&](...){" + *stmt + "});"

	return &nString, nil
}

func (t *Transpiler) TranspileStatement(n *builder.Node) (*string, error) {
	t.log.Debug("wtf3333", n.Type)
	switch n.Type {

	case "c":
		// Direct C/C++ code injection — emit the raw captured source verbatim.
		nString := n.Value.(string)
		return &nString, nil

	case "break":
		s := "break;"
		return &s, nil

	case "continue":
		s := "continue;"
		return &s, nil

	case "if":
		return t.TranspileIfStatement(n)

	case "launch":
		return t.TranspileLaunchStatement(n)

	case "defer":
		return t.TranspileDeferStatement(n)

	case "enum":
		return t.TranspileEnumStatement(n)

	case "kv":
		return t.TranspileKeyValueStatement(n)

	case "map":
		return t.TranspileMapStatement(n)

	case "typedef":
		return t.TranspileTypeDeclaration(n)

	case "struct":
		return t.TranspileStructDeclaration(n)

	case "object":
		return t.TranspileObjectStatement(n)

	// FIXME: Why do we have expressions in here ... ?
	case "literal":
		return t.TranspileLiteralExpression(n)

	case "inc":
		var cppString, err = t.TranspileIncrementExpression(n)
		if err == nil {
			*cppString += ";"
		}

		return cppString, err

	case "dec":
		var cppString, err = t.TranspileDecrementExpression(n)
		if err == nil {
			*cppString += ";"
		}

		return cppString, err

	case "call":
		var cppString, err = t.TranspileCallExpression(n)
		if err == nil {
			*cppString += ";"
		}

		return cppString, err

	case "ident":
		return t.TranspileIdentExpression(n)

	case "use":
		return t.TranspileUseStatement(n)

	case "import":
		// importChan <- n
		return t.TranspileImportStatement(n)

	case "include":
		// includeChan <- n
		return nil, errors.Errorf("Direct C/C++ usage is not implemented yet: include: %+v\n", n)
		// return t.TranspileIncludeStatement(n)

	case "assignment":
		return t.TranspileAssignmentStatement(n)

	case "decl":
		return t.TranspileDeclarationStatement(n)

	case "let":
		return t.TranspileLetStatement(n)

	case "function":
		t.FuncChan <- n
		return nil, nil

	case "return":
		return t.TranspileReturnStatement(n)

	case "block":
		return t.TranspileBlockStatement(n)

	case "while":
		// while is not a user-facing keyword in Express2
		// it's used internally by the tree flattener to convert for-in/for-of loops
		return t.TranspileWhileStatement(n)

	// case "forof":
	// 	return t.TranspileForOfStatement(n)

	case "forin":
		return t.TranspileForInStatement(n)

	case "forover":
		return t.TranspileForOverStatement(n)

	case "forstd":
		return t.TranspileForStdStatment(n)

	case "package":
		// packageChan <- n
		return t.TranspilePackageStatement(n)

	case "selection":
		var exp, err = t.TranspileSelectExpression(n)
		if err != nil {
			return nil, err
		}

		*exp += ";"

		return exp, nil
	}

	return nil, errors.Errorf("Not implemented statement: %+v", n)
}

func (t *Transpiler) TranspilePackageStatement(n *builder.Node) (*string, error) {
	if n.Type != "package" {
		return nil, errors.New("Node is not a package statement")
	}

	pkgNameP, err := t.TranspileExpression(n.Left)
	if err != nil {
		return nil, err
	}
	pkgName := *pkgNameP

	t.log.Debug("STMTS LEN", len(n.Right.Value.([]*builder.Node)))

	stmts := n.Right.Value.([]*builder.Node)

	// Process the package body in two passes:
	// 1. Non-function statements (imports, struct declarations, etc.) inline
	// 2. Functions (regular and method) — transpiled directly and placed inside namespace

	var (
		blockContent string
		methods      = map[string][]string{} // receiver → method strings
		pkgFuncs     []string
	)

	for _, stmt := range stmts {
		if stmt.Type == "function" {
			// Transpile directly, bypassing FuncChan, so we can place inside namespace.
			// Hold funcMu to prevent races on CurrentReceiver/inParamContext with functionWorker.
			receiver, _ := stmt.Metadata["receiver"].(string)
			t.funcMu.Lock()
			t.CurrentReceiver = receiver
			fStr, ferr := t.TranspileFunctionStatement(stmt)
			t.CurrentReceiver = ""
			t.funcMu.Unlock()
			if ferr != nil {
				return nil, errors.Wrapf(ferr, "transpiling package function %q", stmt.Kind)
			}
			if receiver != "" {
				methods[receiver] = append(methods[receiver], *fStr)
			} else {
				pkgFuncs = append(pkgFuncs, *fStr)
			}
			continue
		}

		vStr, serr := t.TranspileStatement(stmt)
		if serr != nil {
			return nil, serr
		}
		if vStr != nil {
			blockContent += *vStr
		}
	}

	// Inject methods into the struct definition in blockContent.
	// blockContent contains "struct Foo { ... };" — find the right struct.
	for receiver, methodList := range methods {
		target := "struct " + receiver
		idx := strings.Index(blockContent, target)
		if idx < 0 {
			// Struct not found in this package block — skip (shouldn't happen normally)
			continue
		}
		// Find the closing "};" for this struct
		closeIdx := strings.Index(blockContent[idx:], "};")
		if closeIdx < 0 {
			continue
		}
		closeIdx += idx
		methodsStr := "\n" + strings.Join(methodList, "\n") + "\n"
		blockContent = blockContent[:closeIdx] + methodsStr + blockContent[closeIdx:]
	}

	// Assemble namespace block: {non-func content + pkg-level functions}
	var funcStr string
	for _, f := range pkgFuncs {
		funcStr += "\n" + f
	}

	nString := "namespace __" + pkgName + "{" + blockContent + funcStr + "}"

	t.log.Debug("NSTRING", nString)

	return &nString, nil
}

func (t *Transpiler) TranspileReturnStatement(n *builder.Node) (*string, error) {
	if n.Type != "return" {
		return nil, errors.New("Node is not a return statement")
	}

	// Return statments come in the form `return` { expr }

	var nString = "return"

	t.log.Debugf("n: %+v", n)

	// LHS (the return expression) is allowed to be empty
	if n.Left != nil {
		exprString, err := t.TranspileExpression(n.Left)
		if err != nil {
			return nil, err
		}

		nString += " " + *exprString
	}

	nString += ";"

	return &nString, nil
}

func (t *Transpiler) TranspileFunctionStatement(n *builder.Node) (*string, error) {
	if n.Type != "function" {
		return nil, errors.New("Node is not an function")
	}

	// Every function needs defer.cpp for the onReturn/onExit defer mechanism
	t.requirePathInclude(t.LibBase + "defer.cpp")

	/*
		A map with keys for `returns` and `args` will be egroups in the Metadata
		`Kind` is the name of the function
		`Value` is the block than needs to be translated
	*/

	if n.Kind == "" {
		return nil, errors.New("Somehow we parsed a function without a name ...")
	}

	// Start out with just the name; we will put the return type later
	var nString = n.Kind

	// args is an `sgroup`
	argsString, err := t.TranspileSGroup(n.Metadata["args"].(*builder.Node))
	if err != nil {
		return nil, err
	}

	// Append the args
	nString += *argsString

	var (
		returns = n.Metadata["returns"]

		returnsString = "void"

		// Start returns off as void
		returnsStringP = &returnsString
	)

	if n.Kind == "main" {
		returnsString = "int"
	}

	if returns != nil {
		// returns is a `type` for now; multiple returns are not supported right now
		returnsStringP, err = t.TranspileEGroup(returns.(*builder.Node))
		if err != nil {
			return nil, err
		}

		// For now shave off the parens since we are not supporting multiple types here
		if len(*returnsStringP) > 2 {
			returnsString = (*returnsStringP)[1 : len(*returnsStringP)-1]
			returnsStringP = &returnsString
		}
	}

	// Prepend the return string with a space
	nString = *returnsStringP + " " + nString

	blockString, err := t.TranspileBlockStatement(n.Value.(*builder.Node))
	if err != nil {
		return nil, err
	}

	addDeferToBlock(blockString)

	nString += *blockString

	return &nString, nil
}

func addDeferToBlock(blockP *string) {
	if blockP == nil {
		return
	}

	var block = *blockP

	block = block[0:1] + "defer onReturn, onExit;\n" + block[1:]

	*blockP = block
}

func (t *Transpiler) TranspileIdentExpression(n *builder.Node) (*string, error) {
	if n.Type != "ident" {
		blob, _ := json.Marshal(n)
		t.log.Debug("bbbbbbb:", string(blob))
		return nil, errors.New("Node is not an ident")
	}

	nString, ok := n.Value.(string)
	if !ok {
		return nil, errors.Errorf("Node value was not a string; %v", n)
	}

	// nil is C++'s nullptr
	if nString == "nil" {
		nString = "nullptr"
	}

	return &nString, nil
}

func (t *Transpiler) TranspileType(n *builder.Node) (*string, error) {
	if n.Type != "type" {
		blob, _ := json.Marshal(n)
		t.log.Debug("blob:", string(blob))
		return nil, errors.New("Node is not a type")
	}

	var nString string
	// Handle both string values (original) and node values (from let inference)
	if val, ok := n.Value.(string); ok {
		nString = val
	} else if n.Kind != "" {
		// Use Kind when Value is not a string (e.g., from let type inference)
		nString = n.Kind
	} else {
		return nil, errors.Errorf("Node value was not a string and Kind is empty; %v", n)
	}

	switch nString {
	case "var":
		t.requirePathInclude(t.LibBase + "var.cpp")

	case "string":
		nString = "std::" + nString
		t.requireStdInclude("string")

	case "map":
		nString = "map"
		t.requireStdInclude("map")
		t.requirePathInclude(t.LibBase + "var.cpp")

	// TODO(scottshotgg): for now every array will just be a list
	case "array":
		t.requireStdInclude("vector")

		elemKind := n.Kind
		if elemKind == "map" {
			// map[] → std::vector<std::map<std::string, var>>
			elemKind = "std::map<std::string, var>"
			t.requireStdInclude("map")
			t.requirePathInclude(t.LibBase + "var.cpp")
		}

		nString = "std::vector<" + elemKind + ">"

		// Handle nested vectors: int[][] has dim=[{none,-1},{none,-1}] → std::vector<std::vector<int>>
		if dims, ok2 := n.Metadata["dim"].([]*builder.Index); ok2 && len(dims) > 1 {
			for i := 1; i < len(dims); i++ {
				nString = "std::vector<" + nString + ">"
			}
		}

	case "pointer":
		nString = "*"
		if n.Left != nil && n.Left.Type == "selection" {
			// c.TYPE* pattern: Left is a selection node like c.FILE
			// Emit the C type name directly (e.g. FILE*)
			var typeName string
			if n.Left.Right != nil {
				typeName, _ = n.Left.Right.Value.(string)
			}
			if typeName == "" {
				return nil, errors.New("c.TYPE* selection: could not get type name")
			}
			nString = typeName + nString
		} else {
			var typeStringP, err = t.TranspileType(n.Left)
			if err != nil {
				return nil, err
			}
			nString = *typeStringP + nString
		}
	}

	// Check if the type is imported or not
	if n.Metadata["package"] != nil {
		var packageName = n.Metadata["package"].(string)
		// Get the type string from Value or Kind
		var typeName string
		if val, ok := n.Value.(string); ok {
			typeName = val
		} else {
			typeName = n.Kind
		}
		if packageName == "c" {
			// C types are global after #include — no namespace prefix needed.
			// c.FILE → FILE   c.DIR → DIR
			nString = typeName
		} else {
			nString = packageName + "::" + typeName
		}
	}

	return &nString, nil
}

// This changes an Express literal to be formatted the way C++ expects
func (t *Transpiler) prepLiteral(n *builder.Node, cpp string) (*string, error) {
	switch n.Kind {
	case "string":
		cpp = "\"" + cpp + "\""

	case "char":
		cpp = "'" + cpp + "'"

	case "struct":
		// Transpile the block for the value
		vString, err := t.TranspileBlockExpression(n.Right)
		if err != nil {
			return nil, errors.Wrap(err, "prepLiteral struct")
		}

		*vString = n.Value.(string) + *vString

		return vString, nil
	}

	return &cpp, nil
}

func (t *Transpiler) TranspileLiteralExpression(n *builder.Node) (*string, error) {
	if n.Type != "literal" {
		return nil, errors.New("Node is not an literal")
	}

	blob, _ := json.Marshal(n)
	t.log.Debug("its me again: n:", string(blob))

	return t.prepLiteral(n, fmt.Sprintf("%v", n.Value))
}

func (t *Transpiler) TranspileArrayExpression(n *builder.Node) (*string, error) {
	if n.Type != "array" {
		return nil, errors.New("Node is not an array")
	}

	var (
		nString = "{ "
		vString *string
		err     error
	)

	value := n.Value.([]*builder.Node)
	for _, v := range value {
		t.log.Debug("v:", *v)
		vString, err = t.TranspileExpression(v)
		if err != nil {
			return nil, err
		}

		t.log.Debug("vString:", vString)

		nString += *vString + ", "
	}

	// Cut off the last comma and space
	nString = nString[:len(nString)-2] + " }"

	return &nString, nil
}

func (t *Transpiler) TranspileAssignmentStatement(n *builder.Node) (*string, error) {
	if n.Type != "assignment" {
		return nil, errors.New("Node is not an assignment")
	}

	var (
		nString = ""
		vString *string
		err     error
	)

	// Left should be ident
	// Right should be general expression
	// This will require some prepping atleast to figure out
	// if we need any pre-statements

	// Translate the ident expression (lhs)
	vString, err = t.TranspileExpression(n.Left)
	if err != nil {
		return nil, err
	}

	nString = *vString + " = "

	// Translate the ident expression (lhs)
	vString, err = t.TranspileExpression(n.Right)
	if err != nil {
		return nil, err
	}

	t.log.Debug("IS THIS THE ONE:", *vString)

	nString += *vString + ";"

	return &nString, nil
}

func (t *Transpiler) TranspileDeclarationStatement(n *builder.Node, inParamCtx ...bool) (*string, error) {
	if n.Type != "decl" {
		return nil, errors.New("Node is not an declaration")
	}

	var (
		nString = ""
	)

	// Left should be ident
	// Right should be general expression
	// This will require some prepping atleast to figure out
	// if we need any pre-statements

	var tt = ""
	var typeOf = &tt
	var err error
	var constPrefix string

	if n.Left.Type != "deref" {
		typeNode, typeOk := n.Value.(*builder.Node)

		// C-style array declaration: char[dim] varName → char varName[dim] = {};
		// Fires for both `int[5] a` (no initializer) and `int[5] a = {}` (explicit empty init).
		isEmptyInit := n.Right != nil &&
			(n.Right.Type == "block" || n.Right.Type == "array") &&
			func() bool {
				nodes, ok := n.Right.Value.([]*builder.Node)
				return ok && len(nodes) == 0
			}()
		if typeOk && typeNode.Value == "array" && typeNode.Kind != "var" && (n.Right == nil || isEmptyInit) {
			if dims, ok2 := typeNode.Metadata["dim"].([]*builder.Index); ok2 && len(dims) > 0 && dims[0].Type != "none" {
				elemType := typeNode.Kind
				if elemType == "string" {
					elemType = "std::string"
				}
				varNameP, varErr := t.TranspileExpression(n.Left)
				if varErr != nil {
					return nil, varErr
				}
				var dimStr string
				d := dims[0]
				switch d.Type {
				case "ident":
					dimStr = d.Value.(string)
				case token.IntType:
					dimStr = fmt.Sprintf("%d", d.Value.(int))
				}
				nString = elemType + " " + *varNameP + "[" + dimStr + "] = {};"
				return &nString, nil
			}
		}

		typeOf, err = t.TranspileExpression(n.Value.(*builder.Node))
		if err != nil {
			return nil, err
		}

		t.log.Debug("TYPE", *typeOf, n.Left)

		// Prepend const for immutable bindings (not params, not pointers, not dynamic var, not mutable, not struct fields)
		if typeOk {
			inParam := len(inParamCtx) > 0 && inParamCtx[0]
			isPtr := typeNode.Kind == "pointer"
			isDynVar := typeNode.Kind == "var"
			isMut, _ := n.Metadata["mutable"].(bool)
			isField, _ := n.Metadata["is_field"].(bool)
			if !inParam && !isPtr && !isDynVar && !isMut && !isField {
				constPrefix = "const "
			}
		}
	}

	// Don't add the type yet

	// LHS is not allowed to be nil
	if n.Left == nil {
		return nil, errors.New("nil Left hand side")
	}

	// Translate the ident expression (lhs)
	vString, err := t.TranspileExpression(n.Left)
	if err != nil {
		return nil, err
	}

	nString += *vString

	// DESIGN DECISION: Auto-zero-initialization of uninitialized declarations
	//
	// Express auto-zero-initializes declarations without explicit initializers.
	// e.g., `int x` emits `int x = 0;` in C++, NOT `int x;` (which is UB in C++).
	//
	// Zero values by type:
	//   int, float  → 0
	//   bool        → false
	//   char        → '\0'
	//   pointer     → nullptr
	//   string      → "" (std::string default-constructs to empty)
	//   struct      → = {} (aggregate init, all fields recursively zeroed)
	//   var         → null (var class default-constructs to nullType)
	//   map         → {} (std::map default-constructs to empty)
	//   array       → = {} (int[5] — C-style, aggregate-init to all zeros)
	//   vector      → [] (int[] — std::vector default-constructs to empty)
	//
	// Alternatives considered:
	//   Option 1 (rejected): Leave uninitialized — causes undefined behavior in C++.
	//   Option 2 (deferred): Require explicit initializer at parse time (reject `int x`).
	//     Pros: forces documentation of intent. Cons: breaks declare-then-branch
	//     patterns without definite-assignment analysis (a future enhancement).
	//   Option 3 (chosen): Auto-zero-initialize at transpile time.
	//     Pros: eliminates UB, matches Go's zero-value philosophy, simple to implement.
	//     Cons: can hide "forgot to set a meaningful value" bugs (manifests as wrong
	//     result, not a crash). Acceptable tradeoff.
	//
	// If definite-assignment analysis is added later, Option 2 could replace this.
	// See also: `let` already requires an initializer; `var` has runtime nullType default.

	// RHS is allowed to be nil to support declarations without values like `string s`
	if n.Right == nil {
		typeNode := n.Value.(*builder.Node)

		// Do not zero-initialize function parameters — they receive their values from callers.
		inParam := len(inParamCtx) > 0 && inParamCtx[0]
		if !inParam {
			switch {
			case typeNode.Kind == "struct":
				nString += " = {}"
			case typeNode.Value == "int" || typeNode.Value == "float":
				nString += " = 0"
			case typeNode.Value == "bool":
				nString += " = false"
			case typeNode.Value == "char":
				nString += ` = '\0'`
			case typeNode.Kind == "pointer":
				nString += " = nullptr"
			default:
				// User-defined struct types (e.g. `S s`): the typeNode has Value="S"
				// and Kind="" at parse time. Look up the scope tree to detect composite types
				// and emit `= {}` so their fields are recursively zeroed.
				if typeName, ok := typeNode.Value.(string); ok {
					if tv := t.Builder.ScopeTree.GetType(typeName); tv != nil && tv.Composite {
						nString += " = {}"
					}
				}
			}
		}

		if *typeOf == "map" {
			var t = "std::map<std::string, var>"
			typeOf = &t
		}

		nString = constPrefix + *typeOf + " " + nString + ";"
		return &nString, nil
	}

	// Translate the ident expression (lhs)
	// May have to change this down the line or something
	switch *typeOf {
	case "map":
		vString, err = t.TranspileMapBlockStatement(n.Right)
		if err != nil {
			return nil, err
		}

		// typeOfBlock, err := t.DeduceMapBlockType(n.Right)
		// fmt.Println("typeOfBlock, err", *typeOfBlock, err)
		// os.Exit(9)

		kvs, ok := n.Right.Value.([]*builder.Node)
		if !ok {
			return nil, errors.New("kvs not ok")
		}

		var (
			varType   = token.VarType
			keyType   = &varType
			valueType = &varType
		)

		if len(kvs) > 0 {
			keyType, err = t.resolveType(kvs[0].Left)
			if err != nil {
				return nil, err
			}

			valueType, err = t.resolveType(kvs[0].Right)
			if err != nil {
				return nil, err
			}

			blob, _ := json.Marshal(kvs[0])
			t.log.Debug("kvblob:", string(blob))
		}

		nString = constPrefix + fmt.Sprintf("std::map<%s, %s> %s", *keyType, *valueType, nString)

	default:
		nString = constPrefix + *typeOf + " " + nString
		// For vector declarations with an empty `[]` initializer (e.g. `int[] v = []`),
		// the std::vector default constructor already produces an empty vector.
		// Emitting `= []` would produce malformed C++, so we just declare without an initializer.
		// Note: n.Right.Value is a nil []*builder.Node stored as interface{}, so == nil is false;
		// we check len == 0 instead.
		// Only applies to vectors (dim type == "none"), NOT to fixed-size arrays (int[4]).
		if typeNode, ok := n.Value.(*builder.Node); ok && typeNode.Value == "array" {
			isVector := false
			if dims, ok2 := typeNode.Metadata["dim"].([]*builder.Index); ok2 && len(dims) > 0 {
				isVector = dims[0].Type == "none"
			}
			if isVector {
				if nodes, ok2 := n.Right.Value.([]*builder.Node); ok2 && len(nodes) == 0 {
					nString += ";"
					return &nString, nil
				}
			}
		}
		vString, err = t.TranspileExpression(n.Right)
	}

	if err != nil {
		return nil, err
	}

	if n.Left.Type == "deref" && n.Left.Kind == "type" {
		nString += *vString
		return &nString, nil
	}

	nString += " = " + *vString + ";"

	// fmt.Println("nString", nString)

	return &nString, nil
}

func (t *Transpiler) TranspileLetStatement(n *builder.Node) (*string, error) {
	if n.Type != "let" {
		return nil, errors.New("Node is not a let statement")
	}

	// Translate the ident expression (lhs)
	vString, err := t.TranspileExpression(n.Left)
	if err != nil {
		return nil, err
	}

	var nString = *vString

	// RHS is allowed to be nil to support declarations without values like `string s`
	if n.Right == nil {
		return nil, errors.New("let statements must have a right-hand side expression")
	}

	// Infer the type from the right-hand side expression
	var typeString string
	switch n.Right.Type {
	case "literal":
		// Infer type from literal kind
		switch n.Right.Kind {
		case "int":
			typeString = "int"
		case "float":
			typeString = "float"
		case "string":
			typeString = "std::string"
			t.requireStdInclude("string")
		case "bool":
			typeString = "bool"
		case "char":
			typeString = "char"
		default:
			return nil, errors.Errorf("Unknown literal kind for let: %s", n.Right.Kind)
		}
	default:
		// For other expression types, use auto type inference
		typeString = "auto"
	}

	// Translate the expression (rhs)
	vString, err = t.TranspileExpression(n.Right)
	if err != nil {
		return nil, err
	}

	nString = "const " + typeString + " " + nString + " = " + *vString + ";"

	// fmt.Println("nString", nString)

	return &nString, nil
}

func (t *Transpiler) resolveType(n *builder.Node) (*string, error) {
	blob, _ := json.Marshal(n)
	t.log.Debug("vvvvvvv:", string(blob))
	switch n.Type {
	case "ident":
		var v = t.Builder.ScopeTree.Get(n.Value.(string))
		if v == nil {
			blob, _ := json.Marshal(t.Builder.ScopeTree)
			t.log.Debug("scopeTree:", string(blob))
			return nil, errors.New("ident not found:" + n.Value.(string))
		}

		var vv, ok = v.Value.(*builder.Node).Value.(string)
		if !ok {
			return nil, errors.New("could not get type")
		}

		return &vv, nil

	case "literal":
		switch n.Kind {
		case "struct":
			var v = n.Value.(string)
			return &v, nil
		}

		return &n.Kind, nil

	case "type":
		var t, ok = n.Value.(string)
		if !ok {
			return nil, errors.New("could not resolveType from type:")
		}

		return &t, nil

	case "block":
		return nil, errors.New("resolveType: cannot resolve type from a block node")
	}

	return nil, errors.New("unknown node type to resolveType: " + n.Type)
}

// func (t *Transpiler)  TranspileIncrementExpression(n *builder.Node) (*string, error) {
// 	if n.Type != "inc" {
// 		return nil, errors.New("Node is not an inc")
// 	}

// 	var (
// 		nString = ""
// 		vString *string
// 		err     error
// 	)

// 	// Translate the ident expression (lhs)
// 	vString, err = t.TranspileExpression(n.Left)
// 	if err != nil {
// 		return nil, err
// 	}

// 	nString += *vString + "++;"

// 	return &nString, nil
// }

func (t *Transpiler) TranspileConditionExpression(n *builder.Node) (*string, error) {
	if n.Type != "comp" {
		return nil, errors.New("Node is not an comp")
	}

	var (
		nString = ""
		vString *string
		err     error
	)

	// Translate the lhs
	vString, err = t.TranspileExpression(n.Left)
	if err != nil {
		return nil, err
	}

	nString += *vString

	// Translate the rhs
	vString, err = t.TranspileExpression(n.Right)
	if err != nil {
		return nil, err
	}

	nString += n.Value.(string) + *vString

	return &nString, nil
}

func (t *Transpiler) TranspileBinOpExpression(n *builder.Node) (*string, error) {
	if n.Type != "binop" {
		return nil, errors.New("Node is not a binop")
	}

	var (
		nString = ""
		vString *string
		err     error
	)

	// Translate the ident expression (lhs)
	vString, err = t.TranspileExpression(n.Left)
	if err != nil {
		return nil, err
	}

	nString += *vString + n.Value.(string)

	// Translate the ident expression (lhs)
	vString, err = t.TranspileExpression(n.Right)
	if err != nil {
		return nil, err
	}

	nString += *vString

	return &nString, nil
}

func (t *Transpiler) TranspileBlockStatement(n *builder.Node) (*string, error) {
	if n == nil {
		var nString = "{}"
		return &nString, nil
	}

	if n.Type != "block" {
		return nil, errors.New("Node is not a block")
	}

	var (
		nString = ""
		vString *string
		err     error
	)

	for _, stmt := range n.Value.([]*builder.Node) {
		// Skip typedef declarations - they're handled at global scope
		if stmt.Type == "typedef" {
			continue
		}

		vString, err = t.TranspileStatement(stmt)
		if err != nil {
			return nil, err
		}

		if vString == nil {
			continue
		}

		t.log.Debug("vString", *vString)

		// if stmt.Type != "function" {
		nString += *vString
		// }
	}

	nString = "{" + nString + "}"

	return &nString, nil
}

func (t *Transpiler) TranspileMapBlockStatement(n *builder.Node) (*string, error) {
	if n.Type != "block" {
		return nil, errors.New("Node is not a block")
	}

	var (
		nString = ""
		vString *string
		err     error
	)

	for _, stmt := range n.Value.([]*builder.Node) {
		if stmt.Type != "kv" {
			return nil, errors.Errorf("All statements in a map have to be key-value pairs: %+v\n", stmt)
		}

		vString, err = t.TranspileStatement(stmt)
		if err != nil {
			return nil, err
		}

		nString += *vString + ","
	}

	nString = "{" + nString + "}"

	return &nString, nil
}

// type PairedType struct {
// 	Left  string
// 	Right string
// }

// TODO: This should go in the checker

// // Maybe use the above
// func (t *Transpiler) DeduceMapBlockType(n *builder.Node) (*string, error) {
// 	// var firstType, secondType string

// 	/*
// 		This is used to determine what type of map the block is so that we
// 		can put the type into the C++ type
// 	*/

// 	if n.Type != "block" {
// 		return nil, errors.New("Node is not a block")
// 	}

// 	var (
// 		nString = ""
// 		// vString *string
// 		// err     error
// 	)

// 	for _, stmt := range n.Value.([]*builder.Node) {
// 		if stmt.Type != "kv" {
// 			return nil, errors.Errorf("All statements in a map have to be key-value pairs: %+v\n", stmt)
// 		}

// 		var leftType = stmt.Left.Kind
// 		var rightType = stmt.Right.Kind
// 		fmt.Printf("stmt kind: left: %s right: %s\n", leftType, rightType)
// 		if leftType == "" {
// 			fmt.Println("t.scope", t.Builder.ScopeTree.Get(stmt.Left.Value.(string)))
// 			os.Exit(9)
// 		}
// 	}

// 	nString = "<var, var>"
// 	return &nString, nil
// }

// TODO: for now all variables will have a `.` in front, later on it
// should only be the public variables
func (t *Transpiler) TranspileBlockExpression(n *builder.Node) (*string, error) {
	if n == nil {
		var nString = "{}"
		return &nString, nil
	}

	if n.Type != "block" {
		return nil, errors.New("Node is not a block")
	}

	var (
		nString = ""
		vString *string
		err     error
	)

	stmts := n.Value.([]*builder.Node)

	// If the block contains kv pairs it's a map literal — emit var{k1, v1, k2, v2, ...}
	// using the var(initializer_list<var>) constructor in var.cpp.
	if len(stmts) > 0 && stmts[0].Type == "kv" {
		t.requireStdInclude("map")
		t.requirePathInclude(t.LibBase + "var.cpp")

		var parts []string
		for _, kv := range stmts {
			k, err := t.TranspileExpression(kv.Left)
			if err != nil {
				return nil, err
			}
			v, err := t.TranspileExpression(kv.Right)
			if err != nil {
				return nil, err
			}
			parts = append(parts, *k, *v)
		}
		nString = "var{" + strings.Join(parts, ", ") + "}"
		return &nString, nil
	}

	// TODO: don't have a type checker so for right now
	// just type check in here
	for _, stmt := range stmts {
		vString, err = t.TranspileStatement(stmt)
		if err != nil {
			return nil, err
		}

		// Don't check this here; leave it to the type checker
		// if stmt.Type != "assignment" {
		// 	return nil, errors.Errorf("Structs can only contain assignment statements; %+v", stmt)
		// }

		*vString = (*vString)[:len(*vString)-1]

		// Add a dot in front
		nString += "." + *vString + ","
	}

	nString = "{" + nString + "}"

	return &nString, nil
}

func (t *Transpiler) TranspileEGroup(n *builder.Node) (*string, error) {
	if n.Type != "egroup" {
		return nil, errors.New("Node is not a egroup")
	}

	var (
		nStrings = make([]string, len(n.Value.([]*builder.Node)))
		vString  *string
		err      error
	)

	for i, e := range n.Value.([]*builder.Node) {
		vString, err = t.TranspileExpression(e)
		if err != nil {
			return nil, err
		}

		nStrings[i] = *vString
	}

	var nString = "(" + strings.Join(nStrings, ",") + ")"

	return &nString, nil
}

func (t *Transpiler) TranspileStreamEGroup(n *builder.Node) (*string, error) {
	if n.Type != "egroup" {
		return nil, errors.New("Node is not a egroup")
	}

	var (
		nStrings = make([]string, len(n.Value.([]*builder.Node)))
		vString  *string
		err      error
	)

	for i, e := range n.Value.([]*builder.Node) {
		vString, err = t.TranspileExpression(e)
		if err != nil {
			return nil, err
		}

		nStrings[i] = *vString
	}

	var nString = strings.Join(nStrings, "<< \" \" << ")

	return &nString, nil
}

func (t *Transpiler) TranspileSGroup(n *builder.Node) (*string, error) {
	if n.Type != "sgroup" {
		return nil, errors.New("Node is not a sgroup")
	}

	var (
		nStrings = make([]string, len(n.Value.([]*builder.Node)))
		vString  *string
		err      error
	)

	// Transpile each parameter declaration directly (inParamCtx=true) so that
	// parameter declarations are not zero-initialized — they receive values from callers.
	// Using the parameter directly avoids the shared t.inParamContext field race.
	for i, s := range n.Value.([]*builder.Node) {
		vString, err = t.TranspileDeclarationStatement(s, true)
		if err != nil {
			return nil, err
		}

		// Shave off the semicolon since we don't need it in a parameter list
		var vvString = *vString
		if vvString[len(vvString)-1] == ';' {
			vvString = (vvString)[:len(vvString)-1]
		}

		nStrings[i] = vvString
	}

	var nString = "(" + strings.Join(nStrings, ",") + ")"

	return &nString, nil
}

func (t *Transpiler) TranspileCallExpression(n *builder.Node) (*string, error) {
	if n.Type != "call" {
		return nil, errors.New("Node is not a call")
	}

	var (
		nString = ""
		vString *string
		err     error
	)

	vString, err = t.TranspileExpression(n.Value.(*builder.Node))
	if err != nil {
		return nil, err
	}

	blob, _ := json.Marshal(vString)
	t.log.Debug("vsTrIng:", string(blob))

	nString += *vString

	var args = n.Metadata["args"]

	blob, _ = json.Marshal(n.Value.(*builder.Node))
	t.log.Debug("n.Value.(*builder.Node):", string(blob))

	// Just do the checking here for now, not sure the merits of making the sgroup function check
	if args == nil {
		nString += "()"
		return &nString, nil
	}

	var argString string

	funcName, ok := n.Value.(*builder.Node).Value.(string)
	if ok {
		if cFuncs[n.Value.(*builder.Node).Value.(string)] {
			t.requireStdInclude("stdio.h")
		}

		if funcName == "len" {
			// len(x) → (x).size()
			// Works for: vectors (int[], var[], etc.) and strings.
			// NOT supported for arrays (int[5]) — size is a compile-time constant.
			argString, err := t.TranspileEGroup(args.(*builder.Node))
			if err != nil {
				return nil, err
			}
			inner := (*argString)[1 : len(*argString)-1] // strip outer parens
			nString = "(" + inner + ").size()"
			return &nString, nil
		} else if funcName == "Println" {
			var argString, err = t.TranspileStreamEGroup(args.(*builder.Node))
			if err != nil {
				return nil, err
			}

			var ending = " std::endl"

			if len(*argString) > 0 {
				ending = "<<" + ending
			}

			nString = "std::cout <<" + *argString + ending

			return &nString, nil
		} else if funcName == "sleep" {
			t.requireStdInclude("unistd.h")
		}
	}

	vString, err = t.TranspileEGroup(args.(*builder.Node))
	if err != nil {
		return nil, err
	}
	blob, _ = json.Marshal(vString)
	t.log.Debug("egroup vstring:", string(blob))

	argString += *vString

	blob, _ = json.Marshal(vString)
	t.log.Debug("argstring:", string(blob))

	nString += argString

	return &nString, nil
}

var cFuncs = map[string]bool{
	"Println": true,
	"printf":  true,
	"sleep":   true,
	"msleep":  true,
	"now":     true,
}

func (t *Transpiler) TranspileForInStatement(n *builder.Node) (*string, error) {
	var nString = fmt.Sprintf("for (auto const& %s : %s)", n.Left.Left.Value.(string), n.Right.Value.(string))

	// Translate the block statement
	vString, err := t.TranspileBlockStatement(n.Value.(*builder.Node))
	if err != nil {
		return nil, err
	}

	nString += *vString

	return &nString, nil
}

func (t *Transpiler) TranspileForOverStatement(n *builder.Node) (*string, error) {
	if n.Type != "forover" {
		return nil, errors.New("Node is not a forover")
	}

	keyIdent := n.Metadata["start"].(*builder.Node).Value.(string)
	collection, err := t.TranspileExpression(n.Metadata["end"].(*builder.Node))
	if err != nil {
		return nil, err
	}

	var nString string
	var bodyPrefix string

	if ident2, ok := n.Metadata["start2"]; ok {
		// Two-variable form: for i, j over x
		// Transpile as an indexed for-loop with both key and value
		valIdent := ident2.(*builder.Node).Value.(string)
		nString = fmt.Sprintf("{auto& _coll = %s; for (int %s = 0; %s < _coll.size(); %s++)",
			*collection, keyIdent, keyIdent, keyIdent)
		bodyPrefix = fmt.Sprintf("auto %s = _coll[%s];", valIdent, keyIdent)
	} else {
		// Single-variable form: for i over x
		// Transpile as an indexed for-loop; i gets the index
		nString = fmt.Sprintf("{auto& _coll = %s; for (int %s = 0; %s < _coll.size(); %s++)",
			*collection, keyIdent, keyIdent, keyIdent)
	}

	// Translate the body
	vString, err := t.TranspileBlockStatement(n.Value.(*builder.Node))
	if err != nil {
		return nil, err
	}

	if bodyPrefix != "" {
		// Inject the value declaration at the start of the block body
		// TranspileBlockStatement returns "{ ... }", so insert after the opening brace
		body := *vString
		nString += "{" + bodyPrefix + body[1:]
	} else {
		nString += *vString
	}

	// Close the outer scoping block
	nString += "}"

	return &nString, nil
}

func (t *Transpiler) TranspileForStdStatment(n *builder.Node) (*string, error) {
	// Change forin to be a block statement containing:
	//	- declare temp var
	//	- declare array/iter
	//	- while tempvar < iter.length
	//	- var = iter[tempvar]
	//	- loop_block
	//	-	increment var

	// return nil, errors.New("not implemented: forof")

	if n.Type != "forstd" {
		return nil, errors.New("Node is not a forstd")
	}

	var (
		nString = "{"
		vString *string
		err     error
	)

	// Make and translate the ident into a declaration
	ds := TransformIdentToDefaultDeclaration(n.Metadata["start"].(*builder.Node))
	vString, err = t.TranspileDeclarationStatement(ds)
	if err != nil {
		return nil, err
	}

	nString += *vString

	// Make and translate the array expression into a declaration
	dss := TransformExpressionToDeclaration(n.Metadata["end"].(*builder.Node))
	vString, err = t.TranspileDeclarationStatement(dss)
	if err != nil {
		return nil, err
	}

	nString += *vString

	// Make and translate the less than operation
	vString, err = t.TranspileConditionExpression(&builder.Node{
		Type:  "comp",
		Value: "<",
		Left:  n.Metadata["start"].(*builder.Node),
		Right: GenerateLengthCall(dss),
	})
	if err != nil {
		return nil, err
	}

	nString += fmt.Sprintf("while(%s)", *vString)

	// Translate the block statement
	vString, err = t.TranspileBlockStatement(n.Value.(*builder.Node))
	if err != nil {
		return nil, err
	}

	nString += *vString

	// Lastly, make and translate an increment statement for the ident
	vString, err = t.TranspileIncrementExpression(&builder.Node{
		Type: "inc",
		Left: n.Metadata["start"].(*builder.Node),
	})
	if err != nil {
		return nil, err
	}

	nString += *vString + "}"

	return &nString, nil
}

// TranspileWhileStatement handles while loops.
// Note: while is NOT a user-facing keyword in Express2.
// It's used internally by the tree flattener to convert for-in/for-of loops.
// User code cannot directly write while loops.
func (t *Transpiler) TranspileWhileStatement(n *builder.Node) (*string, error) {
	/*
		while statements are simple, we already have all the tools:
		`while` `(` expr `)` block
	*/

	if n.Type != "while" {
		return nil, errors.New("Node is not a while")
	}

	var (
		nString = "while ("
		// vString *string
		// err     error
	)

	condition, err := t.TranspileExpression(n.Left)
	if err != nil {
		return nil, err
	}

	nString += *condition + ")"

	block, err := t.TranspileBlockStatement(n.Value.(*builder.Node))
	if err != nil {
		return nil, err
	}

	nString += *block

	return &nString, nil
}

func (t *Transpiler) TranspileIfStatement(n *builder.Node) (*string, error) {
	/*
		Form:
		`if` [expr] [block] {`else` {epxr} [block]}

		Value is the condition
		Left is the block
		Right is the else statement
	*/

	if n.Type != "if" {
		return nil, errors.Errorf("Node is not an if: %+v", n)
	}

	var (
		nString = "if ("
		vString *string
		err     error
	)

	vString, err = t.TranspileExpression(n.Value.(*builder.Node))
	if err != nil {
		return nil, err
	}

	// Add the condition and the parenthesis
	nString += *vString + ")"

	vString, err = t.TranspileBlockStatement(n.Left)
	if err != nil {
		return nil, err
	}

	// Add the block; this should already come with the curly braces
	nString += *vString

	// If the Right child is non-nil then we have an else block in the form of an if statement
	if n.Right != nil {
		// Check whether it is an elseif of just an else
		switch n.Right.Type {
		case "if":
			vString, err = t.TranspileIfStatement(n.Right)
			if err != nil {
				return nil, err
			}

		case "block":
			vString, err = t.TranspileBlockStatement(n.Right)
			if err != nil {
				return nil, err
			}

		default:
			return nil, errors.Errorf("Node is not an if or block: %+v", n.Right)
		}

		// Add the else block and nest the if statement inside of it
		// TODO: research what implication this has in LLVM
		nString += "else " + *vString
	}

	return &nString, nil
}

// func ConstructGreaterThanOperation(op string, lhs, rhs *builder.Node)

func GenerateLengthCall(n *builder.Node) *builder.Node {
	// TODO: need to have a switch here for arrays and stuff
	return &builder.Node{
		Type: "call",
		Value: &builder.Node{
			Type:  "ident",
			Value: "std::size",
		},
		Metadata: map[string]interface{}{
			"args": &builder.Node{
				Type:  "egroup",
				Value: []*builder.Node{n.Left},
			},
		},
	}
}

func TransformExpressionToDeclaration(n *builder.Node) *builder.Node {
	// TODO: Type checker would give type here; use auto for now
	return &builder.Node{
		Type: "decl",
		Value: &builder.Node{
			Type:  "type",
			Value: "auto",
		},
		Left: &builder.Node{
			Type:  "ident",
			Value: "SOMETHING",
		},
		Right: n,
	}
}

func TransformIdentToDefaultDeclaration(n *builder.Node) *builder.Node {
	return &builder.Node{
		Type: "decl",
		Value: &builder.Node{
			Type:  "type",
			Value: "int",
			Kind:  "int",
		},
		Left:     n,
		Metadata: map[string]interface{}{"mutable": true},
		Right: &builder.Node{
			Type:  "literal",
			Value: 0,
		},
	}
}
