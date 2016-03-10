package rocketLog

import "testing"

func TestXML(t *testing.T){
	e := eventFactory("<s></s>")
	if(e.dataType != XML){
		t.Error("Data Type Failure")
	}
}

func TestJSON(t *testing.T){
	e:= eventFactory("{}")
	if(e.dataType != JSON) {
		t.Error("Data Type Failure")
	}
}

func TestRAW(t *testing.T){
	e := eventFactory("Some text")
	if(e.dataType != RAW){
		t.Error("Data Type Failure")
	}
}
