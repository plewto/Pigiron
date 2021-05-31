package main

import (
	"fmt"

	"github.com/plewto/pigiron/op"
)

func main() {
	fmt.Println("Pigiron.main()")
	a, _ := op.MakeOperator("Null", "Alpha")
	fmt.Println(a)
}
