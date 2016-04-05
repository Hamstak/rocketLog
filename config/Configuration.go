package config

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Configuration struct {
	Input      InputInstance
	Processing ProcessingInstance
	Output     OutputInstance
}

type InputInstance struct {
	File []FileInputInstance
}

type OutputInstance struct {
	File       []FileOutputInstance
	Webservice []WebserviceInstance
}

type FileInputInstance struct {
	File string
	Type string
}

type FileOutputInstance struct {
	File string
}

type WebserviceInstance struct {
	Url string
}

type ProcessingInstance struct {
	Regex []RegexInstance
}

type RegexInstance struct {
	Regex   string
	Mapping string
    Name    string
}

func NewConfiguration(fileName string) (*Configuration, error) {
	config := &Configuration{}
	var err error

	dat, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(dat, config)
	if err != nil {
		panic(err)
	}

	err = errorHandle(config)

	return config, err
}

func errorHandle(config *Configuration) error {
	var err error
	err = nil
	if (len(config.Input.File) == 0) && (len(config.Output.File) == 0 && len(config.Output.Webservice) == 0) {
		err = errors.New("No valid inputs or outputs detected in definition of configuration file")
	} else if len(config.Input.File) == 0 {
		err = errors.New("No valid inputs detected in definition of configuration file")
	} else if len(config.Output.File) == 0 && len(config.Output.Webservice) == 0 {
		err = errors.New("No valid outpits detected in definition of configuration file")
	}

	return err
}
