package outputs

import (
	"net/http"
	"log"
	"encoding/json"
	"io/ioutil"
	"strings"
	"rocketlog/events"
	"sync"
)

type NetOutput struct{
	hostname string
	client http.Client
	lock sync.Mutex
}

const ELASTIC_INDEX = "rocketlog"

func NewNetOutput(hostname string) *NetOutput{
	net_output := &NetOutput{
		hostname: hostname,
		client: http.Client{},
	}

	return net_output
}

func (self *NetOutput) getEndpoint(event *event.Event) string {
	return self.hostname + "/" + ELASTIC_INDEX + "/" + event.Index + "/"
}

func isValidJSON(payload string) bool{
	var payload_intermediate map[string]interface{}

	json.Unmarshal([]byte(payload), &payload_intermediate)
	_, err := json.Marshal(payload_intermediate)
	if(err != nil){
		return false
	}

	return true
}

func (self *NetOutput) Write(event *event.Event) {
	if(!isValidJSON(event.Data)){
		log.Print("Couldn't Post Data For ", event)
		return
	}

	self.lock.Lock()
	defer self.lock.Unlock()

	payload := strings.NewReader(event.Data)
	endpoint := self.getEndpoint(event)
	method := http.MethodPost

	request, err := http.NewRequest(method, endpoint, payload)
	if (err != nil) {
		log.Fatal(err)
	}

	response, err := self.client.Do(request)
	if (err != nil) {
		log.Fatal(err)
	}

	if (response.StatusCode != http.StatusCreated){
		log.Print("Response Status Code == ", response.StatusCode)
		log.Print("Response Headers: ", response.Header)
		body, _ := ioutil.ReadAll(response.Body)
		log.Print("Response Body: ", string(body))
		log.Fatal("Failed To Write To ", self)
	}
}

func (self *NetOutput) Close(){

}