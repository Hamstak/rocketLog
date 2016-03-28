package inputs

import (
	"os"
	"strings"
	"testing"
	"time"
)

const stateFile = "./testfiles/state.json"
const streamFile = "testfiles/file-stream.txt"
const streamFile2 = "testfiles/file-stream-2.txt"

func truncateState() {
	os.Truncate(stateFile, 0)
}

func assertSame(expected, actual string, t *testing.T) {
	if strings.Compare(expected, actual) != 0 {
		t.Errorf("Expected: <<%s>>, Actual: <<%s>>", expected, actual)
	}
}

func Test_FileInputStream_GetName(t *testing.T) {
	truncateState()
    fis := NewFileInputStream(streamFile, "test", NewFileState(stateFile))
	assertSame("FileInputStream='"+streamFile+"'", fis.GetName(), t)
}

func Test_FileInputStream_Realloc(t *testing.T){
    truncateState()
    fis := NewFileInputStream(streamFile2, "test", NewFileState(stateFile))
    defer fis.Close()
    
    line, _ := fis.ReadLine()
   
    if(len(line) < 1400){
        t.Error("Line Too Short! ", line)
    }
}

func Test_FileInputStream_ReadLine(t *testing.T) {
	truncateState()
	file, err := os.OpenFile(streamFile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		t.Fatal(err)
	}

	file.WriteString("Line 1\n")

	fis := NewFileInputStream(streamFile, "test", NewFileState(stateFile))

	line1, _ := fis.ReadLine()
	assertSame("Line 1", line1, t)

	file.WriteString("Line 2\n")
	line2, _ := fis.ReadLine()
	assertSame("Line 2", line2, t)

	fis.Close()
}

func Test_FileInputStream_ReadLine_AsyncWait(t *testing.T) {
	truncateState()
	file, err := os.OpenFile(streamFile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		t.Fatal(err)
	}

	file.WriteString("Line 1\n")

	fis := NewFileInputStream(streamFile, "test", NewFileState(stateFile))

	line1, _ := fis.ReadLine()
	assertSame("Line 1", line1, t)

	go func() {
		time.Sleep(time.Second)
		file.WriteString("Line 2\n")
	}()

	line2, _ := fis.ReadLine()
	assertSame("Line 2", line2, t)

	fis.Close()
}

func Test_FileInputStream_ReadLine_Skip(t *testing.T) {
	truncateState()
	file, err := os.OpenFile(streamFile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		t.Fatal(err)
	}

	file.WriteString("Line 1\n")
	file.WriteString("Line 2\n")

	fileInputStream1 := NewFileInputStream(streamFile, "test", NewFileState(stateFile))
	line1, _ := fileInputStream1.ReadLine()
	assertSame("Line 1", line1, t)
	fileInputStream1.Close()

	fileInputStream2 := NewFileInputStream(streamFile, "test", NewFileState(stateFile))
	line2, _ := fileInputStream2.ReadLine()
	assertSame("Line 2", line2, t)
	fileInputStream2.Close()
}
