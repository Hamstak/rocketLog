package main

import (
	"github.com/hamstak/rocketlog/config"
	"github.com/hamstak/rocketlog/inputs"
	"github.com/hamstak/rocketlog/outputs"
	"github.com/hamstak/rocketlog/processors"
)

type RocketInstance struct {
	Inputs     []inputs.Input
	Processors []processors.Processor
	Outputs    []outputs.Output
	Config     *config.Configuration
}

func NewRocketInstance(configuration *config.Configuration) *RocketInstance {
	// Create Rocket IOP Instances
	rocketInputs := make([]inputs.Input, len(configuration.Input.File))
	rocketProcessors := make([]processors.Processor, len(configuration.Processing.Regex))
	rocketOutputs := make([]outputs.Output, len(configuration.Output.File)+len(configuration.Output.Webservice))

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
	for i, webOutputInstance := range configuration.Output.Webservice {
		rocketOutputs[i+offset] = outputs.NewNetOutput(webOutputInstance.Url)
	}

	return &RocketInstance{
		Inputs:     rocketInputs,
		Processors: rocketProcessors,
		Outputs:    rocketOutputs,
		Config:     configuration,
	}
}

func (self *RocketInstance) Close() {
	for _, output := range self.Outputs {
		output.Close()
	}

	for _, input := range self.Inputs {
		input.Close()
	}
}
