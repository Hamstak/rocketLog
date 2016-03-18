package main

import (
	"github.com/hamstak/rocketlog/inputs"
	"github.com/hamstak/rocketlog/processors"
	"github.com/hamstak/rocketlog/outputs"
	"github.com/hamstak/rocketlog/config"
)

type RocketInstance struct {
	rocket_inputs []inputs.Input
	rocket_processors []processors.Processor
	rocket_outputs []outputs.Output
	config_struct *config.Configuration
}

func NewRocketInstance(config_struct *config.Configuration) *RocketInstance {
	// Create Rocket IOP Instances
	rocket_inputs := make([]inputs.Input, len(config_struct.Input.File))
	rocket_processors := make([]processors.Processor, len(config_struct.Processing.Regex))
	rocket_outputs := make([]outputs.Output, len(config_struct.Output.File)+len(config_struct.Output.Webservice))

	file_state := inputs.NewFileState("state.json")

	// Populate Rocket IOP Instances
	for i, input_instance := range config_struct.Input.File {
		rocket_inputs[i] = inputs.NewFileInputStream(input_instance.File, input_instance.Type, file_state)
	}

	for i, regex_instance := range config_struct.Processing.Regex {
		rocket_processors[i] = processors.NewRegexProcessor(regex_instance.Regex, regex_instance.Mapping)
	}

	for i, file_out_instance := range config_struct.Output.File {
		rocket_outputs[i] = outputs.NewFileOutput(file_out_instance.File)
	}

	offset := len(config_struct.Output.File)
	for i, web_out_instance := range config_struct.Output.Webservice {
		rocket_outputs[i + offset] = outputs.NewNetOutput(web_out_instance.Url)
	}

	return &RocketInstance{
		rocket_inputs:rocket_inputs,
		rocket_processors:rocket_processors,
		rocket_outputs:rocket_outputs,
		config_struct: config_struct,
	}
}

func (self *RocketInstance) Close(){
	for _, output := range self.rocket_outputs {
		output.Close()
	}

	for _, input := range self.rocket_inputs {
		input.Close()
	}
}

