package inputs

import (
	"bufio"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

type FileInputStream struct {
	reader           *bufio.Reader
	absoluteFilePath string
	etype            string
	file             *os.File
	lineNumber       int
	state            *FileState
}

func (fileInputStream *FileInputStream) GetInputName() string {
	return "FileInputStream='" + fileInputStream.absoluteFilePath + "'"
}

func NewFileInputStream(path, etype string, state *FileState) *FileInputStream {
	absolute_file_path, err := filepath.Abs(path)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open(absolute_file_path)
	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(file)

	file_input_stream := &FileInputStream{
		file:             file,
		reader:           reader,
		absoluteFilePath: absolute_file_path,
		state:            state,
		etype:            etype,
	}

	file_input_stream.skip()

	return file_input_stream
}

func (self *FileInputStream) GetType() string {
	return self.etype
}

func (self *FileInputStream) ReadByte() (byte, error) {
	return self.reader.ReadByte()
}

func (self *FileInputStream) skip() {
	target_line_number := self.state.Load(self.absoluteFilePath)
	self.skip_to_line(target_line_number)
}

func (self *FileInputStream) skip_to_line(line_number int) {
	for ; self.lineNumber < line_number; self.lineNumber++ {
		for {
			byte, err := self.ReadByte()
			if err != nil {
				log.Fatal(err)
			}

			if byte == '\n' {
				break
			}
		}
	}
}

func (self *FileInputStream) ReadLine() (string, error) {
	// Buffer related numbers
	buffer_len := 1024
	buffer := make([]byte, buffer_len)
	buffer_index := 0

	// Backoff related numbers
	duration := time.Millisecond

	for buffer_index = 0; buffer_index < buffer_len; buffer_index++ {
		current_byte, err := self.ReadByte()

		if err != nil && err != io.EOF {
			log.Fatal(err)
		} else if err == io.EOF { // If EOF sleep and decrement buffer_index.
			buffer_index--
			duration *= 2
			time.Sleep(duration)
		} else {
			if duration > 1 {
				duration /= 2
			}

			if buffer_index == buffer_len {
				newBuffer := make([]byte, buffer_len*2)
				copy(newBuffer[0:buffer_len], buffer[:])
				buffer = newBuffer
			}

			if current_byte == '\n' {
				self.lineNumber++
				self.save_state()
				break
			}

			buffer[buffer_index] = current_byte
		}
	}

	return string(buffer[0:buffer_index]), nil
}

func (self *FileInputStream) save_state() {
	self.state.Save(self.absoluteFilePath, self.lineNumber)
}

func (self *FileInputStream) Close() {
	self.file.Close()
}
