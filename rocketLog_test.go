package rocketLog

import "testing"

func TestXMLDetection(t *testing.T){
	e := eventFactory("<s></s>")
	if(e.dataType != XML){
		t.Error("Data Type Failure")
	}
}

func TestJSONDetection(t *testing.T){
	e:= eventFactory("{}")
	if(e.dataType != JSON) {
		t.Error("Data Type Failure")
	}
}

func TestRAWDetection(t *testing.T){
	e := eventFactory("Some text")
	if(e.dataType != RAW){
		t.Error("Data Type Failure")
	}
}

func TestConfigurationInput(t *testing.T){
	e:= readConfiguration()
	if(e.webservice != "something"){
		t.Error("Some read error")
	}
}
