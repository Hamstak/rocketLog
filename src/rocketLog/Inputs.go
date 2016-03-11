package main

import (
	"bufio"
	"os"
	"log"
	"encoding/json"
	"path/filepath"
)

const STATE_FILE = "./state.json"

type FileInput struct {
	scanner     bufio.Scanner
	line_number int
	abs_path    string
	file *os.File
}


func NewFileInput(path string) FileInput {
	file, err := os.Open(path)
	if(err != nil){
		log.Fatal(err)
	}

	abs_path, err := filepath.Abs(path)
	if(err != nil){
		log.Fatal(err)
	}

	file_scanner := bufio.NewScanner(file)
	fin := FileInput{
		scanner: *file_scanner,
		line_number: int(0),
		abs_path: abs_path,
		file: file,
	}

	return fin
}

func (input FileInput) Close() {
	input.file.Close()
}

func (input FileInput) HasLine() bool {
	return input.scanner.Scan()
}

func (input FileInput) ReadLine() string {
	if (input.scanner.Scan() == false){
		log.Fatal("No Tokens Left")
	}

	input.line_number += 1
	return input.scanner.Text()
}

func (input FileInput) SaveState(){
	file_map := input.LoadState()
	file_map[input.abs_path] = input.line_number

	file, err := os.OpenFile(STATE_FILE, os.O_RDWR | os.O_TRUNC, 0666)
	if(err != nil){
		log.Fatal(err)
	}

	e := json.NewEncoder(file)
	err = e.Encode(file_map)
	if(err != nil){
		log.Fatal(err)
	}

	file.Close()
}

func (input FileInput) LoadState() map[string] int {
	var file *os.File
	var err error

	file, err = os.OpenFile(STATE_FILE, os.O_CREATE | os.O_RDWR, 0666)
	if(err != nil){
		log.Fatal(err)
	}

	decoder := json.NewDecoder(file)
	file_state := make(map[string]int)
	decoder.Decode(&file_state)

	file.Close()
	return file_state
}

