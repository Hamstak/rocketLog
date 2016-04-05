package outputs

import (
	"bufio"
	"github.com/hamstak/rocketlog/event"
	"log"
	"os"
)

// FileOutput is an endpoint for rocketlog events that will write all the
// events to a file. It is used for debugging purposes and is formatted
// in a json-like format.
type FileOutput struct {
	fileName string
	file     *os.File
	writer   *bufio.Writer
}

// NewFileOutput is the constructor for the FileOutput object. It takes a
// parameter for the path at which it should output to.
func NewFileOutput(fileName string) *FileOutput {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		log.Fatal(err)
	}

	writer := bufio.NewWriter(file)

	file_output := &FileOutput{
		fileName: fileName,
		file:     file,
		writer:   writer,
	}

	return file_output
}

// Write writes the event to a file (path specified in the constructor)
func (self *FileOutput) Write(event *event.Event) {
	line := event.Data

	_, err := self.writer.WriteString(line)
	if err != nil {
		log.Fatal(err)
	}

	self.writer.Flush()
}

// Close closes the file descriptor
func (self *FileOutput) Close() {
	self.file.Close()
}

func (fileOutput *FileOutput) ToString() string {
    return "OUTPUT(Type: \"FileOutput\", File: \"" + fileOutput.fileName + "\")"
}