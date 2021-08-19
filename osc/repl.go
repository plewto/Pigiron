package osc

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
	goosc "github.com/hypebeast/go-osc/osc"
	"github.com/plewto/pigiron/config"
	"github.com/plewto/pigiron/piglog"
	"github.com/plewto/pigiron/macro"
)

const (
	COMMENT = "#"
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

// Prompt() displays terminal prompt.
//
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
	s = strings.Split(s, COMMENT)[0]
	command, rawArgs := splitCommand(s)
	acc := make([]string, 0, len(rawArgs))
	for _, a := range strings.Split(rawArgs, ",") {
		arg := strings.TrimSpace(a)
		acc = append(acc, arg)
	}
	return command, acc
}

// Eval() evaluates REPL commands
//
func Eval(command string, args []string) {
	switch {
	case command == "":  // ignore blank lines
	case command == "#": // ignore comment lines
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



// REPL() enters the interactive command loop.
//
func REPL() {

	// Filter input characters as BUG 013 fix.
	//
	accept := func(s string) bool {
		for _, c := range strings.TrimSpace(s) {
			b := byte(c)
			if b < 0x20 || b > 0x7E {
				return false
			}
		}
		return true
	}

	
	for {
		fmt.Print(config.GlobalParameters.TextColor)
		Prompt()
		raw := Read()
		if accept(raw) {
			piglog.Log(fmt.Sprintf("CMD  : %s", raw))
			command, args := parse(raw)
			filename, flag := batchFileExist(command)
			if flag {
				BatchLoad(filename)
			} else {
				if macro.IsMacro(command) {
					expanded, err := macro.Expand(command, args)
					if err != nil {
						fmt.Printf("ERROR: %v\n", err)
						continue
					}
					command, args = parse(expanded)
				}
				Eval(command, args)
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
}
		

// /pig/eval command [,arg1, arg2, ...]
// 
func remoteEval(msg *goosc.Message) ([]string, error) {
	var err error
	for _, raw := range msg.Arguments {
		command, args := parse(fmt.Sprintf("%v", raw))
		command = strings.TrimRight(command, ",")
		Eval(command, args)
		time.Sleep(10 * time.Millisecond)
	}
	return empty, err
}
	
