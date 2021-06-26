package osc

import (
	"fmt"
	goosc "github.com/hypebeast/go-osc/osc"
	
)

type PigServer interface {
	Root() string
	SetRoot(string)
	//Clients() []PigClient
	//AddClient(client PigClient)
	//RemoveAllClients()
	GetResponder() Responder
	GetREPLResponder() Responder
	AddMsgHandler(address string, handler func(msg *goosc.Message))
	ListenAndServe()
	IP() string
	Port() int
	Close()
	Commands() []string
}

type OSCServer struct {
	backingServer *goosc.Server
	dispatcher *goosc.StandardDispatcher
	root string
	//clients []PigClient
	responder Responder
	replResponder Responder
	ip string
	port int
	commands []string
}



func NewServer(ip string, port int, root string) PigServer {
	server := new(OSCServer)
	server.ip = ip
	server.port = port
	server.root = root
	//server.clients = make([]PigClient, 0, 4)
	server.responder = globalResponder
	server.replResponder = replResponder
	addr := fmt.Sprintf("%s:%d", ip, port)
	server.dispatcher = goosc.NewStandardDispatcher()
	server.backingServer = &goosc.Server {
		Addr: addr,
		Dispatcher: server.dispatcher,
	}
	server.commands = make([]string, 0, 16)
	return server
}

func (s *OSCServer) GetResponder() Responder {
	return s.responder
}

func (s *OSCServer) GetREPLResponder() Responder {
	return s.replResponder
}

func (s *OSCServer) AddMsgHandler(address string, handler func(msg *goosc.Message)) {
	s.dispatcher.AddMsgHandler(address, handler)
	s.commands = append(s.commands, address)
}

func (s *OSCServer) Commands() []string {
	return s.commands
}

func (s *OSCServer) ListenAndServe() {
	fmt.Printf("OSC Listening: %s:%d  /%s\n", s.ip, s.port, s.root)
	go s.backingServer.ListenAndServe()
}

func (s *OSCServer) Root() string {
	return s.root
}

func (s *OSCServer) SetRoot(root string) {
	s.root = root
}

	
func (s *OSCServer) IP() string {
	return s.ip
}

func (s *OSCServer) Port() int {
	return s.port
}

func (s *OSCServer) Close() {
	s.backingServer.CloseConnection()
}


func AddOSCHandler(s PigServer, command string, handler func(*goosc.Message)([]string, error)) {
	address := fmt.Sprintf("/%s/%s", s.Root(), command)
	var result = func(msg *goosc.Message) {
		status, err := handler(msg)
		if err != nil {
			s.GetResponder().Error(address, status, err)
			s.GetREPLResponder().Error(address, status, err)
		} else {
			s.GetResponder().Ack(address, status)
			s.GetREPLResponder().Ack(address, status)
		}
	}
	s.AddMsgHandler(address, result)
}
