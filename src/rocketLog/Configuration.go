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
	File []FileInstance
	Webservice []WebserviceInstance
}

type FileInstance struct {
	File string
}

type WebserviceInstance struct  {
	Port string
	Address string
}

type ProcessingInstance struct {
	Regex []RegexInstance
}

type RegexInstance struct {
	Regex string
	Mapping string
}

func readConfiguration() interface{}{
	config := &Configuration{}

	dat, err := ioutil.ReadFile("testfiles/config.yml")
	if(err != nil){
		panic(err)
	}

	err = yaml.Unmarshal(dat, config)
	if(err != nil){
		panic(err)
	}
	log.Print(config)
	return config
}