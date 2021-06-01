package main

import (
	"fmt"

	//"github.com/plewto/pigiron/op"
	"github.com/plewto/pigiron/midi"
)

func main() {
	fmt.Println("Pigiron.main()")
	midi.DumpDevices()
	dev, err := midi.GetInputDevice("E-MU Xmidi 2x2 MIDI 2")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(dev)
	}
	Cleanup()
}


func Cleanup() {
	fmt.Println("pigiron.Cleanup() executes")
	midi.Cleanup()
}
