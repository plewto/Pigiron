package main

import (
	"fmt"
	"bufio"
	"os"
	"strings"
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
	
	
func repl() {
	reader := bufio.NewReader(os.Stdin)
	ip := config.GlobalParameters.OSCServerHost
	port := int(config.GlobalParameters.OSCServerPort)
	root := config.GlobalParameters.OSCServerRoot
	client := goosc.NewClient(ip, port)
	for {
		fmt.Print(prompt)
		raw, _ := reader.ReadString('\n')
		raw = strings.TrimSpace(raw)
		command, args := parseAddress(raw)
		address := fmt.Sprintf("/%s/%s", root, command)
		msg := goosc.NewMessage(address)
		for _, arg := range strings.Split(args, ",") {
			msg.Append(strings.TrimSpace(arg))
		}
		client.Send(msg)
	}
}
		
		
		
