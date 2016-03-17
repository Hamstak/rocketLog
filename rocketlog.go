package main

import (
	"flag"
	"os"
	"log"
	"os/signal"
	"sync"
	"runtime"
	"time"

	"github.com/hamstak/rocketlog/config"
	"github.com/hamstak/rocketlog/inputs"
	"github.com/hamstak/rocketlog/processors"
	"github.com/hamstak/rocketlog/outputs"
	"github.com/hamstak/rocketlog/event"
)

var verbose bool

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


func destroy_rocket_instances(rocket_inputs []inputs.Input, rocket_outputs []outputs.Output){
	for _, output := range rocket_outputs {
		output.Close()
	}

	for _, input := range rocket_inputs {
		input.Close()
	}
}

func handle_interrupt(cb func()){
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	go func() {
		for _ = range signalChan {
			log.Print("\nReceived an interrupt, stopping services...\n")
			cb()
		}
	}()
}

func consume_event(event *event.Event, rocket_processors []processors.Processor, rocket_outputs []outputs.Output){
	for _, processor := range rocket_processors {
		if(processor.Matches(event.Data)){
			event.Data = processor.Process(event.Data)
			output_event(event, rocket_outputs)
		}
	}
}

func output_event(event *event.Event, rocket_outputs []outputs.Output){
	for _, output := range rocket_outputs {
		output.Write(event)
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

	lock := sync.Mutex{}

	handle_interrupt(func(){
		log.Print("Destroying Rocket_Instances")
		destroy_rocket_instances(rocket_inputs, rocket_outputs)
		lock.Lock()
		os.Exit(0)
		lock.Unlock()
	})

	for {
		lock.Lock()

		vlog("working")
		log.Print("Sleeping")

		for i, input := range rocket_inputs {
			line, err := input.ReadLine()
			if(err != nil){
				vlog("Reached end of input ", i)
				input.Flush()
				time.Sleep(time.Second * 1)
			} else {
				e := event.NewEvent(line, "rocket_input", input.GetType())
				consume_event(e, rocket_processors, rocket_outputs)
			}
		}

		lock.Unlock()
		runtime.Gosched()
	}
}