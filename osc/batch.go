package osc

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"github.com/plewto/pigiron/config"
	"github.com/plewto/pigiron/pigpath"
	"github.com/plewto/pigiron/piglog"
)


func batchFileExist(name string) (filename string, flag bool) {
	if len(name) == 0 {
		return "", false
	}
	filename = pigpath.Join(config.GlobalParameters.BatchDirectory, name)
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return filename, false
	} else {
		return filename, true
	}
}


func printBatchError(filename string, err error) {
	fmt.Print(config.GlobalParameters.ErrorColor)
	fmt.Printf("Can not read batch file '%s'\n", filename)
	fmt.Printf("%s\n", err)
}

// BatchLoad() loads batch file.
// A batch file is a sequence of osc commands with identical syntax to
// interactive commands.   Lines beginning with # are ignored.
//
// The filename argument may begin with the special characters:
//    ~/  file is relative to the user's home directory.
//    !/  file is relative to the configuration directory.
//
// Returns non-nil error if the file could not be read.
//
func BatchLoad(filename string) error {
	batchError = false
	inBatchMode = true
	filename = pigpath.SubSpecialDirectories(filename)
	file, err := os.Open(filename)
	if err != nil {
		printBatchError(filename, err)
		inBatchMode = false
		return err
	}
	defer file.Close()
	raw, err := ioutil.ReadAll(file)
	if err != nil {
		printBatchError(filename, err)
		inBatchMode = false
		return err
	}
	fmt.Printf("Loading batch file  '%s'\n", filename)
	lines := strings.Split(string(raw), "\n")
	for i, line := range lines {
		fmt.Printf("Batch [%3d]  %s\n", i+1, line)
		piglog.Log(fmt.Sprintf("BATCH: [line %3d] %s", i, line))
		command, args := parse(line)
		Eval(command, args)
		sleep(10)
		if batchError {
			break
		}
	}
	batchError = false
	inBatchMode = false
	return err
}
