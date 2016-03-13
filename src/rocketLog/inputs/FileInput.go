package inputs

import (
	"bufio"
	"os"
	"log"
	"encoding/json"
	"path/filepath"
	"errors"
)

type FileInput struct {
	scanner		bufio.Scanner
	line_number	int
	abs_path	string
	state_file	string
	file		*os.File
}

func NewFileInput(path, state_file string) *FileInput{
	file, err := os.Open(path)
	if(err != nil){
		log.Fatal(err)
	}

	abs_path, err := filepath.Abs(path)
	if(err != nil){
		log.Fatal(err)
	}

	file_scanner := bufio.NewScanner(file)
	fin := &FileInput{
		scanner: *file_scanner,
		abs_path: abs_path,
		file: file,
		state_file: state_file,
	}

	file_state := fin.loadState()
	fin.SkipTo(file_state[fin.abs_path])

	return fin
}

func (input *FileInput) SkipTo(skip_to int){
	for input.line_number < skip_to {
		input.ReadLine()
	}
}

func (input *FileInput) Close() {
	input.file.Close()
	input.saveState()
}

func (input *FileInput) ReadLine() (string, error) {
	var err error
	err = nil

	if (input.scanner.Scan() == false){
		err = errors.New("No tokens left")
	} else {
		input.line_number++
	}

	return input.scanner.Text(), err
}

func (input *FileInput) saveState(){
	file_map := input.loadState()
	file_map[input.abs_path] = input.line_number

	file, err := os.OpenFile(input.state_file, os.O_RDWR | os.O_TRUNC, 0666)
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

func (input *FileInput) loadState() map[string] int{
	var file *os.File
	var err error

	file, err = os.OpenFile(input.state_file, os.O_CREATE | os.O_RDWR, 0666)
	if(err != nil){
		log.Fatal(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	file_state := make(map[string]int)
	decoder.Decode(&file_state)
	return file_state
}