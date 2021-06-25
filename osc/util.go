package osc

import (
	"fmt"
	"os"
	"path/filepath"
	"github.com/plewto/pigiron/config"
)


func StringSlice(values ...interface{}) []string {
	acc := make([]string, len(values))
	for i, v := range values {
		acc[i] = fmt.Sprintf("%v", v)
	}
	return acc
}


func SubUserHome(filename string) string {
	if len(filename) > 0 && filename[0] == byte('~') {
		home, _ := os.UserHomeDir()
		filename = filepath.Join(home, filename[1:])
	}
	return filename
}
		
		


func Prompt() {
	root := config.GlobalParameters.OSCServerRoot
	fmt.Printf("\n/%s: ", root)
}
