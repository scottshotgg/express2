package main

import (
	"fmt"

	"github.com/scottshotgg/express2/compiler"
)

func main() {
	var file = "compiler/test/test.expr"

	var err = compiler.Compile(file)
	if err != nil {
		fmt.Printf("\nerror: %s\n", err)
	}
}
