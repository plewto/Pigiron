package osc

import (
	"fmt"
	"strings"
	"os"
	"bufio"
	"time"
	"io/ioutil"
	"path/filepath"
	
	goosc "github.com/hypebeast/go-osc/osc"
	"github.com/plewto/pigiron/config"
	
)

var (
	internalClient *goosc.Client
	reader *bufio.Reader

	// if true exit batch mode
	batchError bool = false

	// true while in batch mode.
	// OSC REPL client output suppressed while true
	inBatchMode bool = false
)


func init() {
	host := config.GlobalParameters.OSCServerHost
	port := int(config.GlobalParameters.OSCServerPort)
	internalClient = goosc.NewClient(host, port)
	reader = bufio.NewReader(os.Stdin)
}
	

func sleep(n int) {
	time.Sleep(time.Duration(n) * time.Millisecond)
}

// SubUserHome replaces special characters in filename.
//
// Two special cases are defined:
//
// 1) If the first filename character is '~', the filename is returned
//    relative to the user's home directory.
// 2) If the first character is '!', the filename is returned relative to
//    the configuration directory
//
//    foo    --> foo
//    ~/foo  --> /home/<user>/foo
//    !/foo  --> /home/<user>/.config/pigiron/foo
//    
//   
func SubUserHome(filename string) string {
	if len(filename) == 0 {
		return filename
	}
	prefix := filename[0]
	switch prefix {
	case byte('~'):
		home, _ := os.UserHomeDir()
		filename = filepath.Join(home, filename[1:])
	case byte('!'):
		cfig, _ := os.UserConfigDir()
		filename = filepath.Join(cfig, "pigiron", filename[1:])
	default:
		// ignore
	}
	return filename

}


func Prompt() {
	root := config.GlobalParameters.OSCServerRoot
	fmt.Printf("\n/%s: ", root)
}


func Read() string {
	s, _ := reader.ReadString('\n')
	return s
}


// splits command from arguments  
func splitCommand(s string)(string, string) {
	s = strings.TrimSpace(s)
	pos := strings.Index(s, " ")
	if pos < 0 {
		return s, ""
	}
	command, args := s[:pos], strings.TrimSpace(s[pos:])
	return command, args
}
	

func parse(s string)(string, []string) {
	command, rawArgs := splitCommand(s)
	acc := make([]string, 0, len(rawArgs))
	for _, a := range strings.Split(rawArgs, ",") {
		acc = append(acc, strings.TrimSpace(a))
	}
	return command, acc
}


func Eval(command string, args []string) {
	switch command {
	case "":  // ignore blank lines
	case "#": // ignore comment lines
	default:
		root := config.GlobalParameters.OSCServerRoot
		address := fmt.Sprintf("/%s/%s", root, command)
		msg := goosc.NewMessage(address)
		for _, s := range args {
			msg .Append(s)
		}
		internalClient.Send(msg)
	}
}


func printBatchError(filename string, err error) {
	fmt.Print(config.GlobalParameters.ErrorColor)
	fmt.Printf("Can not read batch file '%s'\n", filename)
	fmt.Printf("%s\n", err)
}

func BatchLoad(filename string) error {
	batchError = false
	inBatchMode = true
	filename = SubUserHome(filename)
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
		

func REPL() {
	for {
		fmt.Print(config.GlobalParameters.TextColor)
		Prompt()
		raw := Read()
		command, args := parse(raw)
		Eval(command, args)
		time.Sleep(10 * time.Millisecond)
	}
}
		
