package main

import (
	"fmt"

	"github.com/plewto/pigiron/op"
)

func main() {
	fmt.Println("Pigiron.main()")
	op1, _ := op.MakeOperator("Null", "Alpha")
	fmt.Println(op1.Info())
}
