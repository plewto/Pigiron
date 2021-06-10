package main

import (
	"fmt"
	"time"

	"github.com/plewto/pigiron/midi"
	"github.com/plewto/pigiron/op"
)

func main() {
	fmt.Println("Pigiron.main()")
	
	op1, _ := op.NewMIDIInput("in", "Arturia")
	cf, _ := op.NewOperator("ChannelFilter", "filter")
	dist, _ := op.NewOperator("Distributor", "dist")
	mon, _ := op.NewOperator("Monitor", "mon")	
	op2, _ := op.NewMIDIOutput("out", "E-MU Xmidi 2x2 MIDI 1")

	op1.Connect(cf)
	cf.Connect(dist)
	dist.Connect(mon)
	dist.Connect(op2)

	// midi.DumpDevices()
	// fmt.Println("--------------------------------")
	op1.PrintTree()
	fmt.Println()
	fmt.Println(cf.Info())

	for {
		op.ProcessInputs()
		time.Sleep(1 * time.Millisecond)
	}
	
	Cleanup()

}

func Ignore(values ...interface{}) {}


func Cleanup() {
	fmt.Println("pigiron.Cleanup()")
	midi.Cleanup()
	op.Cleanup()
	
}
