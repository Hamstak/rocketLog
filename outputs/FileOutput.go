package outputs

import (
	"bufio"
	"github.com/hamstak/rocketlog/event"
	"log"
	"os"
)

type FileOutput struct {
	file_name string
	file      *os.File
	writer    *bufio.Writer
}

func NewFileOutput(file_name string) *FileOutput {
	file, err := os.OpenFile(file_name, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		log.Fatal(err)
	}

	writer := bufio.NewWriter(file)

	file_output := &FileOutput{
		file_name: file_name,
		file:      file,
		writer:    writer,
	}

	return file_output
}

func (self *FileOutput) Write(event *event.Event) {
	line := event.Data

	_, err := self.writer.WriteString(line)
	if err != nil {
		log.Fatal(err)
	}

	self.writer.Flush()
}

func (self *FileOutput) Close() {
	self.file.Close()
}
