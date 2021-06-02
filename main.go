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
	a, err := op.MakeMIDIInput("Alpha", "E-MU")
	b, err := op.MakeMIDIInput("Beta", "E-MU")

	fmt.Println(a.Info())
	fmt.Println("------------------------")
	fmt.Println(b.Info())

	fmt.Printf("a is b ?  : %v\n", a == b)
	
	Ignore(err)
	
	Cleanup()

}

func Ignore(values ...interface{}) {}


func Cleanup() {
	fmt.Println("pigiron.Cleanup() executes")
	midi.Cleanup()
}
