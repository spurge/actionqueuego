package actionserver

import (
	"bufio"
	"net"
)

type ActionServer struct {
	joins   chan net.Conn
	closes  chan bool
	writes  chan []byte
	clients []*ActionClient
}

func NewActionServer() *ActionServer {
	server := &ActionServer{
		make(chan net.Conn),
		make(chan bool),
		make(chan []byte),
		make([]*ActionClient, 0),
	}

	return server
}

func (as ActionServer) Listen() {
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
		case data := <-as.writes:
			as.write(data)
		}
	}
}

func (as ActionServer) Join(conn net.Conn) {
	as.joins <- conn
}

func (as ActionServer) Close() {
	as.closes <- true
}

func (as ActionServer) close() {
	for _, client := range as.clients {
		client.Close()
	}
}

func (as ActionServer) Write(data []byte) {
	as.writes <- data
}

func (as ActionServer) write(data []byte) {
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
