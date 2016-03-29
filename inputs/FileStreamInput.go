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
	relativeFilePath string
	etype            string
	file             *os.File
	lineNumber       int
	state            *FileState
}

func (fileInputStream *FileInputStream) GetName() string {
	return "FileInputStream='" + fileInputStream.relativeFilePath + "'"
}

func NewFileInputStream(path, etype string, state *FileState) *FileInputStream {
	absoluteFilePath, err := filepath.Abs(path)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open(absoluteFilePath)
	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(file)

	fileInputStream := &FileInputStream{
		file:             file,
		reader:           reader,
		absoluteFilePath: absoluteFilePath,
		relativeFilePath: path,
		state:            state,
		etype:            etype,
	}

	fileInputStream.skip()

	return fileInputStream
}

func (fileInputStream *FileInputStream) GetType() string {
	return fileInputStream.etype
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
	bufferLen := 1024
	buffer := make([]byte, bufferLen)
	bufferIndex := 0

	// Backoff related numbers
	duration := time.Millisecond

	for bufferIndex = 0; bufferIndex <= bufferLen; bufferIndex++ {
		currentByte, err := self.ReadByte()

		if err != nil && err != io.EOF {
			log.Fatal(err)
		} else if err == io.EOF { // If EOF sleep and decrement bufferIndex.
			bufferIndex--
			duration *= 2
			time.Sleep(duration)
		} else {
			if duration > 1 {
				duration /= 2
			}

			if bufferIndex == bufferLen {
				newBuffer := make([]byte, bufferLen*2)
				copy(newBuffer[0:bufferLen], buffer[:])
				buffer = newBuffer
				bufferLen *= 2
			}

			if currentByte == '\n' {
				self.lineNumber++
				self.saveState()
				break
			}

			buffer[bufferIndex] = currentByte
		}
	}

	return string(buffer[0:bufferIndex]), nil
}

func (self *FileInputStream) saveState() {
	self.state.Save(self.absoluteFilePath, self.lineNumber)
}

func (self *FileInputStream) Close() {
	self.file.Close()
}
