package inputs

import (
	"os"
	"bufio"
	"io"
	"time"
	"log"
	"path/filepath"
)

type FileInputStream struct {
	reader             *bufio.Reader
	absolute_file_path string
	file               *os.File
	line_number        int
	state              *FileState
}

func NewFileInputStream(path string, state *FileState) *FileInputStream{
	absolute_file_path, err := filepath.Abs(path)
	if(err != nil){
		log.Fatal(err)
	}

	file, err := os.Open(absolute_file_path)
	if(err != nil){
		log.Fatal(err)
	}

	reader := bufio.NewReader(file)

	file_input_stream := &FileInputStream{
		file: file,
		reader: reader,
		absolute_file_path: absolute_file_path,
		state:state,
	}

	file_input_stream.skip()

	return file_input_stream
}

func (self *FileInputStream) ReadByte() (byte, error) {
	return self.reader.ReadByte()
}

func (self *FileInputStream) skip(){
	target_line_number := self.state.Load(self.absolute_file_path)
	self.skip_to_line(target_line_number)
}

func (self *FileInputStream) skip_to_line(line_number int){
	for ; self.line_number < line_number; self.line_number++ {
		for {
			byte, err := self.ReadByte()
			if(err != nil){
				log.Fatal(err)
			}

			if(byte == '\n'){
				break;
			}
		}
	}
}

func (self *FileInputStream) ReadLine() string{
	// Buffer related numbers
	buffer_len := 1024
	buffer := make([]byte, buffer_len)
	buffer_index := 0

	// Backoff related numbers
	duration := time.Millisecond

	for buffer_index = 0; buffer_index < buffer_len; buffer_index++ {
		current_byte, err := self.ReadByte()

		if(err != nil && err != io.EOF){
			log.Fatal(err)
		} else if(err == io.EOF){ // If EOF sleep and decrement buffer_index.
			buffer_index--
			duration *= 2
			time.Sleep(duration)
		} else {
			if (duration > 1){
				duration /= 2
			}

			if(buffer_index == buffer_len){
				new_buffer := make([]byte, buffer_len * 2)
				copy(new_buffer[0:buffer_len], buffer[:])
				buffer = new_buffer
			}

			if(current_byte == '\n'){
				self.line_number++
				self.save_state()
				break;
			}

			buffer[buffer_index] = current_byte
		}
	}

	return string(buffer[0:buffer_index])
}

func (self *FileInputStream) save_state(){
	self.state.Save(self.absolute_file_path, self.line_number)
}

func (self *FileInputStream) Close(){
	self.file.Close()
}
