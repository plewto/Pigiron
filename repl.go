package main

import (
	"fmt"
	"bufio"
	"os"
	"io/ioutil"
	"strings"
	"path/filepath"
	"time"
	goosc "github.com/hypebeast/go-osc/osc"
	"github.com/plewto/pigiron/config"
	"github.com/plewto/pigiron/osc"
)

const remark string = "#"

var (
	replReader = bufio.NewReader(os.Stdin)
	replClient osc.PigClient
)

// subUserHome substitutes leading ~ charater for user home directory.
//
func subUserHome(filename string) string {
	result := filename
	if len(filename) > 0 && string(filename[0]) == "~" {
		home, _ := os.UserHomeDir()
		result = filepath.Join(home, filename[1:])
	}
	return result
}

func printBar() {
	fmt.Println("------------------------")
}


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
		LoadBatchFile(args)
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


func LoadBatchFile(filename string) {
	osc.ClearError()
	filename = subUserHome(filename)
	fmt.Print(config.GlobalParameters.TextColor)
	printBar()
	fmt.Printf("Loading batch file: '%s'\n", filename)
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Print(config.GlobalParameters.ErrorColor)
		fmt.Printf("ERROR: Can not open batch file: '%s'\n", filename)
		fmt.Printf("ERROR: %s\n", err)
		fmt.Print(config.GlobalParameters.TextColor)
	} else {
		lines := strings.Split(string(buf), "\n")
		for i, line := range lines {
			eval(line)
			time.Sleep(10*time.Millisecond)
			if osc.OSCError() {
				fmt.Print(config.GlobalParameters.ErrorColor)
				fmt.Printf("ERROR: batch file line %d\n", i)
				fmt.Printf("ERROR: %s\n", line)
				break
			}
			
		}
	}
	time.Sleep(10*time.Millisecond)
	fmt.Print(config.GlobalParameters.TextColor)
	fmt.Println()
	prompt()
}


	
		
