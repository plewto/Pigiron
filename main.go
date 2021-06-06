package main

import (
	"fmt"

	"github.com/plewto/pigiron/midi"
)

func main() {
	fmt.Println("Pigiron.main()")
	midi.DumpDevices()
	id, err := midi.GetOutputID("MIDI 2")

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(id)
	Cleanup()

}

func Ignore(values ...interface{}) {}


func Cleanup() {
	fmt.Println("pigiron.Cleanup()")
	midi.Cleanup()
	
}
