package event

import (
	"strings"
)

const Raw = "RAW"
const Json = "JSON"
const Xml = "XML"

type Event struct {
	Data     string
	DataType string
	Index    string
}

// NewEvent is the constructor for an Event. It takes a payload string, and an
// index. The index corolates to the elasticsearch index.
func NewEvent(payload, index string) *Event {
	trimmed := strings.Trim(payload, " \t\n")
	dataType := Raw

	if len(trimmed) >= 2 {
		if strings.IndexByte("{[", trimmed[0]) != -1 && strings.IndexByte("}]", trimmed[len(trimmed)-1]) != -1 {
			dataType = Json
		} else if trimmed[0] == '<' && trimmed[len(trimmed)-1] == '>' {
			dataType = Xml
		}
	}

	return &Event{
		Data:     payload,
		DataType: dataType,
		Index:    index,
	}

}

// ToString returns a logging friendly representation of the event (doesn't include payload)
func (event *Event) ToString() string {
	return "EVENT(Index: \"" + event.Index + "\", Type: \"" + event.DataType + "\")"
}
