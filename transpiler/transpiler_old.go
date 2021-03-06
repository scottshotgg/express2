package transpiler

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	token "github.com/scottshotgg/express-token"
	"github.com/scottshotgg/express2/builder"
	"github.com/scottshotgg/express2/tree_flattener"
)

type WithPriority struct {
	Priority int
	Value    string
}

type Transpiler struct {
	LibBase      string
	Name         string
	Builder      *builder.Builder
	AST          *builder.Node
	ASTCloneJSON []byte
	Extra        []string
	Functions    map[string]string
	Imports      map[string]string
	Includes     map[string]string
	Packages     map[string]string
	Types        map[string]string
	Structs      map[string]WithPriority
	GenerateMain bool
	Wg           *sync.WaitGroup
	Wg1          *sync.WaitGroup
	FuncChan     chan *builder.Node
	TypeChan     chan *builder.Node
	StructChan   chan *builder.Node
	IncludeChan  chan *builder.Node
	PackageChan  chan *builder.Node
	ImportChan   chan *builder.Node
	AppendChan   chan string
}

func (t *Transpiler) emit(line string) {
	t.AppendChan <- line
}

func (t *Transpiler) appendWorker() {
	defer t.Wg1.Done()

	var totalFile string

	for a := range t.AppendChan {
		totalFile += a
	}

	fmt.Println("totalFile", totalFile)
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

		if t.Functions[functionName] != "" {
			// FIXME: this is an error
			log.Printf("Function already declared: %+v\n", f)
			os.Exit(9)
		}

		function, err = t.TranspileFunctionStatement(f)
		if err != nil {
			log.Printf("Function error: %+v %+v\n", f, err)
			os.Exit(9)
		}

		t.Functions[functionName] = *function
	}
}

func New(ast *builder.Node, b *builder.Builder, name, libBase string) *Transpiler {
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
		FuncChan:    make(chan *builder.Node, 100),
		TypeChan:    make(chan *builder.Node, 100),
		StructChan:  make(chan *builder.Node, 100),
		IncludeChan: make(chan *builder.Node, 100),
		PackageChan: make(chan *builder.Node, 100),
		ImportChan:  make(chan *builder.Node, 100),
		AppendChan:  make(chan string, 5),
		Wg:          &sync.WaitGroup{},
		Wg1:         &sync.WaitGroup{},
	}

	// go appendWorker(&wg)

	t.ASTCloneJSON, _ = json.Marshal(ast)

	return &t
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

	t.Wg1.Add(1)
	go t.typeWorker()

	t.Wg1.Add(1)
	go t.structWorker()

	t.Wg1.Add(1)
	go t.packageWorker()

	t.Wg.Add(1)
	go t.includeWorker()

	t.Wg.Add(1)
	go t.importWorker()

	// Flatten the tree
	includes, err := tree_flattener.New().Flatten(t.AST)
	if err != nil {
		return err
	}

	for i := range includes {
		t.IncludeChan <- includes[i]
	}

	for i := range nodes {
		blob, _ := json.Marshal(nodes[i])
		fmt.Println("blob:", string(blob))
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
			t.IncludeChan <- nodes[i]

		case "decl":
			// Just transpile the statement for now
			stringP, err := t.TranspileStatement(nodes[i])
			if err != nil {
				fmt.Printf("err %+v\n", err)
				os.Exit(9)
				// return "", err
			}

			t.Extra = append(t.Extra, *stringP)

		case "use":
			// Just transpile the statement for now
			stringP, err := t.TranspileStatement(nodes[i])
			if err != nil {
				fmt.Printf("err %+v\n", err)
				os.Exit(9)
				// return "", err
			}

			t.Extra = append(t.Extra, *stringP)

		case "map":
			// Just transpile the statement for now
			stringP, err := t.TranspileStatement(nodes[i])
			if err != nil {
				fmt.Printf("err %+v\n", err)
				os.Exit(9)
				// return "", err
			}

			t.Extra = append(t.Extra, *stringP)

		case "package":
			t.PackageChan <- nodes[i]

		default:
			return errors.Errorf("Node was not categorized properly: %+v\n", nodes[i])
		}

		// // Just transpile the statement for now
		// stringP, err := t.TranspileStatement(nodes[i])
		// if err != nil {
		// 	fmt.Printf("err %+v\n", err)
		// 	os.Exit(9)
		// 	// return "", err
		// }

		// t.Extra = append(t.Extra, *stringP)
	}

	// Just a fucking dirty ass hackerino
	time.Sleep(1 * time.Second)

	// These are over used. Really the only reason that the function, struct, and type
	// chans were here in the first place was to capture all of the stuff to put it at the top
	// but tbh this should be a semantic parser step before it even gets to the AST

	// Close the channel and alert the worker that we are done
	close(t.PackageChan)
	close(t.FuncChan)
	close(t.TypeChan)
	close(t.StructChan)

	// Wait for everything to be transpiled
	t.Wg1.Wait()

	close(t.ImportChan)
	close(t.IncludeChan)

	// Wait for everything to be transpiled
	t.Wg.Wait()

	if t.Functions["main"] == "" {
		// return "", errors.New("No main function declared")
	}

	return nil
}

func (t *Transpiler) typeWorker() {
	defer t.Wg1.Done()

	var (
		stringP *string
		err     error
	)

	for node := range t.TypeChan {
		stringP, err = t.TranspileStatement(node)
		if err != nil {
			fmt.Printf("err %+v\n", err)
			os.Exit(9)
			// return "", err
		}

		t.Types[node.Left.Value.(string)] = *stringP
	}
}

func (t *Transpiler) structWorker() {
	defer t.Wg1.Done()

	var (
		stringP *string
		err     error
	)

	var i int
	for node := range t.StructChan {
		stringP, err = t.TranspileStatement(node)
		if err != nil {
			fmt.Printf("err %+v\n", err)
			os.Exit(9)
		}

		t.Structs[node.Left.Value.(string)] = WithPriority{
			Priority: i,
			Value:    *stringP,
		}

		i++
	}
}

func (t *Transpiler) packageWorker() {
	defer t.Wg1.Done()

	var (
		stringP *string
		err     error
	)

	var i int
	for node := range t.PackageChan {
		stringP, err = t.TranspileStatement(node)
		if err != nil {
			fmt.Printf("err %+v\n", err)
			os.Exit(9)
		}

		t.Packages[node.Left.Value.(string)] = *stringP

		i++
	}
}

func (t *Transpiler) includeWorker() {
	defer t.Wg.Done()

	var (
		includeStringP *string
		// Why does this shadow ...
		// Is the gofunc "capturing" variables that aren't passed?
		ierr error
	)

	for node := range t.IncludeChan {
		// Might want to make this go through the entire pipeline ...
		includeStringP, ierr = t.TranspileIncludeStatement(node)
		if ierr != nil {
			log.Printf("Error transpiling include statement: %+v\n", ierr)

			// Exit if there is a problem transpiling the import statement
			// and we'll deal with it later
			os.Exit(9)
		}

		// TODO: should really check the deref on all of these, but the usage/running
		// is pretty predictable right now
		t.Includes[node.Left.Value.(string)] = *includeStringP
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
			log.Printf("Error transpiling import statement: %+v\n", ierr)

			// Exit if there is a problem transpiling the import statement
			// and we'll deal with it later
			os.Exit(9)
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
	for _, t := range t.Structs {
		structs[t.Priority] = t.Value
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

	fmt.Println("mainFunc", t.Functions["main"])
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
		fmt.Printf("n %+v\n", n)
		fmt.Printf("n %+v\n", n.Left)
		fmt.Printf("n %+v\n", n.Right)
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

	// If it is the import for libc
	if n.Kind == "c" {
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

	t.Packages[n.Left.Value.(string)] = fmt.Sprintf("namespace %s {\n %s\n }\n", "__"+n.Left.Value.(string), tt.ToCpp())
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

	var selector = "."
	if n.Left.Type == "package" {
		selector = "::"
	}

	var nString = *lhs + selector + *rhs

	return &nString, nil
}

// The type checker should produce arrow function ones as well
func (t *Transpiler) TranspileDerefExpression(n *builder.Node) (*string, error) {
	if n.Type != "deref" {
		return nil, errors.New("Node is not an deref")
	}

	// Left is the ident; right is nothing
	lhs, err := t.TranspileExpression(n.Left)
	if err != nil {
		return nil, err
	}

	var nString = "*" + *lhs

	return &nString, nil
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

func (t *Transpiler) TranspileMapStatement(n *builder.Node) (*string, error) {
	if n.Type != "map" {
		return nil, errors.New("Node is not a map")
	}

	// Transpile the ident
	vString, err := t.TranspileExpression(n.Left)
	if err != nil {
		return nil, err
	}

	// Could just have it add `struct` here but this will show us changes
	var nString = "std::map<std::string, std::string>" + " " + *vString + "= "

	// Transpile the block for the value
	vString, err = t.TranspileMapBlockStatement(n.Right)
	if err != nil {
		return nil, err
	}

	// Include std::map from C++
	t.IncludeChan <- &builder.Node{
		Type: "include",
		Left: &builder.Node{
			Type:  "literal",
			Value: "map",
		},
	}

	nString += *vString + ";"

	return &nString, nil
}

func (t *Transpiler) TranspileLaunchStatement(n *builder.Node) (*string, error) {
	if n.Type != "launch" {
		return nil, errors.New("Node is not a launch node")
	}

	blob, _ := json.Marshal(n)
	fmt.Println("launch:", string(blob))

	blob, _ = json.Marshal(n)
	fmt.Println("launch.left:", string(blob))

	// Transpile the ident
	var vString, err = t.TranspileStatement(n.Left)
	if err != nil {
		return nil, err
	}

	// Include libmill for coroutines
	t.IncludeChan <- &builder.Node{
		Type: "include",
		Kind: "path",
		Left: &builder.Node{
			Type:  "literal",
			Value: t.LibBase + "libmill/libmill.h",
		},
	}

	// includeChan <- &builder.Node{
	// 	Type: "include",
	// 	// This is not supposed to be a `path` import; it is a library feature, stop re-adding that shit
	// 	Left: &builder.Node{
	// 		Type:  "literal",
	// 		Value: "libmill.h",
	// 	},
	// }

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
	var nString = "onReturn.deferStack.push([=](...){" + *stmt + "});"

	return &nString, nil
}

func (t *Transpiler) TranspileStatement(n *builder.Node) (*string, error) {
	fmt.Println("wtf3333", n.Type)
	switch n.Type {

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
		// typeChan <- n
		return t.TranspileTypeDeclaration(n)

	case "struct":
		// structChan <- n
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

	case "function":
		t.FuncChan <- n
		// Right now just grab the name from the string
		// Later on we can issue new function names for lambdas
		// var name = n.Kind
		// return &name, nil
		// return t.TranspileFunctionStatement(n)

	case "return":
		return t.TranspileReturnStatement(n)

	case "block":
		return t.TranspileBlockStatement(n)

	case "while":
		return t.TranspileWhileStatement(n)

	case "forof":
		return t.TranspileForOfStatement(n)

	case "forin":
		return t.TranspileForInStatement(n)

	case "forstd":
		// TODO:

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

	var (
		nString      = "namespace "
		vString, err = t.TranspileExpression(n.Left)
	)

	if err != nil {
		return nil, err
	}

	nString += " " + *vString

	// Get all of the statements inside the package

	fmt.Println("STMTS LEN", len(n.Right.Value.([]*builder.Node)))

	vString, err = t.TranspileBlockStatement(n.Right)
	if err != nil {
		return nil, err
	}

	nString += *vString

	fmt.Println("NSTRING", nString)

	return &nString, nil
}

func (t *Transpiler) TranspileReturnStatement(n *builder.Node) (*string, error) {
	if n.Type != "return" {
		return nil, errors.New("Node is not a return statement")
	}

	// Return statments come in the form `return` { expr }

	var nString = "return"

	fmt.Printf("n: %+v\n", n)

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

	// Include std::map from C++
	t.IncludeChan <- &builder.Node{
		Type: "include",
		Kind: "path",
		Left: &builder.Node{
			Type:  "literal",
			Value: t.LibBase + "defer.cpp",
		},
	}

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
		fmt.Println("bbbbbbb:", string(blob))
		return nil, errors.New("Node is not an ident")
	}

	nString, ok := n.Value.(string)
	if !ok {
		return nil, errors.Errorf("Node value was not a string; %v", n)
	}

	return &nString, nil
}

func (t *Transpiler) TranspileType(n *builder.Node) (*string, error) {
	if n.Type != "type" {
		blob, _ := json.Marshal(n)
		fmt.Println("blob:", string(blob))
		return nil, errors.New("Node is not a type")
	}

	nString, ok := n.Value.(string)
	if !ok {
		return nil, errors.Errorf("Node value was not a string; %v", n)
	}

	switch nString {
	case "string":
		nString = "std::" + nString

		// TODO: switch this to just use a damn string later
		t.IncludeChan <- &builder.Node{
			Type: "include",
			Left: &builder.Node{
				Type:  "literal",
				Value: "string",
			},
		}

	case "map":
		nString = "map"

		t.IncludeChan <- &builder.Node{
			Type: "include",
			Left: &builder.Node{
				Type:  "literal",
				Value: "map",
			},
		}

		t.IncludeChan <- &builder.Node{
			Type: "include",
			Kind: "path",
			Left: &builder.Node{
				Type:  "literal",
				Value: t.LibBase + "var.cpp",
			},
		}

	// TODO(scottshotgg): for now every array will just be a list
	case "array":
		t.IncludeChan <- &builder.Node{
			Type: "include",
			Left: &builder.Node{
				Type:  "literal",
				Value: "vector",
			},
		}

		nString = "std::vector<" + n.Kind + ">"

	case "pointer":
		nString = "*"
		var typeStringP, err = t.TranspileType(n.Left)
		if err != nil {
			return nil, err
		}

		nString = *typeStringP + nString
	}

	// Check if the type is imported or not
	if n.Metadata["package"] != nil {
		nString = n.Metadata["package"].(string) + "::" + n.Value.(string)
	}

	return &nString, nil
}

// This changes an Express literal to be formatted the way C++ expects
func (t *Transpiler) prepLiteral(n *builder.Node, cpp string) *string {
	switch n.Kind {
	case "string":
		cpp = "\"" + cpp + "\""

	case "char":
		cpp = "'" + cpp + "'"

	case "struct":
		// Transpile the block for the value
		vString, err := t.TranspileBlockExpression(n.Right)
		if err != nil {
			fmt.Println("err:", err)
			os.Exit(9)
		}

		*vString = n.Value.(string) + *vString

		return vString
	}

	return &cpp
}

func (t *Transpiler) TranspileLiteralExpression(n *builder.Node) (*string, error) {
	if n.Type != "literal" {
		return nil, errors.New("Node is not an literal")
	}

	blob, _ := json.Marshal(n)
	fmt.Println("its me again: n:", string(blob))

	return t.prepLiteral(n, fmt.Sprintf("%v", n.Value)), nil
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
		fmt.Println("v:", *v)
		vString, err = t.TranspileExpression(v)
		if err != nil {
			return nil, err
		}

		fmt.Println("vString:", vString)

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

	fmt.Println("IS THIS THE ONE:", *vString)

	nString += *vString + ";"

	return &nString, nil
}

func (t *Transpiler) TranspileDeclarationStatement(n *builder.Node) (*string, error) {
	if n.Type != "decl" {
		return nil, errors.New("Node is not an declaration")
	}

	t.Builder.ScopeTree.Declare(n)

	var (
		nString = ""
	)

	// Left should be ident
	// Right should be general expression
	// This will require some prepping atleast to figure out
	// if we need any pre-statements

	typeOf, err := t.TranspileType(n.Value.(*builder.Node))
	if err != nil {
		return nil, err
	}

	log.Println("TYPE", *typeOf, n.Left)

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

	// RHS is allowed to be nil to support declarations without values like `string s`
	if n.Right == nil {
		// If the declaration is a struct, then give it the default init if no expression is provided
		if n.Value.(*builder.Node).Kind == "struct" {
			nString += "={}"
		}

		if *typeOf == "map" {
			var t = "std::map<var, var>"
			typeOf = &t
		}

		// log.Printf("HEY ITS ME %+v\n", n.Value.(*builder.Node))

		nString = *typeOf + " " + nString + ";"
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
			fmt.Println("kvblob:", string(blob))
		}

		nString = fmt.Sprintf("std::map<%s, %s> %s", *keyType, *valueType, nString)

	default:
		nString = *typeOf + " " + nString
		vString, err = t.TranspileExpression(n.Right)
	}

	if err != nil {
		return nil, err
	}

	nString += " = " + *vString + ";"

	// fmt.Println("nString", nString)

	return &nString, nil
}

func (t *Transpiler) resolveType(n *builder.Node) (*string, error) {
	blob, _ := json.Marshal(n)
	fmt.Println("vvvvvvv:", string(blob))
	switch n.Type {
	case "ident":
		var v = t.Builder.ScopeTree.Get(n.Value.(string))
		if v == nil {
			blob, _ := json.Marshal(t.Builder.ScopeTree)
			fmt.Println("scopeTree:", string(blob))
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
		fmt.Println("I AM HERE")
		os.Exit(9)
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
		vString, err = t.TranspileStatement(stmt)
		if err != nil {
			return nil, err
		}

		fmt.Println("vString", *vString)

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

	// TODO: don't have a type checker so for right now
	// just type check in here
	for _, stmt := range n.Value.([]*builder.Node) {
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

	for i, s := range n.Value.([]*builder.Node) {
		vString, err = t.TranspileStatement(s)
		if err != nil {
			return nil, err
		}

		// `
		// {
		// 	"Type": "decl",
		// 	"Value": {
		// 		"Type": "type",
		// 		"Kind": "map",
		// 		"Value": "map"
		// 	},
		// 	"Left": {
		// 		"Type": "ident",
		// 		"Value": "m"
		// 	}
		// },
		// `

		// if s.

		// Shave off the semicolon since we don't need it
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
	fmt.Println("vsTrIng:", string(blob))

	nString += *vString

	var args = n.Metadata["args"]

	blob, _ = json.Marshal(n.Value.(*builder.Node))
	fmt.Println("n.Value.(*builder.Node):", string(blob))

	// Just do the checking here for now, not sure the merits of making the sgroup function check
	if args == nil {
		nString += "()"
		return &nString, nil
	}

	var argString string

	funcName, ok := n.Value.(*builder.Node).Value.(string)
	if ok {
		if cFuncs[n.Value.(*builder.Node).Value.(string)] {
			t.IncludeChan <- &builder.Node{
				Type: "include",
				Left: &builder.Node{
					Type:  "literal",
					Value: "stdio.h",
				},
			}
		}

		if funcName == "Println" {
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
			t.IncludeChan <- &builder.Node{
				Type: "include",
				Left: &builder.Node{
					Type:  "literal",
					Value: "unistd.h",
				},
			}
		}
	}

	vString, err = t.TranspileEGroup(args.(*builder.Node))
	if err != nil {
		return nil, err
	}
	blob, _ = json.Marshal(vString)
	fmt.Println("egroup vstring:", string(blob))

	argString += *vString

	blob, _ = json.Marshal(vString)
	fmt.Println("argstring:", string(blob))

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

func (t *Transpiler) TranspileForOfStatement(n *builder.Node) (*string, error) {
	// Change forin to be a block statement containing:
	//	- declare temp var
	//	- declare array/iter
	//	- while tempvar < iter.length
	//	- var = iter[tempvar]
	//	- loop_block
	//	-	increment var

	// return nil, errors.New("not implemented: forof")

	if n.Type != "forof" {
		return nil, errors.New("Node is not a forof")
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

func (t *Transpiler) TranspileWhileStatement(n *builder.Node) (*string, error) {
	/*
		while statements are simple, we already have all the tools:
		`while` `(` expr `)` block
	*/

	var nString string

	var v = t.Builder.ScopeTree.Get(n.Right.Value.(string))

	blob, _ := json.Marshal(v)
	fmt.Println("vvvvblobbyhill:", string(blob))

	if v != nil && v.Value != nil && v.Value.(*builder.Node).Kind == "map" {
		nString = fmt.Sprintf("for (auto const& set : %s)", n.Right.Value.(string))
		start := n.Metadata["start"]
		if start == nil {
			return nil, errors.New("No start amount ...")
			// return nil
		}

		n.Value.(*builder.Node).Value = append([]*builder.Node{{
			Type: "decl",
			Value: &builder.Node{
				Type: "type",
				// Kind: "int",
				Value: "auto",
			},
			Left: start.(*builder.Node),
			Right: &builder.Node{
				Type:  "literal",
				Kind:  "auto",
				Value: "set.first",
			},
		}}, n.Value.(*builder.Node).Value.([]*builder.Node)...)

		blob, _ = json.Marshal(n.Value.(*builder.Node).Value)
		fmt.Println("blobbyhill:", string(blob))

	} else {
		start := n.Metadata["start"]
		if start == nil {
			return nil, errors.New("No start amount ...")
			// return nil
		}

		/*
			while statements are simple, we already have all the tools:
			`while` `(` expr `)` block
		*/

		if n.Type != "while" {
			return nil, errors.New("Node is not a while")
		}

		var (
			nString = fmt.Sprintf("int %s=0;{ while(", start.(*builder.Node).Value)
			// vString *string
			// err     error
		)

		fmt.Printf("transpile expr: %+v\n", n.Left.Right)
		condition, err := t.TranspileExpression(n.Left)
		if err != nil {
			return nil, err
		}

		nString += *condition + ")"

		block, err := t.TranspileBlockStatement(n.Value.(*builder.Node))
		if err != nil {
			return nil, err
		}

		var b = *block

		b = b[:len(b)-1] + start.(*builder.Node).Value.(string) + "++;" + b[len(b)-1:]

		nString += b + "}"

		return &nString, nil
	}

	// Translate the block statement
	vString, err := t.TranspileBlockStatement(n.Value.(*builder.Node))
	if err != nil {
		return nil, err
	}

	nString += *vString

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
	fmt.Println("n", n)

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
		},
		Left: n,
		Right: &builder.Node{
			Type:  "literal",
			Value: 0,
		},
	}
}
