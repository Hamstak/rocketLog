package outputs

import (
	"testing"
	"os"
	"github.com/hamstak/rocketlog/event"
)

func Test_FileOuput_FileWrite(t *testing.T) {
	file_name := "./testfiles/output.txt"
	var output Output

	output = NewFileOutput(file_name)

	event := event.NewEvent("Hello World!", "Test-Index")
	output.Write(event)

	output.Close()

	os.Remove(file_name)
}
