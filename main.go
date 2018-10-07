package main

import (
	"encoding/json"
	"fmt"

	"github.com/scottshotgg/express-ast"
	"github.com/scottshotgg/express-lex"
)

var (
	simpleTest = `
	function something() {
		stuff := "woah"
		var thing = "yeah"
	}

	ay()

	int something = 3 - 89

	some := 1 + i < 0

	for i := 0, i < 10, i++ {
		something := "here"
	}

	for i in something {
		something := "here"
		return true
	}

	{
		woah := "random string"
	}

	if something {
		int i = 0
	} else {
		int i = 1
	}
	
	char[] me`
)

// func CompressTokens(lexTokens []token.Token) ([]token.Token, error) {
// 	compressedTokens := []token.Token{}

// 	alreadyChecked := false

// 	for i := 0; i < len(lexTokens)-1; i++ {
// 		// This needs to be simplified
// 		if lexTokens[i].Type == "ASSIGN" || lexTokens[i].Type == "SEC_OP" || lexTokens[i].Type == "PRI_OP" && lexTokens[i+1].Type == "ASSIGN" || lexTokens[i+1].Type == "SEC_OP" || lexTokens[i+1].Type == "PRI_OP" {
// 			compressedToken, ok := token.TokenMap[lexTokens[i].Value.String+lexTokens[i+1].Value.String]
// 			// fmt.Println("added \"" + lexTokens[i].Value.String + lexTokens[i+1].Value.String + "\"")
// 			if ok {
// 				compressedTokens = append(compressedTokens, compressedToken)
// 				i++

// 				// If we were able to combine the last two tokens and make a new one, mark it
// 				if i == len(lexTokens)-1 {
// 					alreadyChecked = true
// 				}

// 				continue
// 			}
// 		}

// 		// Filter out the white space
// 		if lexTokens[i].Type == "WS" {
// 			continue
// 		}

// 		compressedTokens = append(compressedTokens, lexTokens[i])
// 	}

// 	// If it hasn't been already checked and the last token is not a white space, then append it
// 	if !alreadyChecked && lexTokens[len(lexTokens)-1].Type != "WS" {
// 		compressedTokens = append(compressedTokens, lexTokens[len(lexTokens)-1])
// 	}

// 	return compressedTokens, nil
// }

func main() {
	// Lex the source code
	tokens, err := lex.New(simpleTest).Lex()
	if err != nil {
		fmt.Println("err", err)
	}

	tokens, err = ast.CompressTokens(tokens)
	if err != nil {
		fmt.Println("err", err)
	}

	for _, t := range tokens {
		fmt.Println(t)
	}

	builder := &ast.ASTBuilder{
		Tokens: tokens,
	}

	p, err := builder.BuildAST()
	if err != nil {
		fmt.Println("err", err)
	}

	pJSON, _ := json.Marshal(p)
	fmt.Println()
	fmt.Println(string(pJSON))

}
