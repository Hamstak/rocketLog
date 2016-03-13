package processors

import (
	"testing"
	"log"
	"os"
	"strings"
	"os/exec"
)

const WHO_AM_I = " 11:39:26 up 3 days,  2:44,  3 users,  load average: 0.23, 0.13, 0.14\nUSER     TTY      FROM             LOGIN@   IDLE   JCPU   PCPU WHAT\nmoth     :0       :0               Thu07   ?xdm?   2:09m  0.52s init --user\nmoth     pts/4    :0.0             Fri10   45:11m  0.43s  0.43s bash\nmoth     pts/24   :0.0             Sat11    5.00s  0.40s  0.00s /tmp/go-build698847472/rocketlog/processors/_test/processors.test -test.v=true\n"

func TestMain(m *testing.M){
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	ret := m.Run()
	os.Exit(ret)
}

func Test_RegexProcessor_Process(t *testing.T) {
	var processor Processor

	processor = NewRegexProcessor("([a-zA-Z.]*) ([0-9]*)", "{ \"host\" : \"(1)\", \"port\" : \"(2)\", \"whoami\" : \"`whoami`\" }")
	output, err := exec.Command("whoami").Output()
	if(err != nil){
		log.Fatal(err)
	}
	whoami := strings.Trim(string(output), "\n")
	expected := "{ \"host\" : \"somerandomhost.com\", \"port\" : \"1234\", \"whoami\" : \"" + whoami +"\" }"
	received := processor.Process("somerandomhost.com 1234 that isnt needed")
	if(strings.Compare(received, expected) != 0){
		t.Error("Expected: <<", expected, ">>, Recieved: <<", received, ">>")
	}
}
