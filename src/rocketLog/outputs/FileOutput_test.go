package outputs

import (
	"testing"
	"os"
	"log"
)


func TestMain(m *testing.M){
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	ret := m.Run()
	os.Exit(ret)
}

func Test_FileOuput_FileWrite(t *testing.T) {
	file_name := "./testfiles/output.txt"

	output := NewFileOutput(file_name)
	output.Write("Hello World!")
	output.Write("SecondLine")
	output.Close()

	os.Remove(file_name)
}
