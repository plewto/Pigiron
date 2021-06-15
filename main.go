package main

import (
	"fmt"
	//"time"
	"github.com/plewto/pigiron/config"
	"github.com/plewto/pigiron/osc"
	"github.com/plewto/pigiron/midi"
	"github.com/plewto/pigiron/op"

)

func main() {
	fmt.Println("Pigiron.main()")
	config.DumpGlobalParameters()

	osc.AckGlobal("/pig/foo", []string{"Alpha", "Beta", "Gamma"})
	
	Cleanup()
}

func Ignore(values ...interface{}) {}


func Cleanup() {
	fmt.Println("pigiron.Cleanup()")
	midi.Cleanup()
	op.Cleanup()
	
}
