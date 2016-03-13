package rocketLog

import "testing"

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

func TestConfigurationInput(t *testing.T){
	e:= readConfiguration()
	if(e.webservice != "something"){
		t.Error("Some read error")
	}
}
