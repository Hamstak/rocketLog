package main

import (
	"testing"
)

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

func TestConfigurationInputGeneral(t *testing.T){
	c, err := ReadConfiguration("testfiles/config.yml")
	if( err != nil){
		panic(err)
	}
	if (c.Input.Webservice[0].portAddress != "https://0.0.0.0:0000/"){
		t.Error("Baking failure")
	}
}

