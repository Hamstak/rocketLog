package main

import (
	"flag"
	"rocketlog/config"
	"os"
	"log"
	"rocketlog/inputs"
	"rocketlog/processors"
	"rocketlog/outputs"
)

var verbose bool
var config_yml config.Configuration

func vlog(args ...interface{}){
	if(verbose){
		log.Print(args)
	}
}

func fatal(args ...interface{}){
	log.Fatal(args)
}

func validate_config_path(config_path string){
	_, err := os.Stat(config_path)
	if(err != nil){
		log.Fatal("config_file: ", config_path, " not found")
	}
}

func vlog_config(config_struct *config.Configuration){
	for i, file := range config_struct.Input.File {
		vlog("File Input: ", i)
		vlog("\tFile:", file.File)
		vlog("\tType:", file.Type)
	}

	for i, regex := range config_struct.Processing.Regex {
		vlog("Regex Processor: ", i)
		vlog("\tRegex: ", regex.Regex)
		vlog("\tMapping: ", regex.Mapping)
	}

	for i, file := range config_struct.Output.File {
		vlog("File Output: ", i)
		vlog("\tFile:", file.File)
	}

	for i, web := range config_struct.Output.Webservice {
		vlog("Web Output: ", i)
		vlog("\tUrl: ", web.Url)
	}
}

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

func destroy_rocket_instances(rocket_inputs []inputs.Input, rocket_outputs []outputs.Output){
	for _, output := range rocket_outputs {
		output.Close()
	}

	for _, input := range rocket_inputs {
		input.Close()
	}
}

func main(){
	var config_path string
	flag.StringVar(&config_path, "config", "./configuration.yml", "The configuration file")
	flag.BoolVar(&verbose, "verbose", false, "print verbose output")
	flag.Parse()

	validate_config_path(config_path)

	vlog("Loading Configuration")
	config_struct, err := config.NewConfiguration(config_path)
	if(err != nil){
		fatal(err)
	}
	vlog_config(config_struct)
	vlog("Configuration loaded")

	vlog("Creating Rocket_Instances")
	rocket_inputs := make([]inputs.Input, len(config_struct.Input.File))
	rocket_processors := make([]processors.Processor, len(config_struct.Processing.Regex))
	rocket_outputs := make([]outputs.Output, len(config_struct.Output.File) + len(config_struct.Output.Webservice))
	populate_rocket_instances(config_struct, rocket_inputs, rocket_processors, rocket_outputs)


	vlog("Destroying Rocket_Instances")
	destroy_rocket_instances(rocket_inputs, rocket_outputs)
}