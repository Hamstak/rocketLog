package event

import (
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	ret := m.Run()
	os.Exit(ret)
}

func Test_NewEvent_XMLDetection(t *testing.T) {
	e := NewEvent("<s></s>", "TEST")
	if e.DataType != Xml {
		t.Error("Data Type Failure")
	}
}

func Test_NewEvent_JSONDetection(t *testing.T) {
	e := NewEvent("{}", "TEST")
	if e.DataType != Json {
		t.Error("Data Type Failure")
	}
}

func Test_NewEvent_RAWDetection(t *testing.T) {
	e := NewEvent("Some text", "TEST")
	if e.DataType != Raw {
		t.Error("Data Type Failure")
	}
}
