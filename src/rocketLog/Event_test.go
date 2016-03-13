package rocketLog

import (
	"testing"
)

func TestXMLDetection(t *testing.T){
	e := NewEvent("<s></s>", "TEST", "TEST")
	if(e.DataType != XML){
		t.Error("Data Type Failure")
	}
}

func TestJSONDetection(t *testing.T){
	e:= NewEvent("{}", "TEST", "TEST")
	if(e.DataType != JSON) {
		t.Error("Data Type Failure")
	}
}

func TestRAWDetection(t *testing.T){
	e := NewEvent("Some text", "TEST", "TEST")
	if(e.DataType != RAW){
		t.Error("Data Type Failure")
	}
}

func TestConfigurationInputGeneral(t *testing.T){
	c, err := ReadConfiguration("testfiles/config.yml")
	if( err != nil){
		panic(err)
	}
	if (c.Input.Webservice[0].portAddress != "https://0.0.0.0:0000/"){
		t.Error("Baking failure")
	}
}

