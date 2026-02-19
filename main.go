package main

import (
	"fmt"
	"os"
)

func main() {
	fileName := "test.md"
	if len(os.Args) > 1 {
		fileName = os.Args[1]
	}

	content, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	lexer := NewLexer(string(content))
	var tokens []Token
	for {
		token := lexer.ReadNextToken()
		tokens = append(tokens, token)
		if token.Type == EOF {
			break
		}
	}

	parser := NewParser(tokens)
	html := parser.Parse()
	fmt.Println(html)
}
