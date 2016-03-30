package inputs

import (
	"bufio"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

// FileInputStream is a rocketlog input for reading from files.
type FileInputStream struct {
	reader           *bufio.Reader
	absoluteFilePath string
	relativeFilePath string
	etype            string
	file             *os.File
	lineNumber       int
	state            *FileState
}

// GetName returns the name for the FileInputStream
func (fileInputStream *FileInputStream) GetName() string {
	return "FileInputStream='" + fileInputStream.relativeFilePath + "'"
}

// NewFileInputStream is the constructor for the FileInputStream
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

// GetType returns the elasticsearch type that this input will produce events for.
func (fileInputStream *FileInputStream) GetType() string {
	return fileInputStream.etype
}

// ReadByte reads a single byte from a file.
func (fileInputStream *FileInputStream) ReadByte() (byte, error) {
	return fileInputStream.reader.ReadByte()
}

func (fileInputStream *FileInputStream) skip() {
	targetLineNumber := fileInputStream.state.Load(fileInputStream.absoluteFilePath)
	fileInputStream.skipToLine(targetLineNumber)
}

func (fileInputStream *FileInputStream) skipToLine(lineNumber int) {
	for ; fileInputStream.lineNumber < lineNumber; fileInputStream.lineNumber++ {
		for {
			byte, err := fileInputStream.ReadByte()
			if err != nil {
				log.Fatal(err)
			}

			if byte == '\n' {
				break
			}
		}
	}
}

// ReadLine reads a line from the file. If it doesn't have
func (fileInputStream *FileInputStream) ReadLine() (string, error) {
	maxDuration := time.Second * 30

	// Buffer related numbers
	bufferLen := 1024
	buffer := make([]byte, bufferLen)
	bufferIndex := 0

	// Backoff related numbers
	duration := time.Millisecond

	for bufferIndex = 0; bufferIndex <= bufferLen; bufferIndex++ {
		currentByte, err := fileInputStream.ReadByte()

		if err != nil && err != io.EOF {
			log.Fatal(err)
		} else if err == io.EOF { // If EOF sleep and decrement bufferIndex.
			bufferIndex--

			if duration*2 < maxDuration {
				duration *= 2
			}

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
				fileInputStream.lineNumber++
				fileInputStream.saveState()
				break
			}

			buffer[bufferIndex] = currentByte
		}
	}

	return string(buffer[0:bufferIndex]), nil
}

func (fileInputStream *FileInputStream) saveState() {
	fileInputStream.state.Save(fileInputStream.absoluteFilePath, fileInputStream.lineNumber)
}

// Close closes the file descriptor
func (fileInputStream *FileInputStream) Close() {
	fileInputStream.file.Close()
}
