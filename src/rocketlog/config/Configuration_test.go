package config

import (
	"testing"
	"log"
	"os"
	"strings"
)

func TestMain(m *testing.M){
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	ret := m.Run()
	os.Exit(ret)
}

func Test_Configuration_Input(t *testing.T){
	c, err := NewConfiguration("testfiles/config.yml")
	if( err != nil){
		panic(err)
	}

	var expected, actual string

	// Test file input
	expected = "./input.txt"
	actual = c.Input.File[0].File
	compare(expected, actual, t)

	expected = "input-file"
	actual = c.Input.File[0].Type
	compare(expected, actual, t)
}


func Test_Configuration_Processing(t *testing.T){
	c, err := NewConfiguration("testfiles/config.yml")
	if( err != nil){
		panic(err)
	}

	var expected, actual string

	// Test first regex
	expected = "^(this)*"
	actual = c.Processing.Regex[0].Regex
	compare(expected, actual, t)

	expected = "(1) thing"
	actual = c.Processing.Regex[0].Mapping
	compare(expected, actual, t)

	// Test second regex
	expected = "^(foo)*"
	actual = c.Processing.Regex[1].Regex
	compare(expected, actual, t)

	expected = "(1)(1)"
	actual = c.Processing.Regex[1].Mapping
	compare(expected, actual, t)
}

func Test_Configuration_OutputGeneral(t *testing.T){
	c, err := NewConfiguration("testfiles/config.yml")
	if( err != nil){
		panic(err)
	}

	var expected, actual string

	// Test first file output
	expected = "./output.txt"
	actual = c.Output.File[0].File
	compare(expected, actual, t)

	// Test second file output
	expected = "./output2.txt"
	actual = c.Output.File[1].File
	compare(expected, actual, t)

	// Test first elastic output
	expected = "http://127.0.0.1:9200"
	actual = c.Output.Webservice[0].Url
	compare(expected, actual, t)

	// Test second elastic output
	expected = "http://127.0.0.1:9201"
	actual = c.Output.Webservice[1].Url
	compare(expected, actual, t)
}

func compare(expected, actual string, t *testing.T){
	if (strings.Compare(expected, actual) != 0){
		t.Error("Expected: <<", expected, ">> Actual <<", actual, ">>")
	}
}
