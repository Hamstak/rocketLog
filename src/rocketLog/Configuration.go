package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
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

func readConfiguration() Configuration{
	var config Configuration

	//config.Input.File = make([]FileInstance, 1)
	//config.Input.Webservice = make([]WebserviceInstance, 1)
	//config.Input.File[0].File = "stdin"
	//config.Input.Webservice[0].Address = "0.0.0.0"
	//config.Input.Webservice[0].Port = "0000"
	//
	//config.Processing.Regex = make([]RegexInstance, 1)
	//config.Processing.Regex[0].Regex = "^(this)*"
	//config.Processing.Regex[0].Mapping = "(1) thing"
	//
	//config.Output.File = make([]FileInstance, 1)
	//config.Output.Webservice = make([]WebserviceInstance, 1)
	//config.Output.File[0].File = "stdout"
	//config.Output.Webservice[0].Port = "0000"
	//config.Output.Webservice[0].Address = "0.0.0.0"


	//dat, err := yaml.Marshal(config)
	//if (err != nil){
	//	panic(err)
	//}
	//log.Print(dat)
	//err = ioutil.WriteFile("config_test.yml",dat, 0644 )
	//if(err != nil){
	//	panic(err)
	//}
	dat, err := ioutil.ReadFile("config_test.yml")
	if(err != nil){
		panic(err)
	}
	err = yaml.Unmarshal(dat, config)
	return config
}