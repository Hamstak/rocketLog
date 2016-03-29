package outputs

import (
	"errors"
	"github.com/hamstak/rocketlog/event"
	"log"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"
)

func startElasticSearch(t *testing.T) error {
	cmd := exec.Command("docker-compose", "-f", "./testfiles/docker-compose.yml", "up", "-d")
	err := cmd.Run()
	if err != nil {
		log.Print(err)
		t.Error(err)
	}

	for i := 0; i < 30; i++ {
		resp, err := http.Get("http://elasticsearch:9200")
		if err == nil && resp.StatusCode == 200 {
			return nil
		}
		time.Sleep(time.Second)
		log.Print("Elasticsearch loading..")
	}

	return errors.New("Elastic timeout")
}

func stopElasticSearch(t *testing.T) {
	cmd := exec.Command("docker-compose", "-f", "./testfiles/docker-compose.yml", "down")
	err := cmd.Run()
	if err != nil {
		t.Error(err)
	}
}

func Test_NetOuput_Write(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping NetOutput_Write (elsaticsearch) test")
	}

	err := startElasticSearch(t)
	if err != nil {
		log.Print(err)
		stopElasticSearch(t)
		t.Fail()
		return
	}

	var output Output
	output = NewNetOutput("http://elasticsearch:9200")
	output.Write(event.NewEvent("{ \"Foo\":\"Bar\" }", "test-index"))
	output.Close()

	stopElasticSearch(t)
	os.RemoveAll("./testfiles/esdata")
}
