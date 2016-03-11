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
	ret := m.Run()
	os.Exit(ret)
}

func TestFileRead(t *testing.T) {
	os.Remove(STATE_FILE)

	input := NewFileInput("./testfiles/input.txt")
	defer input.Close()

	line := input.ReadLine()
	expected := "line-1"

	if(strings.Compare(line, expected) != 0){
		log.Print(input)
		t.Error("<<" + line + ">> != <<" + expected + ">>")
	}

	expected_line_number := 1
	if(input.line_number != expected_line_number){
		t.Error( fmt.Sprintf("Line Number: %d != %d", input.line_number, expected_line_number))
	}
	os.Remove(STATE_FILE)
}

func TestSaveLoadState(t *testing.T){
	os.Remove(STATE_FILE)

	input := NewFileInput("./testfiles/input.txt")
	defer input.Close()
	input.line_number = 47
	input.saveState()

	input2 := NewFileInput("./testfiles/input.txt")
	defer input2.Close()
	input2.loadState()

	if(input2.line_number != input.line_number){
		t.Error( fmt.Sprintf("Input1 %d != Input2 %d", input.line_number, input2.line_number))
	}
	os.Remove(STATE_FILE)
}

func TestResumingLines(t *testing.T){
	os.Remove(STATE_FILE)
	input := NewFileInput("./testfiles/input.txt")

	line1 := input.ReadLine()
	line2 := input.ReadLine()

	if(strings.Compare(line1, "line-1") != 0 || strings.Compare(line2, "line-2") != 0){
		t.Error("Input was not as expected")
	}

	input.Close()

	input = NewFileInput("./testfiles/input.txt")
	expected := 2
	if(input.line_number != expected){
		t.Error("Expected Line Number:", expected, "Got:", input.line_number)
	}
	os.Remove(STATE_FILE)
}