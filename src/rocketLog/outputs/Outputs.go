package outputs

import "rocketLog"

type Output interface {
	Write(data *rocketLog.Event)
	Close()
}