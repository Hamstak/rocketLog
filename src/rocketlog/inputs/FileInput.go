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
	elastic_type	string
}

func NewFileInput(path, state_file, elastic_type string) *FileInput{
	abs_path, err := filepath.Abs(path)
	if(err != nil){
		log.Fatal(err)
	}

	file, err := os.OpenFile(abs_path, os.O_RDONLY, 0666)
	if(err != nil){
		log.Fatal(err)
	}

	file_scanner := bufio.NewScanner(file)
	fin := &FileInput{
		scanner: *file_scanner,
		abs_path: abs_path,
		file: file,
		state_file: state_file,
		elastic_type: elastic_type,
	}

	file_state := fin.loadState()
	fin.skipTo(file_state[fin.abs_path])

	return fin
}

func (input *FileInput) skipTo(skip_to int) error{
	for input.line_number < skip_to {
		_, err := input.ReadLine()
		if(err != nil){
			return err
		}
	}

	return nil
}

func (self *FileInput) Flush(){
	self.file.Close()
	self.saveState()

	var err error
	self.file, err = os.OpenFile(self.abs_path, os.O_RDONLY, 0666)
	if(err != nil){
		log.Fatal(err)
	}

	last_line_number := self.line_number
	self.line_number = 0
	self.scanner = *bufio.NewScanner(self.file)

	if(self.skipTo(last_line_number) != nil){
		self.line_number = 0
		self.Flush()
	}

}

func (input *FileInput) Close() {
	input.file.Close()
	input.saveState()
}

func (self *FileInput) GetType() string {
	return self.elastic_type
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