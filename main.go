package main

import (
	"fmt"
	//"time"

	"github.com/plewto/pigiron/midi"
	"github.com/plewto/pigiron/op"
)

func main() {
	fmt.Println("Pigiron.main()")

	in, _ := op.NewMIDIInput("in", "Arturia")
	filter, _ := op.NewOperator("ChannelFilter", "filter")
	dist, _ := op.NewOperator("Distributor", "dist")
	mon, _ := op.NewOperator("Monitor", "mon")
	out, _ := op.NewMIDIOutput("out", "MIDI 1")

	in.Connect(filter)
	filter.Connect(dist)
	dist.Connect(mon)
	mon.Connect(out)

	in.PrintTree()
	
	// for {
	// 	op.ProcessInputs()
	// 	time.Sleep(1 * time.Millisecond)
	// }

	fmt.Println(dist.Info())
	
	Cleanup()

}

func Ignore(values ...interface{}) {}


func Cleanup() {
	fmt.Println("pigiron.Cleanup()")
	midi.Cleanup()
	op.Cleanup()
	
}
