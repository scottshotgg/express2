package main

import (
	"fmt"
	"os"

	"github.com/scottshotgg/express2/compiler"
)

func main() {
	var file = "compiler/test/test.expr"

	var err = compiler.Compile(file)
	if err != nil {
		fmt.Println("error: %+v", err)
		os.Exit(9)
	}
}
