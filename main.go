package main

import (
	"fmt"
	"gossg/parser"
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

	lexer := parser.NewLexer(string(content))
	var tokens []parser.Token
	for {
		token := lexer.ReadNextToken()
		tokens = append(tokens, token)
		if token.Type == parser.EOF {
			break
		}
	}

	p := parser.NewParser(tokens)
	html := p.Parse()
	fmt.Println(html)
}
