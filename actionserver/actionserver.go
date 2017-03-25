package actionserver

import (
	"bufio"
	"net"
	//"net/http"
	"log"
)

type ActionServer struct {
	joins       chan net.Conn
	closes      chan bool
	writes      chan []byte
	clients     []*ActionClient
	netListener net.Listener
}

func NewActionServer(netListener net.Listener) *ActionServer {
	server := &ActionServer{
		make(chan net.Conn),
		make(chan bool),
		make(chan []byte),
		make([]*ActionClient, 0),
		netListener,
	}

	return server
}

func (as *ActionServer) Listen() {
	go as.acceptingClients()

loop:
	for {
		select {
		case close := <-as.closes:
			if close {
				as.close()
				break loop
			}
		case conn := <-as.joins:
			client := NewActionClient(conn)
			as.clients = append(as.clients, client)
			log.Println("Joined clients: ", len(as.clients))
		case data := <-as.writes:
			as.write(data)
		}
	}
}

func (as *ActionServer) acceptingClients() {
	for {
		conn, err := as.netListener.Accept()

		if err == nil {
			as.Join(conn)
		}
	}
}

func (as *ActionServer) Join(conn net.Conn) {
	as.joins <- conn
}

func (as *ActionServer) Close() {
	as.closes <- true
}

func (as *ActionServer) close() {
	log.Println("Closing clients: ", len(as.clients))
	for _, client := range as.clients {
		client.Close()
	}

	as.netListener.Close()
}

func (as *ActionServer) Write(data []byte) {
	as.writes <- data
}

func (as *ActionServer) write(data []byte) {
	log.Println("Writing clients: ", len(as.clients))
	for _, client := range as.clients {
		client.Write(data)
	}
}

type ActionClient struct {
	writes chan []byte
	closes chan bool
	conn   net.Conn
	writer *bufio.Writer
}

func NewActionClient(conn net.Conn) *ActionClient {
	client := &ActionClient{
		make(chan []byte),
		make(chan bool),
		conn,
		bufio.NewWriter(conn),
	}

	go client.listen()

	return client
}

func (c ActionClient) listen() {
loop:
	for {
		select {
		case close := <-c.closes:
			if close {
				c.close()
				break loop
			}
		case data := <-c.writes:
			c.write(data)
		}
	}
}

func (c ActionClient) Write(data []byte) {
	c.writes <- data
}

func (c ActionClient) write(data []byte) {
	c.writer.Write(data)
	defer c.writer.Flush()
}

func (c ActionClient) Close() {
	c.closes <- true
}

func (c ActionClient) close() {
	c.conn.Close()
}
