package outputs

import (
	"errors"
	"github.com/hamstak/rocketlog/event"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

func startElasticSearch(t *testing.T) error {
	for i := 0; i < 30; i++ {
		resp, err := http.Get("http://localhost:9200")
		if err == nil && resp.StatusCode == 200 {
			return nil
		}
		time.Sleep(time.Second)
		log.Print("Elasticsearch loading..")
	}

	return errors.New("Elastic timeout")
}

func Test_NetOuput_Write(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping NetOutput_Write (elsaticsearch) test")
	}

	err := startElasticSearch(t)
	if err != nil {
		log.Print(err)
		t.Fail()
		return
	}

	var output Output
	output = NewNetOutput("http://localhost:9200")
	output.Write(event.NewEvent("{ \"Foo\":\"Bar\" }", "test-index"))
	output.Close()

	os.RemoveAll("./testfiles/esdata")
}
