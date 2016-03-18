package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/hamstak/rocketlog/config"
	"github.com/hamstak/rocketlog/event"
	"github.com/hamstak/rocketlog/inputs"
)

var verbose bool
var check bool

func validate_config_path(config_path string) {
	_, err := os.Stat(config_path)
	if err != nil {
		log.Fatal("config_file: ", config_path, " not found")
	}
}

func LoadConfiguration() *config.Configuration{
	var config_path string
	flag.StringVar(&config_path, "config", "./configuration.yml", "The configuration file")
	flag.BoolVar(&verbose, "verbose", false, "print verbose output")
	flag.BoolVar(&check, "check", false, "just check if the configuration is valid")
	flag.Parse()

	validate_config_path(config_path)

	config_struct, err := config.NewConfiguration(config_path)
	if err != nil {
		log.Fatal(err)
	}
	return config_struct
}

func PrintConfiguration(config_struct *config.Configuration){
	for i, file := range config_struct.Input.File {
		log.Println("File Input: ", i)
		log.Println("\tFile:", file.File)
		log.Println("\tType:", file.Type)
	}

	for i, regex := range config_struct.Processing.Regex {
		log.Println("Regex Processor: ", i)
		log.Println("\tRegex: ", regex.Regex)
		log.Println("\tMapping: ", regex.Mapping)
	}

	for i, file := range config_struct.Output.File {
		log.Println("File Output: ", i)
		log.Println("\tFile:", file.File)
	}

	for i, web := range config_struct.Output.Webservice {
		log.Println("Web Output: ", i)
		log.Println("\tUrl: ", web.Url)
	}
}



func handle_close_interrupt(lock *sync.Mutex, rocket_instance *RocketInstance) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	go func() {
		for _ = range signalChan {
			log.Print("\nReceived an interrupt, stopping services...\n")
			log.Print("Destroying Rocket_Instances")
			lock.Lock()
			rocket_instance.Close()
			os.Exit(0)
			lock.Unlock()
		}
	}()
}



func ModifyEvents(input, output chan *event.Event, rocket_instance *RocketInstance) {
	for {
		event := <-input
		if(modify_an_event(event, rocket_instance)){
			output <- event
		} else {
			log.Print("Unmatched Event: ", event)
		}

	}
}

func modify_an_event(event *event.Event, rocket_instance *RocketInstance) bool{
	for _, processor := range rocket_instance.rocket_processors {
		if processor.Matches(event.Data) {
			event.Data = processor.Process(event.Data)
			return true
		}
	}

	return false
}

func ConsumeEvents(events chan *event.Event, rocket_instance *RocketInstance) {
	for {
		e := <-events
		for _, output := range rocket_instance.rocket_outputs {
			output.Write(e)
		}
	}
}

func ProduceEvents(output chan *event.Event, rocket_instance *RocketInstance){
	for _, producer := range rocket_instance.rocket_inputs {
		go produce_event_for_input(output, producer)
	}
}

func produce_event_for_input(output chan *event.Event, input inputs.Input){
	for {
		line, err := input.ReadLine()
		if(err != nil){
			log.Fatal(err)
		}

		e := event.NewEvent(line, input.GetType())
		output <- e
	}
}

func main() {
	config_struct := LoadConfiguration()

	if(check){
		log.Print("Configuration Loaded OK")
		os.Exit(0)
	}

	if(verbose){
		PrintConfiguration(config_struct)
	}

	rocket_instance := NewRocketInstance(config_struct)
	lock := &sync.Mutex{}

	handle_close_interrupt(lock, rocket_instance)

	input_to_modify := make(chan *event.Event, 1)
	modify_to_output := make(chan *event.Event, 1)

	go ProduceEvents(input_to_modify, rocket_instance)
	go ModifyEvents(input_to_modify, modify_to_output, rocket_instance)
	ConsumeEvents(modify_to_output, rocket_instance)
}
