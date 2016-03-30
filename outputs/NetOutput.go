package outputs

import (
	"encoding/json"
	"github.com/hamstak/rocketlog/event"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
)

// NetOutput writes objects to the elasticsearch api.
type NetOutput struct {
	hostname string
	client   http.Client
	lock     sync.Mutex
}

// ElasticIndex is the elastic search index in which all rocketlog events will be output to.
const ElasticIndex = "rocketlog"

// NewNetOutput is the constructor for a NetOutput Object. Creates a NetOutput that connects to an instance of elasticsearch
func NewNetOutput(hostname string) *NetOutput {
	newOutput := &NetOutput{
		hostname: hostname,
		client:   http.Client{},
	}

	return newOutput
}

func (netOutput *NetOutput) getEndpoint(event *event.Event) string {
	return netOutput.hostname + "/" + ElasticIndex + "/" + event.Index + "/"
}

func isValidJSON(payload string) bool {
	var payloadIntermediate map[string]interface{}

	json.Unmarshal([]byte(payload), &payloadIntermediate)
	_, err := json.Marshal(payloadIntermediate)
	if err != nil {
		return false
	}

	return true
}

// Write Writes the event to elasticsearch
func (netOutput *NetOutput) Write(event *event.Event) {
	if !isValidJSON(event.Data) {
		log.Print("Couldn't Post Data For ", event)
		return
	}

	netOutput.lock.Lock()
	defer netOutput.lock.Unlock()

	payload := strings.NewReader(event.Data)
	endpoint := netOutput.getEndpoint(event)
	method := http.MethodPost

	request, err := http.NewRequest(method, endpoint, payload)
	if err != nil {
		log.Fatal(err)
	}

	response, err := netOutput.client.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	if response.StatusCode != http.StatusCreated {
		log.Print("Response Status Code == ", response.StatusCode)
		log.Print("Response Headers: ", response.Header)
		body, _ := ioutil.ReadAll(response.Body)
		log.Print("Response Body: ", string(body))
		log.Fatal("Failed To Write To ", netOutput)
	}
}

// Closes the NetOutput object
func (netOutput *NetOutput) Close() {

}
