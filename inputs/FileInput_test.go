package inputs;

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

func tryDelete(filename string){
	err := os.Remove(filename)
	if(err != nil){
		log.Print(err)
	}
}

func Test_FileInput_FileRead(t *testing.T) {
	state_file := "./testfiles/test-file-read.json"
	tryDelete(state_file)

	input := NewFileInput("./testfiles/input.txt", state_file, "")
	
	line, err := input.ReadLine()
	if (err != nil){
		log.Print("No tokens found")
	}
	expected := "line-1"

	if(strings.Compare(line, expected) != 0){
		t.Error("<<" + line + ">> != <<" + expected + ">>")
	}

	expected_line_number := 1
	if(input.line_number != expected_line_number){
		t.Error( fmt.Sprintf("Line Number: %d != %d", input.line_number, expected_line_number))
	}

	input.Close()
	tryDelete(state_file)
}

func Test_FileInput_SaveLoadState(t *testing.T){
	state_file := "./testfiles/test-save-load-state.json"
	tryDelete(state_file)

	input := NewFileInput("./testfiles/input.txt", state_file, "")
	input.line_number = 3
	input.Close()

	input2 := NewFileInput("./testfiles/input.txt", state_file, "")
	input2.loadState()
	input2.Close()

	if(input2.line_number != input.line_number){
		t.Error( fmt.Sprintf("Input1 %d != Input2 %d", input.line_number, input2.line_number))
	}
	tryDelete(state_file)
}

func Test_FileInput_ResumingLines(t *testing.T){
	state_file := "./testfiles/test-resuming-lines.json"
	tryDelete(state_file)

	// Read some code and then save the state
	input := NewFileInput("./testfiles/input.txt", state_file, "")

	line1, err1 := input.ReadLine()
	line2, err2 := input.ReadLine()

	if( err1 != nil){
		log.Print(err1)
	}

	if(err2 != nil){
		log.Print(err2)
	}

	if(strings.Compare(line1, "line-1") != 0 || strings.Compare(line2, "line-2") != 0){
		t.Error("Input was not as expected")
	}

	input.Close()

	// Load the file state
	input = NewFileInput("./testfiles/input.txt", state_file, "")
	expected := 2
	if(input.line_number != expected){
		t.Error("Expected Line Number:", expected, "Got:", input.line_number)
	}

	tryDelete(state_file)
}