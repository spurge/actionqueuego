package actionqueue

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"
)

var filename string = "/tmp/action-queue-test"

func cleanup(queue *ActionQueue) {
	queue.Close()
	os.Remove(filename)
}

func addRandomActions(queue *ActionQueue, delay time.Duration) int {
	rand.Seed(time.Now().UnixNano())

	total := rand.Intn(100)

	for i := 0; i < total; i++ {
		queue.AddAction(fmt.Sprintf("Action %d", i))
		time.Sleep(delay)
	}

	return total
}

func testEntry(entry *ActionEntry, tested int, err error, t *testing.T) {
	if entry.pos != tested {
		t.Error(fmt.Sprintf("Position %d != %d", entry.pos, tested))
	}

	if entry.def != fmt.Sprintf("Action %d", tested) {
		t.Error(fmt.Sprintf("Action %d != %s", tested, entry.def))
	}

	if err != nil {
		t.Error(err)
	}
}

func TestAddAction(t *testing.T) {
	queue, err := NewActionQueue(filename)

	defer cleanup(queue)

	if err != nil {
		t.Fatal(err)
	}

	pos, err := queue.AddAction("{\"test\":\"value\"}")

	if err != nil {
		t.Fatal(err)
	}

	if pos <= 0 {
		t.Fatal("Position was", pos)
	}
}

func TestReadAllHistory(t *testing.T) {
	queue, _ := NewActionQueue(filename)
	tested := 0

	defer cleanup(queue)

	total := addRandomActions(queue, 0)

	callback := func(entry *ActionEntry, err error) {
		testEntry(entry, tested, err, t)
		tested++
	}

	count, err := queue.ReadHistory(callback, 0, -1)

	if err != nil {
		t.Fatal(err)
	}

	if count != total || tested != total {
		t.Fatal(fmt.Sprintf("Count %d != %d != %d", total, count, tested))
	}
}

func TestReadHistoryParts(t *testing.T) {
	queue, _ := NewActionQueue(filename)
	tested := 0

	defer cleanup(queue)

	total := 67
	begin := 23
	stop := 60

	for i := 0; i < total; i++ {
		queue.AddAction(fmt.Sprintf("Action %d", i))
	}

	callback := func(entry *ActionEntry, err error) {
		testEntry(entry, tested+begin, err, t)
		tested++
	}

	count, err := queue.ReadHistory(callback, begin, stop)

	if err != nil {
		t.Fatal(err)
	}

	if count != tested {
		t.Fatal(fmt.Sprintf("Count %d > %d != %d", total, count, tested))
	}
}

func TestTailHistory(t *testing.T) {
	queue, _ := NewActionQueue(filename)
	tested := 0

	defer cleanup(queue)

	done := make(chan bool)

	callback := func(entry *ActionEntry, err error) {
		testEntry(entry, tested, err, t)
		tested++
	}

	go queue.TailHistory(callback, 0, done)
	total := addRandomActions(queue, 10*time.Millisecond)
	done <- true

	if tested != total {
		t.Fatal(fmt.Sprintf("Tested %d != %d", total, tested))
	}
}
