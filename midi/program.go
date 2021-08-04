package midi


import (
	"github.com/rakyll/portmidi"
)


type ProgramBank interface {
	ChannelSelector
	ProgramRange() (byte, byte)
	CurrentProgram() byte
	ChangeProgram(event portmidi.Event)     // ignore out of bounds values
	String() string
}



	
