package lex

import (
	"io/ioutil"
	"strconv"

	"github.com/pkg/errors"

	"github.com/scottshotgg/express-token"
)

// Lexer holds all the needed variables to appropriately lex
type Lexer struct {
	source      []rune
	Accumulator string
	Tokens      []token.Token
}

// Lexemes are the specific symbols the lexer needs to recognize
var Lexemes = []string{
	"var",
	"int",
	"float",
	"string",
	"bool",
	"char",
	"object",

	":",
	"=",
	"+",
	"-",
	"*",
	"/",
	"(",
	")",
	"{",
	"}",
	"[",
	"]",
	"\"",
	"'",
	";",
	",",
	"#",
	"!",
	"<",
	">",
	"@",
	"\\",
	// "„",
	" ",
	"\n",
	"\t",

	// "select",
	// "SELECT",
	// "FROM",
	// "WHERE",
}

// New returns a new lexer attached to the provided source
func New(source string) *Lexer {
	return &Lexer{
		source: []rune(source),
	}
}

// NewFromFile returns a lexer attached to a specific file
func NewFromFile(path string) (*Lexer, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return New(string(data)), nil
}

// LexLiteral is used for determining whether something is a ident or literal
// If it is a literal is it a string, char, int, float, or bool
func (meta *Lexer) LexLiteral() (token.Token, error) {
	// Make a token and set the default value to bool; this is just because its the
	// first case in the switch and everything below sets it, so it makes the code a bit
	// cleaner
	// We COULD do this with tokens in the tokenMap for true and false
	t := token.Token{
		Type: token.Literal,
		Value: token.Value{
			True:   false,
			Type:   token.BoolType,
			String: meta.Accumulator,
		},
	}

	switch meta.Accumulator {
	// Default value is false, we only need to catch the case to keep it out of the default
	case "false":

	// Check if its true
	case "true":
		t.Value.True = true

	// Else move on and figure out what kind of number it is (or an ident)
	default:
		// Figure out from the two starting characters
		base := 10
		if len(meta.Accumulator) > 2 {
			switch meta.Accumulator[:2] {
			// Binary
			case "0b":
				base = 2

			// Octal
			case "0o":
				base = 8

			// Hex
			case "0x":
				base = 16
			}
		}

		// If the base is not 10 anymore, shave off the 0b, 0o, or 0x
		if base != 10 {
			meta.Accumulator = meta.Accumulator[2:]
		}

		// Attempt to parse an int from the accumulator
		value, err := strconv.ParseInt(meta.Accumulator, base, 64)

		// Convert the int64 to an int for now
		// I'll switch this when I'm ready to deal with different bit sizes
		t.Value.True = int(value)
		t.Value.Type = token.IntType

		// TODO: need to make something for scientific notation with carrots and e
		// If it errors, check to see if it is an float
		if err != nil {
			// Attempt to parse a float from the accumulator
			t.Value.True, err = strconv.ParseFloat(meta.Accumulator, 64)
			t.Value.Type = token.FloatType
			if err != nil {
				// If it's not a float, check whether it is a keyword
				keyword, ok := token.TokenMap[meta.Accumulator]
				if ok {
					t = keyword
				} else {
					// If it is not a keyword or a parse-able number, assume that it is an ident (for now)
					t.Type = token.Ident
					t.Value = token.Value{
						String: meta.Accumulator,
					}
				}
			}
		}
	}

	return t, nil
}

// Lex is the primary function used to lex the source string into tokens
func (meta *Lexer) Lex() ([]token.Token, error) {
	for index := 0; index < len(meta.source); index++ {
		char := string(meta.source[index])

		// Else see if it's recognized lexeme
		lexemeToken, ok := token.TokenMap[char]

		// If it is not a recognized lexeme, add it to the accumulator and move on
		if !ok {
			meta.Accumulator += char
			continue
		}

		// Filter out the comments
		switch lexemeToken.Value.Type {
		case "div":
			index++
			if index < len(meta.source)-1 {
				switch meta.source[index] {
				case '/':
					for {
						index++
						if index == len(meta.source) || meta.source[index] == '\n' {
							break
						}
					}

				case '*':
					for {
						index++
						if index == len(meta.source) || (meta.source[index] == '*' && meta.source[index+1] == '/') {
							index++
							break
						}
					}

				default:
					meta.Tokens = append(meta.Tokens, token.TokenMap[char])
				}
			}

			continue

		// Use the lexer to parse strings
		case "squote":
			fallthrough

		case "dquote":
			// If the accumulator is not empty, check it before parsing the string
			if meta.Accumulator != "" {
				ts, err := meta.LexLiteral()
				if err != nil {
					return []token.Token{}, err
				}

				meta.Tokens = append(meta.Tokens, ts)
				meta.Accumulator = ""
			}

			stringLiteral := ""

			index++
			for string(meta.source[index]) != lexemeToken.Value.String {
				// If there is an escaping backslash in the string then just increment over
				// it so that the next accumulate and increment will pickup the next char naturally
				if string(meta.source[index]) == "\\" {
					index++
				}

				stringLiteral += string(meta.source[index])

				index++
			}

			// Don't allow strings to use single quotes like JS
			stringType := token.StringType
			if lexemeToken.Value.Type == "squote" {
				if len(stringLiteral) > 1 {
					return []token.Token{}, errors.Errorf("Too many values in character literal declaration: %s", stringLiteral)
				}

				stringType = token.CharType
			}

			meta.Tokens = append(meta.Tokens, token.Token{
				ID:   0,
				Type: token.Literal,
				Value: token.Value{
					Type:   stringType,
					True:   stringLiteral,
					String: stringLiteral,
				},
			})

			continue

		case "period":
			// For now just accumulate the period and evaluate it later during parsing
			meta.Accumulator += char
			continue
		}

		// If the accumulator is not empty, check it
		if meta.Accumulator != "" {
			ts, err := meta.LexLiteral()
			if err != nil {
				return []token.Token{}, err
			}

			meta.Tokens = append(meta.Tokens, ts)
		}

		// Append the current token and reset the accumulator
		meta.Tokens = append(meta.Tokens, lexemeToken)
		meta.Accumulator = ""
	}

	// If the accumulator is not empty, check it
	if meta.Accumulator != "" {
		ts, err := meta.LexLiteral()
		if err != nil {
			return []token.Token{}, err
		}

		meta.Tokens = append(meta.Tokens, ts)
	}

	return meta.Tokens, nil
}
