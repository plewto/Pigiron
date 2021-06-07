package main

import (
	"fmt"

	"github.com/plewto/pigiron/midi"
	"github.com/plewto/pigiron/op"
)

func main() {
	fmt.Println("Pigiron.main()")
	midi.DumpDevices()

	dummy, _ := op.NewOperator("Dummy", "Alpha")

	fmt.Println(dummy.Info())
	Cleanup()

}

func Ignore(values ...interface{}) {}


func Cleanup() {
	fmt.Println("pigiron.Cleanup()")
	midi.Cleanup()
	
}
