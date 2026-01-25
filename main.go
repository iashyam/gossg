package main

import (
	"bufio"
	"fmt"
	"os"
)

func readLines(filepath string) ([]string, error) {
	// Read the file and return the lines as a slice of strings

	file, err := os.Open(filepath)
	lines := []string{}
	if err != nil {
		return lines, err
	}
	defer file.Close()

	newScanner := bufio.NewScanner(file)
	for newScanner.Scan() {
		line := newScanner.Text()
		lines = append(lines, line)
	}
	return lines, nil
}

func main() {
	// FILE := "README.md"
	// lines, err := readLines(FILE)
	// if err!=nil{
	// 	fmt.Printf("got error %s", err)
	// }
	// for _, line := range(lines){
	// 	line = headings(line)
	// 	fmt.Println(line)
	// }
	text := "start text *italics* `x=0`**emphasis** hello world *** this is nested bold italics***"
	fmt.Println(text)
	lex := NewLexer(text)
	token := lex.ReadNextToken()
	for token.Type != 0 {
		fmt.Println(token)
		token = lex.ReadNextToken()
	}
}
