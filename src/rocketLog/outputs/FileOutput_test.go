package outputs

import (
	"testing"
	"os"
	"rocketLog"
)

func Test_FileOuput_FileWrite(t *testing.T) {
	file_name := "./testfiles/output.txt"
	var output Output

	output = NewFileOutput(file_name)

	event := rocketLog.NewEvent("Hello World!","Test-Producer", "Test-Index")
	output.Write(event)

	output.Close()

	os.Remove(file_name)
}
