package inputs

import (
	"os"
	"encoding/json"
	"log"
	"sync"
)

type FileState struct {
	state_file_path string
	mutex *sync.Mutex
}

func NewFileState(state_file_path string) *FileState {
	mutex := &sync.Mutex{}
	return &FileState{
		state_file_path: state_file_path,
		mutex: mutex,
	}

}

func (self *FileState) Save(filename string, line_number int){
	file_state_map := self.loadMap()
	file_state_map[filename] = line_number
	self.saveMap(file_state_map)
}

func (self *FileState) Load(filename string) int {
	file_state_map := self.loadMap()
	return file_state_map[filename]
}

func (self *FileState) loadMap() map[string] int {
	var file *os.File
	var err error

	self.mutex.Lock()
	defer self.mutex.Unlock()

	file, err = os.OpenFile(self.state_file_path, os.O_CREATE | os.O_RDWR, 0666)
	defer file.Close()
	if(err != nil){
		log.Fatal(err)
	}


	decoder := json.NewDecoder(file)
	file_state := make(map[string]int)
	decoder.Decode(&file_state)
	return file_state
}

func (self *FileState) saveMap(file_state_map map[string] int){
	self.mutex.Lock()
	defer self.mutex.Unlock()

	file, err := os.OpenFile(self.state_file_path, os.O_RDWR | os.O_TRUNC, 0666)
	defer file.Close()
	if(err != nil){
		log.Fatal(err)
	}

	e := json.NewEncoder(file)
	err = e.Encode(file_state_map)
	if(err != nil){
		log.Fatal(err)
	}


}
