package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Configuration struct{
	Input ioInstance
	Processing ProcessingInstance
	Output ioInstance
}

type ioInstance struct {
	Index string
	File []FileInstance
	Webservice []WebserviceInstance
}

type FileInstance struct {
	File string
	Index string 	`omitempty`
}

type WebserviceInstance struct  {
	Port string
	Address string
	Index string `omitempty`
	portAddress string
}

type ProcessingInstance struct {
	Regex []RegexInstance
}

type RegexInstance struct {
	Regex string
	Mapping string
}

func ReadConfiguration() interface{}{
	config := &Configuration{}

	dat, err := ioutil.ReadFile("testfiles/config.yml")
	if(err != nil){
		panic(err)
	}

	err = yaml.Unmarshal(dat, config)
	if(err != nil){
		panic(err)
	}

	for i := 0; i < len(config.Input.Webservice); i++{
		config.Input.Webservice[i] = "https://" + config.Input.Webservice[i].Address + ":" + config.Input.Webservice[i].Port + "/"
	}
	
	return config
}