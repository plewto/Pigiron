package main

import (
	"fmt"
	"bufio"
	"os"
	"strings"
	
	goosc "github.com/hypebeast/go-osc/osc"

	
	"github.com/plewto/pigiron/config"
	"github.com/plewto/pigiron/osc"
	
	
)

const remark string = "#"

var (
	replReader = bufio.NewReader(os.Stdin)
	replClient osc.PigClient
)

func prompt() {
	root := config.GlobalParameters.OSCServerRoot
	fmt.Printf("/%s/ ", root)
}

func repl() {
	root := config.GlobalParameters.REPLRoot
	host := config.GlobalParameters.REPLHost
	port := int(config.GlobalParameters.REPLPort)
	replClient = osc.NewClient(host, port, root, "")
	for {
		prompt()
		s := read()
		eval(s)
	}
	
}


func read() string {
	s, err := replReader.ReadString('\n')
	if err != nil {
		fmt.Printf("\nERROR: $s\n", err)
		s = ""
	}
	return s
}


func dispatch(command string, args string) {
	switch command {
	case "":
		// ignore blank line
	case "#" :
		// ignore comment
	case "batch": 
		// ReadBatchFile(args, client)
	default: // transmit OSC
		root := config.GlobalParameters.REPLRoot
		address := fmt.Sprintf("/%s/%s", root, command)
		msg := goosc.NewMessage(address)
		for _, arg := range strings.Split(args, ",") {
			msg.Append(strings.TrimSpace(arg))
		}
		replClient.Send(msg)
	}
}

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


func eval(s string) {
	raw := strings.TrimSpace(s)
	command, args := parseAddress(raw)
	dispatch(command, args)
}
