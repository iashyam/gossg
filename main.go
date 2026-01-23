package main

import (
	"bufio"
	"fmt"
	"os"
)

func readLines(filepath string)([]string, error){
	// Read the file and return the lines as a slice of strings

	file, err := os.Open(filepath)
	lines := []string{} 
	if err!=nil{
		return lines, err
	}
	defer file.Close()

	newScanner := bufio.NewScanner(file)
	for newScanner.Scan(){
		line := newScanner.Text()
		lines = append(lines, line)
	}
	return lines, nil
}

func main(){
	FILE := "README.md"
	lines, err := readLines(FILE)
	if err!=nil{
		fmt.Printf("got error %s", err)
	}
	for _, line := range(lines){
		fmt.Println(line)
	}
}
