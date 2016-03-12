package outputs

import (
	"os"
	"log"
	"bufio"
)

type FileOutput struct {
	file_name string
	file *os.File
	writer *bufio.Writer
}

func NewFileOutput(file_name string) *FileOutput {
	file, err := os.OpenFile(file_name, os.O_RDWR | os.O_CREATE, 0666)

	if(err != nil){
		log.Fatal(err)
	}

	file_output := &FileOutput{
		file_name: file_name,
		file: file,
		writer: bufio.NewWriter(file),
	}

	return file_output
}

func (self *FileOutput) Write(line string){
	self.writer.WriteString(line)
	self.writer.Flush()
}

func (self *FileOutput) Close(){
	self.file.Close()
}