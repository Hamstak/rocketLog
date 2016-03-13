package rocketLog

import (
	"strings"
	//"fmt"
	//"encoding/json"
	//"encoding/xml"
)

const RAW = "RAW"
const JSON = "JSON"
const XML = "XML"

type Event struct{
	Data     string
	Producer string
	DataType string
	Index    string
}

func NewEvent(payload, producer, index string) *Event{
	trimmed := strings.Trim(payload, " \t\n")
	dataType := RAW

	if(strings.IndexByte("{[", trimmed[0]) != -1 && strings.IndexByte("}]", trimmed[len(trimmed) - 1]) != -1){
		dataType = JSON
	}else if (trimmed[0] == '<' && trimmed[len(trimmed) - 1] == '>'){
		dataType = XML
	}

	return &Event{
		Data: payload,
		Producer: producer,
		DataType: dataType,
		Index: index,
	}

}
