package outputs

import "rocketlog/events"

type Output interface {
	Write(data *event.Event)
	Close()
}