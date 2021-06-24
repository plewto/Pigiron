package osc

import (
	"fmt"
	goosc "github.com/hypebeast/go-osc/osc"
	
)

type PigServer interface {
	Root() string
	SetRoot(string)
	// Client() PigClient
	// SetClient(PigClient)
	Clients() []PigClient
	AddClient(client PigClient)
	AddMsgHandler(command string, handler func(msg *goosc.Message))
	ListenAndServe()
	IP() string
	Port() int
	Close()
}

type OSCServer struct {
	backingServer *goosc.Server
	dispatcher *goosc.StandardDispatcher
	root string
	clients []PigClient
	ip string
	port int
}


func NewServer(ip string, port int, root string) PigServer {
	server := new(OSCServer)
	server.ip = ip
	server.port = port
	server.root = root
	//server.client = globalClient
	server.clients = make([]PigClient, 0, 4)
	addr := fmt.Sprintf("%s:%d", ip, port)
	server.dispatcher = goosc.NewStandardDispatcher()
	server.backingServer = &goosc.Server {
		Addr: addr,
		Dispatcher: server.dispatcher,
	}
	return server
}

func (s *OSCServer) AddMsgHandler(command string, handler func(msg *goosc.Message)) {
	addr := fmt.Sprintf("/%s/%s", s.root, command)
	s.dispatcher.AddMsgHandler(addr, handler)
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

// func (s *OSCServer) Client() PigClient {
// 	return s.client
// }

// func (s *OSCServer) SetClient(client PigClient) {
// 	s.client = client
// }

func (s *OSCServer) Clients() []PigClient {
	return s.clients
}

func (s *OSCServer) AddClient(client PigClient) {
	s.clients = append(s.clients, client)
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


func AddOSCHandler(s PigServer, address string, handler func(*goosc.Message)([]string, error)) {
	address = fmt.Sprintf("/%s/%s", s.Root(), address)
	var result = func(msg *goosc.Message) {
		status, err := handler(msg)
		if err != nil {
			for _, c := range s.Clients() {
				c.Error(address, status)
			}
		} else {
			for _, c := range s.Clients() {
				c.Ack(address, status)
			}
		}
	}
	s.AddMsgHandler(address, result)
}
