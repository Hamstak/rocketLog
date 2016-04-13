package config

import (
	"errors"
	"github.com/hamstak/rocketlog/inputs"
	"github.com/hamstak/rocketlog/outputs"
	"github.com/hamstak/rocketlog/processors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// Configuration contains the information for constructing rocketlog instances.
type Configuration struct {
	Input      inputConfiguration
	Processing processingConfiguration
	Output     outputInstance
}

type inputConfiguration struct {
	File   []fileInputConfiguration
	Syslog []syslogServerConfiguration
}

type outputInstance struct {
	File          []fileOutputConfiguration
	Elasticsearch []elasticsearchConfiguration
}

type syslogServerConfiguration struct {
	Port     int
	Protocol string
	Type     string
}

type fileInputConfiguration struct {
	File string
	Type string
}

type fileOutputConfiguration struct {
	File string
}

type elasticsearchConfiguration struct {
	Url string
}

type processingConfiguration struct {
	Regex []regexConfiguration
}

type regexConfiguration struct {
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

type RocketInstance struct {
	Inputs     []inputs.Input
	Processors []processors.Processor
	Outputs    []outputs.Output
	Config     *Configuration
}

// CreateInstance a rocketlog instance from the configuration.
func (configuration *Configuration) CreateInstance() *RocketInstance {
	// Create Rocket IOP Instances
	rocketInputs := make([]inputs.Input, len(configuration.Input.File))
	rocketProcessors := make([]processors.Processor, len(configuration.Processing.Regex))
	rocketOutputs := make([]outputs.Output, len(configuration.Output.File)+len(configuration.Output.Elasticsearch))

	file_state := inputs.NewFileState("state.json")

	// Populate Rocket IOP Instances
	for i, inputInstance := range configuration.Input.File {
		rocketInputs[i] = inputs.NewFileInputStream(inputInstance.File, inputInstance.Type, file_state)
	}

	for i, regexInstance := range configuration.Processing.Regex {
		rocketProcessors[i] = processors.NewRegexProcessor(regexInstance.Name, regexInstance.Regex, regexInstance.Mapping)
	}

	for i, fileOutputInstance := range configuration.Output.File {
		rocketOutputs[i] = outputs.NewFileOutput(fileOutputInstance.File)
	}

	offset := len(configuration.Output.File)
	for i, webOutputInstance := range configuration.Output.Elasticsearch {
		rocketOutputs[i+offset] = outputs.NewNetOutput(webOutputInstance.Url)
	}

	return &RocketInstance{
		Inputs:     rocketInputs,
		Processors: rocketProcessors,
		Outputs:    rocketOutputs,
		Config:     configuration,
	}
}

// Close shuts down and cleans up the rocket instance.
func (rocketInstance *RocketInstance) Close() {
	for _, output := range rocketInstance.Outputs {
		output.Close()
	}

	for _, input := range rocketInstance.Inputs {
		input.Close()
	}
}

func errorHandle(config *Configuration) error {
	var err error
	err = nil
	if (len(config.Input.File) == 0) && (len(config.Output.File) == 0 && len(config.Output.Elasticsearch) == 0) {
		err = errors.New("No valid inputs or outputs detected in definition of configuration file")
	} else if len(config.Input.File) == 0 {
		err = errors.New("No valid inputs detected in definition of configuration file")
	} else if len(config.Output.File) == 0 && len(config.Output.Elasticsearch) == 0 {
		err = errors.New("No valid outpits detected in definition of configuration file")
	}

	return err
}
