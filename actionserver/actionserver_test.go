package actionserver

import (
	"actionqueue.go/actionserver/mock_net"
	"fmt"
	"github.com/golang/mock/gomock"
	"sync"
	"testing"
	"time"
)

func TestListenAndJoin(t *testing.T) {
	mockctrl := gomock.NewController(t)
	defer mockctrl.Finish()

	conn1 := mock_net.NewMockConn(mockctrl)
	conn2 := mock_net.NewMockConn(mockctrl)

	gomock.InOrder(
		conn1.EXPECT().
			Write([]byte("Some test 0")).
			Return(len([]byte("Some test 0")), nil),
		conn1.EXPECT().
			Write([]byte("Some test 1")).
			Return(len([]byte("Some test 1")), nil),
		conn1.EXPECT().
			Write([]byte("Some test 2")).
			Return(len([]byte("Some test 2")), nil),
		conn1.EXPECT().
			Write([]byte("Last test after close == false")).
			Return(len([]byte("Last test after close == false")), nil),
		conn1.EXPECT().Close().Return(nil),
	)

	gomock.InOrder(
		conn2.EXPECT().
			Write([]byte("Last test after close == false")).
			Return(len([]byte("Last test after close == false")), nil),
		conn2.EXPECT().Close().Return(nil),
	)

	var wg sync.WaitGroup
	wg.Add(2)
	server := NewActionServer()

	go func() {
		server.Listen()

		// Give the clients some time to close down
		time.Sleep(10 * time.Millisecond)

		wg.Done()
	}()

	go func() {
		server.Join(conn1)

		for i := 0; i < 3; i++ {
			server.Write([]byte(fmt.Sprintf("Some test %d", i)))
		}

		server.Join(conn2)

		// Shall not close on false
		server.closes <- false

		server.Write([]byte("Last test after close == false"))

		// But close now...
		server.Close()

		wg.Done()
	}()

	wg.Wait()
}
