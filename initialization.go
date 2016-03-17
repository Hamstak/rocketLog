package main

import (
	"github.com/hamstak/rocketlog/config"
	"github.com/hamstak/rocketlog/inputs"
	"github.com/hamstak/rocketlog/processors"
	"github.com/hamstak/rocketlog/outputs"
)


func populate_rocket_instances(
config_struct *config.Configuration,
rocket_inputs []inputs.Input,
rocket_processors []processors.Processor,
rocket_outputs []outputs.Output){

	for i, input_instance := range config_struct.Input.File {
		rocket_inputs[i] = inputs.NewFileInput(input_instance.File, "state.json", input_instance.Type)
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
}
