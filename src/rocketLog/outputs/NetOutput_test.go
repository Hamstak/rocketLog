package outputs

import (
	"testing"
	"os"
)

func Test_NetOuput_Write(t *testing.T) {
	file_name := "./testfiles/output.txt"

	var output Output
	output = NewFileOutput(file_name)
	output.Write("Hello World!")
	output.Write("SecondLine")
	output.Close()

	os.Remove(file_name)
}
