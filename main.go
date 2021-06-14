package main

import (
	"fmt"
	//"time"
	_ "github.com/plewto/pigiron/config"
	"github.com/plewto/pigiron/midi"
	"github.com/plewto/pigiron/op"
)

func main() {
	fmt.Println("Pigiron.main()")


	
	
	Cleanup()

}

func Ignore(values ...interface{}) {}


func Cleanup() {
	fmt.Println("pigiron.Cleanup()")
	midi.Cleanup()
	op.Cleanup()
	
}
