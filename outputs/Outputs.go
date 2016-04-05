package outputs

import "github.com/hamstak/rocketlog/event"

// Output is an endpoint for any rocketlog event to be written to.
type Output interface {
	Write(data *event.Event)
    ToString() string
	Close()
}
