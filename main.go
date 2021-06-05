package main

import (
	"fmt"
)

func main() {
	fmt.Println("Pigiron.main()")

	Cleanup()

}

func Ignore(values ...interface{}) {}


func Cleanup() {
	fmt.Println("pigiron.Cleanup() executes")
}
