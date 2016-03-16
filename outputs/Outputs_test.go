package outputs

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
