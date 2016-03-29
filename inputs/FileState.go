package inputs

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

// FileState is the shared resource that all FileInputStreams share to keep track
// of which line in the file the FileInputStream left off at.
type FileState struct {
	stateFilePath string
	mutex         *sync.Mutex
}

// NewFileState is the constructor for a FileState
func NewFileState(stateFilePath string) *FileState {
	mutex := &sync.Mutex{}
	return &FileState{
		stateFilePath: stateFilePath,
		mutex:         mutex,
	}

}

// Save takes the given line number and path, and saves it to its local map.
// It then takes the map and saves it to the state file.
func (fileState *FileState) Save(filename string, line_number int) {
	file_state_map := fileState.loadMap()
	file_state_map[filename] = line_number
	fileState.saveMap(file_state_map)
}

// Load reads the state file and updates the FileState's map to include the
// latest changes from the file. It then returns the line number for the
// specified file.
func (fileState *FileState) Load(filename string) int {
	fileStateMap := fileState.loadMap()
	return fileStateMap[filename]
}

func (self *FileState) loadMap() map[string]int {
	var file *os.File
	var err error

	self.mutex.Lock()
	defer self.mutex.Unlock()

	file, err = os.OpenFile(self.stateFilePath, os.O_CREATE|os.O_RDWR, 0666)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	decoder := json.NewDecoder(file)
	file_state := make(map[string]int)
	decoder.Decode(&file_state)
	return file_state
}

func (self *FileState) saveMap(fileStateMap map[string]int) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	file, err := os.OpenFile(self.stateFilePath, os.O_RDWR|os.O_TRUNC, 0666)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	e := json.NewEncoder(file)
	err = e.Encode(fileStateMap)
	if err != nil {
		log.Fatal(err)
	}

}
