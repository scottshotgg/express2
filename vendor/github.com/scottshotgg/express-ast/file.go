package ast

import (
	"strings"
	"sync"
)

// File represents a file that is being compiled
type File struct {
	// TODO: think about how file and block can abstractly be the same thing
	Name string
	// This might need to include an array for functions; do this for blocks too
	Statements []Statement
}

// Length returns the list of statements in the file
func (f *File) Length() int {
	// TODO: this will have to do something to recurse down the chain and figure out blocks and add that to the total
	// return len(f.Statements)

	// for _, stmt := range f.Statements {
	// 	// TODO: statement should define a .Length() function that will return the length of the statement node
	// }

	return -1
}

// NewFile returns a new file and sets the filename
func NewFile(filename string) *File {
	return &File{
		Name: filename,
	}
}

// AddStatement appends a statement to the file
func (f *File) AddStatement(stmt Statement) {
	f.Statements = append(f.Statements, stmt)
}

func (f *File) Kind() NodeType { return FileNode }

func (f *File) String() string {
	stmts := make([]string, len(f.Statements))

	wg := sync.WaitGroup{}
	for i, stmt := range f.Statements {
		wg.Add(1)

		go func(i int, stmt Statement, wg *sync.WaitGroup) {
			stmts[i] = stmt.String()

			wg.Done()
		}(i, stmt, &wg)
	}

	wg.Wait()

	if f.Name == "main.expr" {
		// TODO: need to do something else with the imports
		return "#include <string>\n int main() {" + strings.Join(stmts, "") + "}"
	}

	return strings.Join(stmts, "\n")
}
