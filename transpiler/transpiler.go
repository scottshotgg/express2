package transpiler

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/scottshotgg/express2/builder"
	"github.com/scottshotgg/express2/tree_flattener"
)

type Transpiler struct {
	Name         string
	ASTCloneJSON []byte
	AST          *builder.Node
	Extra        []string
	Functions    map[string]string
	Types        map[string]string
	Includes     map[string]string
	Imports      map[string]string
	Structs      []string
	GenerateMain bool
}

func emit(line string) {
	appendChan <- line
}

var (
	wg          sync.WaitGroup
	wg1         sync.WaitGroup
	funcChan    = make(chan *builder.Node, 100)
	typeChan    = make(chan *builder.Node, 100)
	structChan  = make(chan *builder.Node, 100)
	includeChan = make(chan *builder.Node, 100)
	importChan  = make(chan *builder.Node, 100)
	appendChan  = make(chan string, 5)
)

/*
	Transpile needs to work like this:
	- recurse through each statement
	- if the statement contains ANY block, then flatten on the node
*/

func appendWorker(wg *sync.WaitGroup) {
	defer wg.Done()

	var totalFile string

	for a := range appendChan {
		totalFile += a
	}

	fmt.Println("totalFile", totalFile)
}

func (t *Transpiler) functionWorker(wg *sync.WaitGroup) {
	defer wg.Done()

	var (
		function     *string
		err          error
		functionName string
	)

	for f := range funcChan {
		functionName = f.Kind

		if t.Functions[functionName] != "" {
			// FIXME: this is an error
			log.Printf("Function already declared: %+v\n", f)
			os.Exit(9)
		}

		function, err = TranspileFunctionStatement(f)
		if err != nil {
			log.Printf("Function error: %+v %+v\n", f, err)
			os.Exit(9)
		}

		t.Functions[functionName] = *function
	}
}

func New(ast *builder.Node, name string) *Transpiler {
	var t = Transpiler{
		Name:      name,
		AST:       ast,
		Functions: map[string]string{},
		Types:     map[string]string{},
		Imports:   map[string]string{},
		Includes:  map[string]string{},
	}

	// go appendWorker(&wg)

	t.ASTCloneJSON, _ = json.Marshal(ast)

	return &t
}

// This will give use some problems with multiple compilers ...

func (t *Transpiler) Transpile() (string, error) {
	// Extract the nodes
	var (
		// flattenedImports []*builder.Node
		nodes = t.AST.Value.([]*builder.Node)
		// stringP *string
		err error

		// cpp string
	)

	// Spin off workers for each type of statement

	wg1.Add(1)
	go t.functionWorker(&wg1)

	wg1.Add(1)
	go t.typeWorker(&wg1)

	wg1.Add(1)
	go t.structWorker(&wg1)

	wg.Add(1)
	go t.includeWorker(&wg)

	wg.Add(1)
	go t.importWorker(&wg)

	// Flatten the tree
	includes, err := tree_flattener.Flatten(t.AST)
	if err != nil {
		return "", err
	}

	for i := range includes {
		includeChan <- includes[i]
	}

	for i := range nodes {
		// TODO: Switch on the statement type to figure out how to process it
		// TODO: Flatten anything with a scope

		// TODO: need to put the function into the function chan here?

		switch nodes[i].Type {
		case "function":
			funcChan <- nodes[i]

		case "struct":
			structChan <- nodes[i]

		case "typedef":
			typeChan <- nodes[i]

		case "import":
			includeChan <- nodes[i]

		case "map":
			// Just transpile the statement for now
			stringP, err := TranspileStatement(nodes[i])
			if err != nil {
				fmt.Printf("err %+v\n", err)
				os.Exit(9)
				// return "", err
			}

			t.Extra = append(t.Extra, *stringP)

		default:
			return "", errors.Errorf("Node was not categorized properly: %+v\n", nodes[i])
		}
	}

	// Close the channel and alert the worker that we are done
	close(funcChan)
	close(typeChan)
	close(structChan)

	// Wait for everything to be transpiled
	wg1.Wait()

	close(importChan)
	close(includeChan)

	// Wait for everything to be transpiled
	wg.Wait()

	if t.Functions["main"] == "" {
		return "", errors.New("No main function declared")
	}

	return t.ToCpp(), nil
}

func (t *Transpiler) typeWorker(wg *sync.WaitGroup) {
	defer wg.Done()

	var (
		stringP *string
		err     error
	)

	for node := range typeChan {
		stringP, err = TranspileStatement(node)
		if err != nil {
			fmt.Printf("err %+v\n", err)
			os.Exit(9)
			// return "", err
		}

		t.Types[node.Left.Value.(string)] = *stringP
	}
}

func (t *Transpiler) structWorker(wg *sync.WaitGroup) {
	defer wg.Done()

	var (
		stringP *string
		err     error
	)

	for node := range structChan {
		stringP, err = TranspileStatement(node)
		if err != nil {
			fmt.Printf("err %+v\n", err)
			os.Exit(9)
			// return "", err
		}

		t.Structs = append(t.Structs, *stringP)
	}
}

func (t *Transpiler) includeWorker(wg *sync.WaitGroup) {
	defer wg.Done()

	var (
		includeStringP *string
		// Why does this shadow ...
		// Is the gofunc "capturing" variables that aren't passed?
		ierr error
	)

	for node := range includeChan {
		// Might want to make this go through the entire pipeline ...
		includeStringP, ierr = TranspileIncludeStatement(node)
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

func (t *Transpiler) importWorker(wg *sync.WaitGroup) {
	defer wg.Done()

	var (
		importStringP *string
		// Why does this shadow ...
		// Is the gofunc "capturing" variables that aren't passed?
		ierr error
	)

	for node := range importChan {
		// Might want to make this go through the entire pipeline ...
		importStringP, ierr = TranspileImportStatement(node)
		if ierr != nil {
			log.Printf("Error transpiling import statement: %+v\n", ierr)

			// Exit if there is a problem transpiling the import statement
			// and we'll deal with it later
			os.Exit(9)
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

	if len(output) > 0 {
		output = append(output, "\n")
	}

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

	if len(typesString) == len(typesString) {
		typesString += "// none\n"
	}
	var structsString = "\n\n// Structs:\n"

	for _, t := range t.Structs {
		structsString += t + "\n"
	}

	if len(structsString) == len(structsString) {
		structsString += "// none\n"
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

	return "\n// Prototypes:\n" + strings.Join(prototypes, "\n") +
		"\n\n// Functions:" + functionString +
		fmt.Sprintf("\n// Main:\n// generated: %v\n%s", t.GenerateMain, mainFunc)
}

func TranspileTypeDeclaration(n *builder.Node) (*string, error) {
	// Format should be:
	// `type` [ident] `=` [type]
	// Left is the ident
	// Right is the type

	// TODO: gonna have to do something to actually enable this type in the parser/compiler

	if n.Type != "typedef" {
		return nil, errors.New("Node is not a typedef")
	}

	var nString = "typedef "

	var cpp, err = TranspileType(n.Right)
	if err != nil {
		return nil, err
	}

	nString += *cpp + " "

	// This will allow technically allow idents to be made from general expressions; not sure if we should keep this or not
	// Might have to change it to TranspileIdent
	cpp, err = TranspileExpression(n.Left)
	if err != nil {
		return nil, err
	}

	nString += *cpp + ";"

	return &nString, nil
}

func TranspileStructDeclaration(n *builder.Node) (*string, error) {
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
	var vString, err = TranspileExpression(n.Left)
	if err != nil {
		return nil, err
	}

	// Could just have it add `struct` here but this will show us changes
	var nString = n.Type + " " + *vString

	// Transpile the block for the value
	vString, err = TranspileBlockStatement(n.Right)
	if err != nil {
		return nil, err
	}

	nString += *vString + ";"

	return &nString, nil
}

func TranspileIncludeStatement(n *builder.Node) (*string, error) {
	if n.Type != "include" {
		return nil, errors.New("Node is not an include")
	}

	lhs, err := TranspileExpression(n.Left)
	if err != nil {
		return nil, err
	}

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

func TranspileImportStatement(n *builder.Node) (*string, error) {
	if n.Type != "import" {
		return nil, errors.New("Node is not an inc")
	}

	lhs, err := TranspileExpression(n.Left)
	if err != nil {
		return nil, err
	}

	*lhs = "#include<" + *lhs + ">"

	return lhs, nil
}

func TranspileIncrementExpression(n *builder.Node) (*string, error) {
	if n.Type != "inc" {
		return nil, errors.New("Node is not an inc")
	}

	lhs, err := TranspileExpression(n.Left)
	if err != nil {
		return nil, err
	}

	// Put parenthesis around it
	*lhs = "(" + *lhs + ")++"

	return lhs, nil
}

func TranspileIndexExpression(n *builder.Node) (*string, error) {
	/*
		Left is an expression
		Right is an expression
	*/

	if n.Type != "index" {
		return nil, errors.New("Node is not an index")
	}

	lhs, err := TranspileExpression(n.Left)
	if err != nil {
		return nil, err
	}

	rhs, err := TranspileExpression(n.Right)
	if err != nil {
		return nil, err
	}

	var nString = *lhs + "[" + *rhs + "]"

	return &nString, nil
}

func TranspileExpression(n *builder.Node) (*string, error) {
	switch n.Type {

	case "literal":
		return TranspileLiteralExpression(n)

	case "ident":
		return TranspileIdentExpression(n)

	case "comp":
		return TranspileConditionExpression(n)

	case "binop":
		return TranspileBinOpExpression(n)

	case "array":
		return TranspileArrayExpression(n)

	case "call":
		return TranspileCallExpression(n)

	case "index":
		return TranspileIndexExpression(n)

	case "block":
		log.Println("blockExpression", n)
		return TranspileBlockExpression(n)

		// case "selection":
		// 	return TranspileSelectExpression(n)
	}

	return nil, errors.Errorf("Not implemented expression: %+v", n)
}

func TranspileMapStatement(n *builder.Node) (*string, error) {
	if n.Type != "map" {
		return nil, errors.New("Node is not a map")
	}

	// Transpile the ident
	var vString, err = TranspileExpression(n.Left)
	if err != nil {
		return nil, err
	}

	// Could just have it add `struct` here but this will show us changes
	var nString = "std::map<std::string, std::string>" + " " + *vString + "= "

	// Transpile the block for the value
	vString, err = TranspileMapBlockStatement(n.Right)
	if err != nil {
		return nil, err
	}

	// Include std::map from C++
	includeChan <- &builder.Node{
		Type: "include",
		Left: &builder.Node{
			Type:  "literal",
			Value: "map",
		},
	}

	nString += *vString + ";"

	return &nString, nil
}

func TranspileLaunchStatement(n *builder.Node) (*string, error) {
	if n.Type != "launch" {
		return nil, errors.New("Node is not a launch node")
	}

	// Transpile the ident
	var vString, err = TranspileStatement(n.Left)
	if err != nil {
		return nil, err
	}

	// Include libmill for coroutines
	// includeChan <- &builder.Node{
	// 	Type: "include",
	// 	Kind: "path",
	// 	Left: &builder.Node{
	// 		Type:  "literal",
	// 		Value: "../lib/libmill/libmill.h",
	// 	},
	// }
	includeChan <- &builder.Node{
		Type: "include",
		// Kind: "path",
		Left: &builder.Node{
			Type:  "literal",
			Value: "libmill.h",
		},
	}

	var nString = "go([=](...){" + *vString + "}());"

	return &nString, nil
}

func TranspileEnumBlockStatement(n *builder.Node) (*string, error) {
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

		vString, err = TranspileStatement(stmt)
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

func TranspileEnumStatement(n *builder.Node) (*string, error) {
	if n.Type != "enum" {
		return nil, errors.New("Node is not a map")
	}

	var enum, err = TranspileEnumBlockStatement(n.Left)
	if err != nil {
		return nil, err
	}

	var nString = "enum " + *enum + ";"

	return &nString, err
}

func TranspileKeyValueStatement(n *builder.Node) (*string, error) {
	var left, err = TranspileExpression(n.Left)
	if err != nil {
		return nil, err
	}

	right, err := TranspileExpression(n.Right)
	if err != nil {
		return nil, err
	}

	var nString = "{" + *left + "," + *right + "}"

	return &nString, nil
}

func TranspileDeferStatement(n *builder.Node) (*string, error) {
	var stmt, err = TranspileStatement(n.Left)
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

func TranspileStatement(n *builder.Node) (*string, error) {
	switch n.Type {

	case "launch":
		return TranspileLaunchStatement(n)

	case "defer":
		return TranspileDeferStatement(n)

	case "enum":
		return TranspileEnumStatement(n)

	case "kv":
		return TranspileKeyValueStatement(n)

	case "map":
		return TranspileMapStatement(n)

	case "typedef":
		return TranspileTypeDeclaration(n)

	case "struct":
		return TranspileStructDeclaration(n)

	// FIXME: Why do we have expressions in here ... ?
	case "literal":
		return TranspileLiteralExpression(n)

	case "inc":
		var cppString, err = TranspileIncrementExpression(n)
		if err == nil {
			*cppString += ";"
		}

		return cppString, err

	case "call":
		var cppString, err = TranspileCallExpression(n)
		if err == nil {
			*cppString += ";"
		}

		return cppString, err

	case "ident":
		return TranspileIdentExpression(n)

	case "import":
		return TranspileImportStatement(n)

	case "include":
		return nil, errors.Errorf("Direct C/C++ usage is not implemented yet: include: %+v\n", n)
		// return TranspileIncludeStatement(n)

	case "assignment":
		return TranspileAssignmentStatement(n)

	case "decl":
		var stmt, err = TranspileDeclarationStatement(n)
		fmt.Println("stmt, err", *stmt, err)
		return stmt, err

	case "function":
		funcChan <- n
		// Right now just grab the name from the string
		// Later on we can issue new function names for lambdas
		var name = n.Kind
		return &name, nil
		// function, err = TranspileFunctionStatement(n)

	case "return":
		return TranspileReturnStatement(n)

	case "block":
		return TranspileBlockStatement(n)

	case "while":
		return TranspileWhileStatement(n)

	case "forof":
		return TranspileForOfStatement(n)

	case "forin":
		return TranspileForInStatement(n)

	case "forstd":
	}

	return nil, errors.Errorf("Not implemented statement: %+v", n)
}

func TranspileReturnStatement(n *builder.Node) (*string, error) {
	if n.Type != "return" {
		return nil, errors.New("Node is not a return statement")
	}

	// Return statments come in the form `return` { expr }

	var nString = "return"

	fmt.Printf("n: %+v\n", n)

	// LHS (the return expression) is allowed to be empty
	if n.Left != nil {
		exprString, err := TranspileExpression(n.Left)
		if err != nil {
			return nil, err
		}

		nString += " " + *exprString
	}

	nString += ";"

	return &nString, nil
}

func TranspileFunctionStatement(n *builder.Node) (*string, error) {
	if n.Type != "function" {
		return nil, errors.New("Node is not an function")
	}

	// Include std::map from C++
	includeChan <- &builder.Node{
		Type: "include",
		Kind: "path",
		Left: &builder.Node{
			Type:  "literal",
			Value: "../lib/defer.cpp",
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
	argsString, err := TranspileSGroup(n.Metadata["args"].(*builder.Node))
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

	if returns != nil {
		// returns is a `type` for now; multiple returns are not supported right now
		returnsStringP, err = TranspileType(returns.(*builder.Node))
		if err != nil {
			return nil, err
		}
	}

	// Prepend the return string with a space
	nString = *returnsStringP + " " + nString

	blockString, err := TranspileBlockStatement(n.Value.(*builder.Node))
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

func TranspileIdentExpression(n *builder.Node) (*string, error) {
	if n.Type != "ident" {
		return nil, errors.New("Node is not an ident")
	}

	nString, ok := n.Value.(string)
	if !ok {
		return nil, errors.Errorf("Node value was not a string; %v", n)
	}

	return &nString, nil
}

func TranspileType(n *builder.Node) (*string, error) {
	if n.Type != "type" {
		return nil, errors.New("Node is not an type")
	}

	nString, ok := n.Value.(string)
	if !ok {
		return nil, errors.Errorf("Node value was not a string; %v", n)
	}

	if nString == "string" {
		nString = "std::" + nString

		// TODO: switch this to just use a damn string later
		includeChan <- &builder.Node{
			Type: "include",
			Left: &builder.Node{
				Type:  "literal",
				Value: "string",
			},
		}
	}

	return &nString, nil
}

// This changes an Express literal to be formatted the way C++ expects
func prepLiteral(kind, cpp string) *string {
	switch kind {
	case "string":
		cpp = "\"" + cpp + "\""
	}

	return &cpp
}

func TranspileLiteralExpression(n *builder.Node) (*string, error) {
	if n.Type != "literal" {
		return nil, errors.New("Node is not an literal")
	}

	return prepLiteral(n.Kind, fmt.Sprintf("%v", n.Value)), nil
}

func TranspileArrayExpression(n *builder.Node) (*string, error) {
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
		vString, err = TranspileExpression(v)
		if err != nil {
			return nil, err
		}

		nString += *vString + ", "
	}

	// Cut off the last comma and space
	nString = nString[:len(nString)-2] + " }"

	return &nString, nil
}

func TranspileAssignmentStatement(n *builder.Node) (*string, error) {
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
	vString, err = TranspileExpression(n.Left)
	if err != nil {
		return nil, err
	}

	nString = *vString + " = "

	// Translate the ident expression (lhs)
	vString, err = TranspileExpression(n.Right)
	if err != nil {
		return nil, err
	}

	nString += *vString + ";"

	return &nString, nil
}

func TranspileDeclarationStatement(n *builder.Node) (*string, error) {
	if n.Type != "decl" {
		return nil, errors.New("Node is not an declaration")
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

	vString, err = TranspileType(n.Value.(*builder.Node))
	if err != nil {
		return nil, err
	}

	log.Println("TYPE", *vString, n.Left)

	nString = *vString + " "

	// LHS is not allowed to be nil
	if n.Left == nil {
		return nil, errors.New("nil Left hand side")
	}

	// Translate the ident expression (lhs)
	vString, err = TranspileExpression(n.Left)
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

		// log.Printf("HEY ITS ME %+v\n", n.Value.(*builder.Node))

		nString += ";"
		return &nString, nil
	}

	// Translate the ident expression (lhs)
	vString, err = TranspileExpression(n.Right)
	if err != nil {
		return nil, err
	}

	nString += " = " + *vString + ";"

	// fmt.Println("nString", nString)

	return &nString, nil
}

// func TranspileIncrementExpression(n *builder.Node) (*string, error) {
// 	if n.Type != "inc" {
// 		return nil, errors.New("Node is not an inc")
// 	}

// 	var (
// 		nString = ""
// 		vString *string
// 		err     error
// 	)

// 	// Translate the ident expression (lhs)
// 	vString, err = TranspileExpression(n.Left)
// 	if err != nil {
// 		return nil, err
// 	}

// 	nString += *vString + "++;"

// 	return &nString, nil
// }

func TranspileConditionExpression(n *builder.Node) (*string, error) {
	if n.Type != "comp" {
		return nil, errors.New("Node is not an comp")
	}

	var (
		nString = ""
		vString *string
		err     error
	)

	// Translate the lhs
	vString, err = TranspileExpression(n.Left)
	if err != nil {
		return nil, err
	}

	nString += *vString

	// Translate the rhs
	vString, err = TranspileExpression(n.Right)
	if err != nil {
		return nil, err
	}

	nString += n.Value.(string) + *vString

	return &nString, nil
}

func TranspileBinOpExpression(n *builder.Node) (*string, error) {
	if n.Type != "binop" {
		return nil, errors.New("Node is not a binop")
	}

	var (
		nString = ""
		vString *string
		err     error
	)

	// Translate the ident expression (lhs)
	vString, err = TranspileExpression(n.Left)
	if err != nil {
		return nil, err
	}

	nString += *vString + n.Value.(string)

	// Translate the ident expression (lhs)
	vString, err = TranspileExpression(n.Right)
	if err != nil {
		return nil, err
	}

	nString += *vString

	return &nString, nil
}

func TranspileBlockStatement(n *builder.Node) (*string, error) {
	if n.Type != "block" {
		return nil, errors.New("Node is not a block")
	}

	var (
		nString = ""
		vString *string
		err     error
	)

	for _, stmt := range n.Value.([]*builder.Node) {
		fmt.Println("stmt", stmt)

		vString, err = TranspileStatement(stmt)
		if err != nil {
			return nil, err
		}

		fmt.Println("vString", *vString)

		if stmt.Type != "function" {
			nString += *vString
		}
	}

	nString = "{" + nString + "}"

	return &nString, nil
}

func TranspileMapBlockStatement(n *builder.Node) (*string, error) {
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

		vString, err = TranspileStatement(stmt)
		if err != nil {
			return nil, err
		}

		nString += *vString + ","
	}

	nString = "{" + nString + "}"

	return &nString, nil
}

// TODO: for now all variables will have a `.` in front, later on it
// should only be the public variables
func TranspileBlockExpression(n *builder.Node) (*string, error) {
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
		vString, err = TranspileStatement(stmt)
		if err != nil {
			return nil, err
		}

		if stmt.Type != "assignment" {
			return nil, errors.Errorf("Structs can only contain assignment statements; %+v", stmt)
		}

		*vString = (*vString)[:len(*vString)-1]

		// Add a dot in front
		nString += "." + *vString + ","
	}

	nString = "{" + nString + "}"

	return &nString, nil
}

func TranspileEGroup(n *builder.Node) (*string, error) {
	if n.Type != "egroup" {
		return nil, errors.New("Node is not a egroup")
	}

	var (
		nStrings = make([]string, len(n.Value.([]*builder.Node)))
		vString  *string
		err      error
	)

	for i, e := range n.Value.([]*builder.Node) {
		vString, err = TranspileExpression(e)
		if err != nil {
			return nil, err
		}

		nStrings[i] = *vString
	}

	var nString = "(" + strings.Join(nStrings, ",") + ")"

	return &nString, nil
}

func TranspileSGroup(n *builder.Node) (*string, error) {
	if n.Type != "sgroup" {
		return nil, errors.New("Node is not a sgroup")
	}

	var (
		nStrings = make([]string, len(n.Value.([]*builder.Node)))
		vString  *string
		err      error
	)

	for i, s := range n.Value.([]*builder.Node) {
		vString, err = TranspileStatement(s)
		if err != nil {
			return nil, err
		}

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

func TranspileCallExpression(n *builder.Node) (*string, error) {
	if n.Type != "call" {
		return nil, errors.New("Node is not a call")
	}

	var (
		nString = ""
		vString *string
		err     error
	)

	vString, err = TranspileIdentExpression(n.Value.(*builder.Node))
	if err != nil {
		return nil, err
	}

	nString += *vString

	var args = n.Metadata["args"]

	// Just do the checking here for now, not sure the merits of making the sgroup function check
	if args != nil {
		vString, err = TranspileEGroup(args.(*builder.Node))
		if err != nil {
			return nil, err
		}

		nString += *vString
	} else {
		nString += "()"
	}

	return &nString, nil
}

func TranspileForInStatement(n *builder.Node) (*string, error) {
	// Change forin to be a block statement containing:
	//	- declare var
	//	- declare array/iter
	//	- while var < iter.length
	//	- loop_block
	//	-	increment var

	// return nil, errors.New("not implemented: forin")

	if n.Type != "forin" {
		return nil, errors.New("Node is not a forin")
	}

	var (
		nString = "{"
		vString *string
		err     error
	)

	// Make and translate the ident into a declaration
	ds := TransformIdentToDefaultDeclaration(n.Metadata["start"].(*builder.Node))
	vString, err = TranspileDeclarationStatement(ds)
	if err != nil {
		return nil, err
	}

	nString += *vString

	// Make and translate the array expression into a declaration
	dss := TransformExpressionToDeclaration(n.Metadata["end"].(*builder.Node))
	vString, err = TranspileDeclarationStatement(dss)
	if err != nil {
		return nil, err
	}

	nString += *vString

	// Make and translate the less than operation
	vString, err = TranspileConditionExpression(&builder.Node{
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
	vString, err = TranspileBlockStatement(n.Value.(*builder.Node))
	if err != nil {
		return nil, err
	}

	nString += *vString

	// Lastly, make and translate an increment statement for the ident
	vString, err = TranspileIncrementExpression(&builder.Node{
		Type: "inc",
		Left: n.Metadata["start"].(*builder.Node),
	})
	if err != nil {
		return nil, err
	}

	nString += *vString + "}"

	return &nString, nil
}

func TranspileForOfStatement(n *builder.Node) (*string, error) {
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
	vString, err = TranspileDeclarationStatement(ds)
	if err != nil {
		return nil, err
	}

	nString += *vString

	// Make and translate the array expression into a declaration
	dss := TransformExpressionToDeclaration(n.Metadata["end"].(*builder.Node))
	vString, err = TranspileDeclarationStatement(dss)
	if err != nil {
		return nil, err
	}

	nString += *vString

	// Make and translate the less than operation
	vString, err = TranspileConditionExpression(&builder.Node{
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
	vString, err = TranspileBlockStatement(n.Value.(*builder.Node))
	if err != nil {
		return nil, err
	}

	nString += *vString

	// Lastly, make and translate an increment statement for the ident
	vString, err = TranspileIncrementExpression(&builder.Node{
		Type: "inc",
		Left: n.Metadata["start"].(*builder.Node),
	})
	if err != nil {
		return nil, err
	}

	nString += *vString + "}"

	return &nString, nil
}

func TranspileWhileStatement(n *builder.Node) (*string, error) {
	/*
		while statements are simple, we already have all the tools:
		`while` `(` expr `)` block
	*/

	if n.Type != "while" {
		return nil, errors.New("Node is not a forof")
	}

	var (
		nString = "{ while("
		// vString *string
		// err     error
	)

	fmt.Printf("transpile expr: %+v\n", n.Left.Right)
	condition, err := TranspileExpression(n.Left)
	if err != nil {
		return nil, err
	}

	nString += *condition + ")"

	block, err := TranspileBlockStatement(n.Value.(*builder.Node))
	if err != nil {
		return nil, err
	}

	nString += *block + "}"

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
