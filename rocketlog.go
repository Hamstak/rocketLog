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

	"net/http"

	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/expvar"
)

var verbose bool
var check bool

var producedEventsCounter metrics.Counter
var consumedEventsCounter metrics.Counter
var modifiedEventsCounter metrics.Counter
var unmatchedEventsCounter metrics.Counter

var counterMutex *sync.Mutex

func validateConfigPath(configPath string) {
	_, err := os.Stat(configPath)
	if err != nil {
		log.Fatal("config_file: ", configPath, " not found")
	}
}

func loadConfiguration() *config.Configuration {
	var configPath string
	flag.StringVar(&configPath, "config", "./configuration.yml", "The configuration file")
	flag.BoolVar(&verbose, "verbose", false, "print verbose output")
	flag.BoolVar(&check, "check", false, "just check if the configuration is valid")
	flag.Parse()

	validateConfigPath(configPath)

	configStruct, err := config.NewConfiguration(configPath)
	if err != nil {
		log.Fatal(err)
	}
	return configStruct
}

func handleCloseInterrupt(lock *sync.Mutex, rocketInstance *RocketInstance) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	go func() {
		for _ = range signalChan {
			log.Print("\nReceived an interrupt, stopping services...\n")
			log.Print("Destroying Rocket_Instances")
			lock.Lock()
			rocketInstance.Close()
			os.Exit(0)
			lock.Unlock()
		}
	}()
}

func modifyEvents(input, output chan *event.Event, rocketInstance *RocketInstance) {
	for {
		event := <-input

		err := modifyASingleEvent(event, rocketInstance)
		if err != nil {
			incremementCounter(unmatchedEventsCounter, 1)
			continue
		}

		incremementCounter(modifiedEventsCounter, 1)
		output <- event
	}
}

func modifyASingleEvent(event *event.Event, rocketInstance *RocketInstance) error {
	for _, processor := range rocketInstance.Processors {
		if processor.Matches(event.Data) {
			event.Data = processor.Process(event.Data)
			if verbose {
				logModifyEvent(event, processor)
			}
			return nil
		}
	}

	logUnmatchedEvent(event)
	return newNotMatchedError(event)
}

type notMatched struct {
	Event *event.Event
}

func (error *notMatched) Error() string {
	return "Couldn't find a match for " + error.Event.ToString()
}

func newNotMatchedError(event *event.Event) *notMatched {
	return &notMatched{
		Event: event,
	}
}

func consumeEvents(events chan *event.Event, rocketInstance *RocketInstance) {
	for {
		e := <-events
		incremementCounter(consumedEventsCounter, 1)
		for _, output := range rocketInstance.Outputs {
			output.Write(e)
			if verbose {
				logConsumeEvent(e, output)
			}
		}
	}
}

func produceEvents(output chan *event.Event, rocketInstance *RocketInstance) {
	for _, producer := range rocketInstance.Inputs {
		go produceEventForInput(output, producer)

		if verbose {
			logCreateProduceThread(producer)
		}
	}
}

func produceEventForInput(output chan *event.Event, input inputs.Input) {
	for {
		line, err := input.ReadLine()
		if err != nil {
			log.Fatal(err)
		}

		e := event.NewEvent(line, input.GetType())
		output <- e

		incremementCounter(producedEventsCounter, 1)

		if verbose {
			logEnqueueEvent(e, input)
		}
	}
}

func incremementCounter(counter metrics.Counter, number uint64) {
	counterMutex.Lock()
	counter.Add(number)
	counterMutex.Unlock()
}

func main() {
	configStruct := loadConfiguration()

	producedEventsCounter = expvar.NewCounter("produced_events")
	consumedEventsCounter = expvar.NewCounter("consumed_events")
	modifiedEventsCounter = expvar.NewCounter("modified_events")
	unmatchedEventsCounter = expvar.NewCounter("unmatched_events")

	counterMutex = &sync.Mutex{}

	if verbose {
		logConfiguration(configStruct)
	}

	if check {
		log.Print("Configuration Loaded OK")
		os.Exit(0)
	}

	rocketInstance := NewRocketInstance(configStruct)
	lock := &sync.Mutex{}

	handleCloseInterrupt(lock, rocketInstance)

	inputToModify := make(chan *event.Event, 1)
	modifyToOutput := make(chan *event.Event, 1)

	go produceEvents(inputToModify, rocketInstance)
	go modifyEvents(inputToModify, modifyToOutput, rocketInstance)
	go consumeEvents(modifyToOutput, rocketInstance)

	http.ListenAndServe(":1234", nil)
}
