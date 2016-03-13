package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"errors"
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
	Type string
}

type WebserviceInstance struct  {
	Port string
	Address string
	Type string
	portAddress string
}

type ProcessingInstance struct {
	Regex []RegexInstance
}

type RegexInstance struct {
	Regex string
	Mapping string
}

func ReadConfiguration(fileName string) (*Configuration, error){
	config := &Configuration{}
	var err error


	dat, err := ioutil.ReadFile(fileName)
	if(err != nil){
		panic(err)
	}

	err = yaml.Unmarshal(dat, config)
	if(err != nil){
		panic(err)
	}

	for i := 0; i < len(config.Input.Webservice); i++{
		config.Input.Webservice[i].portAddress = "https://" + config.Input.Webservice[i].Address + ":" + config.Input.Webservice[i].Port + "/"
	}

	err = errorHandle(config)

	return config, err
}

func errorHandle(config *Configuration) error{
	var err error
	err = nil
	if ((len(config.Input.Webservice) == 0 && len(config.Input.File) == 0 )&& (len(config.Output.File) == 0 && len(config.Output.Webservice) == 0)){
		err = errors.New("No valid inputs or outputs detected in definition of configuration file")
	}else if(len(config.Input.Webservice) == 0 && len(config.Input.File) == 0){
		err = errors.New("No valid inputs detected in definition of configuration file")
	}else if(len(config.Output.File) == 0 && len(config.Output.Webservice) == 0){
		err = errors.New("No valid outpits detected in definition of configuration file")
	}

	return err
}