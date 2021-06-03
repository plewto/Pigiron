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
	mon, err := op.MakeOperator("Monitor", "mon")
	a.Connect(mon)

	if err != nil {
		panic(err)
	}
	
	// fmt.Println(a.Info())
	// fmt.Println("--------------------------------")
	// fmt.Println(mon.Info())
	// fmt.Println("--------------------------------")
	// a.PrintTree()

	fmt.Println("Ready....")
	
	for {}
	
	Cleanup()

}

func Ignore(values ...interface{}) {}


func Cleanup() {
	fmt.Println("pigiron.Cleanup() executes")
	midi.Cleanup()
}
