package main

import (
	"fmt"
	"rocketLog/inputs"
)

func main() {
	fmt.Println("Hello, World!");
	finput := inputs.NewFileInput("./input.txt")
	finput.Close()
}

