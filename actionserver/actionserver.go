package actionserver

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
)

type ActionServer struct {
	join    chan net.Conn
	close   chan bool
	write   chan []byte
	clients []*ActionClient
}

func NewActionServer() *ActionServer {
	return &ActionServer{
		make(chan net.Conn),
		make(chan bool),
		make(chan []byte),
		make([]*ActionClient, 0),
	}
}

func (as ActionServer) Listen() {
	for {
		select {
		case client := <-as.join:
			as.Join(client)
			break
		case close := <-as.close:
			return
		case data := <-as.write:
			as.Write(data)
			break
		}
	}
}

func (as ActionServer) Join(conn net.Conn) {
	client := NewActionClient(conn)
	as.clients = append(as.clients, client)
	go client.Listen()
}

func (as ActionServer) Close() {
	for _, client := range as.clients {
		client.close <- true
	}
}

func (as ActionServer) Write(data []byte) {
	for _, client := range as.clients {
		client.out <- data
	}
}

type ActionClient struct {
	out    chan []byte
	close  chan bool
	conn   *net.Conn
	writer *bufio.Writer
}

func NewActionClient(conn net.Conn) *ActionClient {
	return &ActionClient{
		make(chan []byte),
		make(chan bool),
		&conn,
		bufio.NewWriter(conn),
	}
}

func (c ActionClient) Listen() {
	for {
		select {
		case close := <-c.close:
			c.Close()
			return
		case data := <-c.out:
			c.Write(data)
			break
		}
	}
}

func (c ActionClient) Write(data []byte) {
	for _, b := range data {
		c.writer.WriteByte(b)
	}

	defer c.writer.Flush()
}

func (c ActionClient) Close() {
	//c.conn.
}
