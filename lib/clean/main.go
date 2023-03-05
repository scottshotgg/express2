package main

// var cToExprTypeMap = map[string]string {
// 	"void": "",

// 	"int": "int",
// 	"double": "float",
// 	"char": "char",
// 	"long int": "int",
// }

// var typeReplaceMap = map[string]string {
// 	"void": "func",

// 	"double": "float",
// 	"long int": "int",
// 	"wint_t": "int",

// }

// func main() {
// 	var f, err = os.Open("libc_raw.txt")
// 	if err != nil {
// 		fmt.Println("something")
// 		panic(err)
// 	}

// 	var (
// 		br    = bufio.NewReader(f)
// 		lines []string
// 	)

// 	for {
// 		var l, _, err = br.ReadLine()
// 		if err != nil {
// 			if err == io.EOF {
// 				break
// 			}
// 		}

// 		var split = strings.Split(string(l), "\t")
// 		for _, s := range split {
// 			if strings.Contains(s, ";") {
// 				lines = append(lines, s)
// 				fmt.Println(s)

// 			}
// 		}
// 	}
// }
