package main

import (
	"fmt"

	midi "github.com/plewto/pigiron/midi"
	op "github.com/plewto/pigiron/op"
	

)

func main() {
	fmt.Println("Pigiron.main()")
	midi.DumpDevices()
	fmt.Println("")
	a, _ := op.MakeMIDIInput("a", "Arturia")
	c, _ := op.MakeOperator("ChannelFilter", "filter")
	mon, _ := op.MakeOperator("Monitor", "monitor")

	// Enable all channels
	for i:=1; i<17; i++ {
		c.SelectChannel(i)
	}

	a.Connect(c)
	c.Connect(mon)
	
	fmt.Println(c)

	
	// fmt.Println("Ready....")
	// for {}

	
	
	Cleanup()

}

func Ignore(values ...interface{}) {}


func Cleanup() {
	fmt.Println("pigiron.Cleanup() executes")
	op.Cleanup()
	midi.Cleanup()
	
}
