package actionserver

import (
	"actionqueue.go/actionserver/mock_net"
	"fmt"
	"github.com/golang/mock/gomock"
	"testing"
)

func TestListenAndJoin(t *testing.T) {
	mockctrl := gomock.NewController(t)
	defer mockctrl.Finish()

	server := NewActionServer()

	go server.Listen()

	conn := mock_net.NewMockConn(mockctrl)

	gomock.InOrder(
		conn.EXPECT().Write([]byte("Some test 0")),
		conn.EXPECT().Write([]byte("Some test 1")),
		conn.EXPECT().Write([]byte("Some test 2")),
		conn.EXPECT().Write([]byte("Last test after close == false")),
	)

	server.join <- conn

	if len(server.clients) != 1 {
		t.Fatal(fmt.Sprintf("Clients 1 != %d", len(server.clients)))
	}

	for i := 0; i < 3; i++ {
		server.write <- []byte(fmt.Sprintf("Some test %d", i))
	}

	// Shall not close on false
	server.close <- false

	server.write <- []byte("Last test after close == false")

	// But close now...
	server.close <- true
}
