package processors

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

func Test_RegexProcessor_Process(t *testing.T) {
	var processor Processor

	processor = NewRegexProcessor("([a-zA-Z.]*) ([0-9]*)", "{ \"host\" : \"(1)\", \"port\" : \"(2)\" }")

	expected := "{ \"host\" : \"somerandomhost.com\", \"port\" : \"1234\" }"
	received := processor.Process("somerandomhost.com 1234 that isnt needed")
	if(strings.Compare(received, expected) != 0){
		t.Error("Expected: <<", expected, ">>, Recieved: <<", received, ">>")
	}
}
