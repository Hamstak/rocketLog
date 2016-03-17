package inputs

import (
	"testing"
)

func Test_FileState_SaveLoad(t *testing.T) {
	expected := 109
	file_name := "./some-random-nonexisting-file.txt"

	file_state := NewFileState("./testfiles/state.json")
	file_state.Save(file_name, expected)

	file_state_2 := NewFileState("./testfiles/state.json")
	actual := file_state_2.Load(file_name)

	if(expected != actual){
		t.Errorf("Expecting %d, Actuall %d", expected, actual)
	}
}