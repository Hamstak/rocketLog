package processors

import (
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	ret := m.Run()
	os.Exit(ret)
}

func Test_RegexProcessor_Process(t *testing.T) {
	processor := NewRegexProcessor("([a-zA-Z.]*) ([0-9]*)", "{ \"host\" : \"(1)\", \"port\" : \"(2)\", \"whoami\" : \"`whoami`\" }")
	output, err := exec.Command("whoami").Output()
	if err != nil {
		log.Fatal(err)
	}

	who_am_i := strings.Trim(string(output), "\n")
	expected := "{ \"host\" : \"somerandomhost.com\", \"port\" : \"1234\", \"whoami\" : \"" + who_am_i + "\" }"
	received := processor.Process("somerandomhost.com 1234 that isnt needed")

	t.Log(received)

	if strings.Compare(received, expected) != 0 {
		t.Error("\nExpected: <<", expected, ">>\nRecieved: <<", received, ">>")
		t.Fail()
	}

}
