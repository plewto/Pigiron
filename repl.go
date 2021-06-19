package main

import (
	"fmt"
	"bufio"
	"os"
	"strings"
	"path/filepath"
	"io/ioutil"
	"time"
	
	goosc "github.com/hypebeast/go-osc/osc"
	"github.com/plewto/pigiron/config"
)

var prompt = "PIG: "

func parseAddress(s string) (command string, args string) {
	pos := strings.Index(s, " ")

	if pos == -1 {
		command = strings.TrimSpace(s)
		args = ""
	} else {
		command = s[:pos]
		args = strings.TrimSpace(s[pos:])
	}
	return command, args
}
	

// repl provides local interactive control.
//
// At the PIG: prompt the user may directly enter OSC commands.
// The format is almost identical to remote commands with the following
// exceptions:
//
//   1) The osc address root is not required, it is automatically supplied
//      i.e.  enter 'exit'  not /pig/exit
//
//   2) Arguments to the command are separated by commas.
//      i.e.  new-operator Monitor, alpha
//
func repl() {
	reader := bufio.NewReader(os.Stdin)
	ip := config.GlobalParameters.OSCServerHost
	port := int(config.GlobalParameters.OSCServerPort)
	client := goosc.NewClient(ip, port)
	for {
		fmt.Print(prompt)
		s := read(reader)
		eval(s, client)
	}
}


func read(reader *bufio.Reader) string {
	s, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		s = ""
	}
	return s
}


func dispatch(command string, args string, client *goosc.Client) {
	switch command {
	case "":
		// ignore blank line
	case "#" :
		// ignore comment
	case "batch": 
		readBatchFile(args, client)
	default: // transmit OSC
		root := config.GlobalParameters.OSCServerRoot
		address := fmt.Sprintf("/%s/%s", root, command)
		msg := goosc.NewMessage(address)
		for _, arg := range strings.Split(args, ",") {
			msg.Append(strings.TrimSpace(arg))
		}
		client.Send(msg)
	}
}

		
func eval(s string, client *goosc.Client) {
	raw := strings.TrimSpace(s)
	command, args := parseAddress(raw)
	dispatch(command, args, client)
}

func readBatchFile(filename string, client *goosc.Client) {
	delay := 100*time.Millisecond
	if len(filename) > 1 && filename[0:2] == "~/" {
		home, _ := os.UserHomeDir()
		filename = filepath.Join(home, filename[2:])
		
	}
	lines, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("ERROR: Can not load batch file '%s'\n", filename)
		fmt.Printf("ERROR: %s\n", err)
	} else {
		for _, line := range strings.Split(string(lines), "\n") {
			eval(line, client)
			time.Sleep(delay)
		}
	}
}
