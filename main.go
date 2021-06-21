package main

import (
	"fmt"
	// "time"
	"os"
	
	"github.com/plewto/pigiron/config"
	"github.com/plewto/pigiron/osc"
	"github.com/plewto/pigiron/midi"
	"github.com/plewto/pigiron/op"

)


func main() {
	fmt.Println("Pigiron")
	config.DumpGlobalParameters()
	osc.Listen()
	go repl()
	
	for { // main loop
		if osc.Exit {
			Exit()
		}
	}
}

// func main() {
// 	fmt.Println("Main TEST")
// }


func Exit() {
	fmt.Println("Pigiron exit.")
	Cleanup()
	os.Exit(0)
}



func Cleanup() {
	midi.Cleanup()
	op.Cleanup()
	
}
