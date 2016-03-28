package outputs

import "github.com/hamstak/rocketlog/event"

type Output interface {
	Write(data *event.Event)
	Close()
}
