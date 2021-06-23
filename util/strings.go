package util

import (
	"fmt"
	"github.com/plewto/pigiron/config"
)


func StringSlice(values ...interface{}) []string {
	acc := make([]string, len(values))
	for i, v := range values {
		acc[i] = fmt.Sprintf("%v", v)
	}
	return acc
}


func Prompt() {
	root := config.GlobalParameters.OSCServerRoot
	fmt.Printf("\n/%s/ ", root)
}
