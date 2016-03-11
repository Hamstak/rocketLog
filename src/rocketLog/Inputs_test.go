package main

import (
	"testing"
	"strings"
	"fmt"
	"log"
	"os"
)

func TestMain(m *testing.M){
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	os.Exit(m.Run())
}

func TestFileRead(t *testing.T) {
	input := NewFileInput("./testfiles/input.txt")

	log.Print(input.line_number)
	line := input.ReadLine()
	log.Print(line)
	log.Print(input.line_number)

	expected := "192.168.99.1 - - [11/Mar/2016:06:05:42 +0000] \"GET /index.html HTTP/1.1\" 304 - \"http://192.168.99.101:32773/\" \"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/48.0.2564.116 Safari/537.36\""

	if(strings.Compare(line, expected) != 0){
		log.Print(input)
		t.Error("<<" + line + ">> != <<" + expected + ">>")
	}

	expected_line_number := 1
	if(input.line_number != expected_line_number){
		t.Error( fmt.Sprintf("Line Number: %d != %d", input.line_number, expected_line_number))
	}

	input.Close()
}

func TestSaveLoadState(t *testing.T){
	input := NewFileInput("./testfiles/input.txt")
	input.line_number = 47
	input.SaveState()

	input2 := NewFileInput("./testfiles/input.txt")
	input2.LoadState()

	if(input2.line_number != input.line_number){
		t.Error( fmt.Sprintf("Input1 %d != Input2 %d", input.line_number, input2.line_number))
	}

	input.Close()
	input2.Close()

	defer os.Remove(STATE_FILE)
}