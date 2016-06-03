package main

import (
	"github.com/hamstak/rocketlog/config"
	"github.com/hamstak/rocketlog/event"
	"github.com/hamstak/rocketlog/inputs"
	"github.com/hamstak/rocketlog/outputs"
	"github.com/hamstak/rocketlog/processors"

	"log"
)

const modifyEventString = "Modifing "
const consumeEventString = "Consuming"
const enqueueEventString = "Enqueing "

func logUnmodifiedEvent(e *event.Event) {

}

func logModifyEvent(e *event.Event, processor processors.Processor) {
	log.Println(modifyEventString, e.ToString(), "By", processor.ToString())
}

func logConsumeEvent(e *event.Event, output outputs.Output) {
	log.Println(consumeEventString, e.ToString(), "To", output.ToString())
}

func logCreateProduceThread(producer inputs.Input) {
	log.Println("Creating Producer Thread For -", producer.ToString())
}

func logEnqueueEvent(event *event.Event, input inputs.Input) {
	log.Println(enqueueEventString, event.ToString(), "From", input.ToString())
}

func logConfiguration(configStruct *config.Configuration) {
	for i, file := range configStruct.Input.File {
		log.Println("File Input: ", i)
		log.Println("\tFile:", file.File)
		log.Println("\tType:", file.Type)
	}

	for i, regex := range configStruct.Processing.Regex {
		log.Println("Regex Processor: ", i)
		log.Println("\tRegex: ", regex.Regex)
		log.Println("\tMapping: ", regex.Mapping)
		log.Println("\tName: ", regex.Name)
	}

	for i, file := range configStruct.Output.File {
		log.Println("File Output: ", i)
		log.Println("\tFile:", file.File)
	}

	for i, web := range configStruct.Output.Webservice {
		log.Println("Web Output: ", i)
		log.Println("\tUrl: ", web.Url)
	}
}
