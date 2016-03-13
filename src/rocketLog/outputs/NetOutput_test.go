package outputs

import (
	"testing"
	"os"
	"os/exec"
	"rocketLog"
	"time"
	"net/http"
	"log"
	"errors"
)


func startElasticSearch(t *testing.T) error {
	cmd := exec.Command("docker-compose", "-f", "./testfiles/docker-compose.yml", "up", "-d")
	err := cmd.Run()
	if(err != nil){
		log.Print(err)
		t.Error(err)
	}

	for i := 0; i < 30; i++ {
		resp, err := http.Get("http://elasticsearch:9200")
		if(err == nil && resp.StatusCode == 200){
			return nil
		}
		time.Sleep(time.Second)
		log.Print("Elasticsearch loading..")
	}

	return errors.New("Elastic timeout")
}

func stopElasticSearch(t *testing.T){
	cmd := exec.Command("docker-compose", "-f", "./testfiles/docker-compose.yml", "down")
	err := cmd.Run()
	if(err != nil){
		t.Error(err)
	}
}



func Test_NetOuput_Write(t *testing.T) {
	err := startElasticSearch(t)
	if(err != nil){
		log.Print(err)
		stopElasticSearch(t)
		t.Fail()
		return
	}

	var output Output
	output = NewNetOutput("http://elasticsearch:9200")
	output.Write(rocketLog.NewEvent("{ \"Foo\":\"Bar\" }", "TEST", "test-index"))
	output.Close()

	stopElasticSearch(t)
	os.RemoveAll("./testfiles/esdata")
}
