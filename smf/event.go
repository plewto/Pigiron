package smf

import (
	"github.com/rakyll/portmidi"
)


type Event interface {
	Delta() int
	Events() []*portmidi.Event
}
	
