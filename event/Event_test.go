package event

import (
	"testing"
	"log"
	"os"
)

func TestMain(m *testing.M){
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	ret := m.Run()
	os.Exit(ret)
}

func Test_NewEvent_XMLDetection(t *testing.T){
	e := NewEvent("<s></s>", "TEST")
	if(e.DataType != XML){
		t.Error("Data Type Failure")
	}
}

func Test_NewEvent_JSONDetection(t *testing.T){
	e:= NewEvent("{}", "TEST")
	if(e.DataType != JSON) {
		t.Error("Data Type Failure")
	}
}

func Test_NewEvent_RAWDetection(t *testing.T){
	e := NewEvent("Some text", "TEST")
	if(e.DataType != RAW){
		t.Error("Data Type Failure")
	}
}

