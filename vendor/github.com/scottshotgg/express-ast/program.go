package ast

import "fmt"

// Program represents the following form:
// [ statement ]*
type Program struct {
	// Name string // Don't know if I should include this or not
	Files map[string]*File
}

// Length returns the length of files in the program
func (p *Program) Length() int {
	return len(p.Files)
}

// NewProgram returns a new program
func NewProgram() *Program {
	return &Program{
		Files: map[string]*File{},
	}
}

// AddFile appends a file to the program
func (p *Program) AddFile(f *File) {
	p.Files[f.Name] = f
}

func (p *Program) Kind() NodeType { return ProgramNode }

func (p *Program) String() string {
	// For now:
	// - leave it like this and loop over the files
	return fmt.Sprintf("%+v", *p)
}
