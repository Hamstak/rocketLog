package event

import (
	"strings"
)

const RAW = "RAW"
const JSON = "JSON"
const XML = "XML"

type Event struct {
	Data     string
	DataType string
	Index    string
}

func NewEvent(payload, index string) *Event {
	trimmed := strings.Trim(payload, " \t\n")
	dataType := RAW

	if strings.IndexByte("{[", trimmed[0]) != -1 && strings.IndexByte("}]", trimmed[len(trimmed)-1]) != -1 {
		dataType = JSON
	} else if trimmed[0] == '<' && trimmed[len(trimmed)-1] == '>' {
		dataType = XML
	}

	return &Event{
		Data:     payload,
		DataType: dataType,
		Index:    index,
	}

}
