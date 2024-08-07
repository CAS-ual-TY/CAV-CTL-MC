package main

import (
	"cav/golang/parser"
	"fmt"
	"os"
)

func main() {
	for _, arg := range os.Args {
		fmt.Print(" " + arg)
	}
	fmt.Println()

	if len(os.Args) < 2 {
		fmt.Println("Usage: main <file>")
		os.Exit(1)
	}

	file := os.Args[1]

	wd, _ := os.Getwd()
	fmt.Println("Working directory: " + wd)

	_, flas, err := parser.PARSER.ParseFile(fmt.Sprintf("%s/%s", wd, file))
	if err != nil {
		fmt.Println("Failed to parse file:")
		fmt.Println(err)
		return
	}

	fmt.Println("Formula Results:")
	for _, fla := range flas {
		fmt.Println(fla.String() + ":")
		fmt.Println(fla.Check().String())
	}
}
