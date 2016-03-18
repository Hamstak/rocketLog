package inputs

import (
	"testing"
	"strings"
	"os"
	"time"
)

const STATE_FILE = "./testfiles/state.json"
const STREAM_FILE = "testfiles/file-stream.txt"

func truncate_state(){
	os.Truncate(STATE_FILE, 0)
}

func assertSame(expected, actual string, t *testing.T){
	if(strings.Compare(expected, actual) != 0){
		t.Errorf("Expected: <<%s>>, Actual: <<%s>>", expected, actual)
	}
}

func Test_FileInputStream_ReadLine(t *testing.T) {
	truncate_state()
	file, err := os.OpenFile(STREAM_FILE, os.O_CREATE | os.O_RDWR | os.O_TRUNC, 0644)
	if(err != nil){
		t.Fatal(err)
	}

	file.WriteString("Line 1\n")

	fis := NewFileInputStream(STREAM_FILE, "test", NewFileState(STATE_FILE))

	line1, _ := fis.ReadLine()
	assertSame("Line 1", line1, t)

	file.WriteString("Line 2\n")
	line2, _ := fis.ReadLine()
	assertSame("Line 2", line2, t)

	fis.Close()
}

func Test_FileInputStream_ReadLine_AsyncWait(t *testing.T){
	truncate_state()
	file, err := os.OpenFile(STREAM_FILE, os.O_CREATE | os.O_RDWR | os.O_TRUNC, 0644)
	if(err != nil){
		t.Fatal(err)
	}

	file.WriteString("Line 1\n")

	fis := NewFileInputStream(STREAM_FILE, "test", NewFileState(STATE_FILE))

	line1, _ := fis.ReadLine()
	assertSame("Line 1", line1, t)

	go func(){
		time.Sleep(time.Second)
		file.WriteString("Line 2\n")
	}()

	line2, _ := fis.ReadLine()
	assertSame("Line 2", line2, t)

	fis.Close()
}


func Test_FileInputStream_ReadLine_Skip(t *testing.T){
	truncate_state()
	file, err := os.OpenFile(STREAM_FILE, os.O_CREATE | os.O_RDWR | os.O_TRUNC, 0644)
	if(err != nil){
		t.Fatal(err)
	}

	file.WriteString("Line 1\n")
	file.WriteString("Line 2\n")

	file_input_stream_1 := NewFileInputStream(STREAM_FILE, "test", NewFileState(STATE_FILE))
	line1, _ := file_input_stream_1.ReadLine()
	assertSame("Line 1", line1, t)
	file_input_stream_1.Close()

	file_input_stream_2 := NewFileInputStream(STREAM_FILE, "test", NewFileState(STATE_FILE))
	line2, _ := file_input_stream_2.ReadLine()
	assertSame("Line 2", line2, t)
	file_input_stream_2.Close()
}