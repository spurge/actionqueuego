package actionqueue

import (
  "fmt"
  "os"
  "testing"
)

var filename string = "/tmp/action-queue-test"

func cleanup(queue ActionQueue) {
  queue.Close()
  os.Remove(filename)
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

  queue.AddAction("Action 1")
  queue.AddAction("Action 2")
  queue.AddAction("Action 3")

  callback := func (entry ActionEntry, err error) {
    tested++

    if entry.pos + 1 != tested {
      t.Error(fmt.Sprintf("Position %d != %d", entry.pos + 1, tested))
    }

    if entry.def != fmt.Sprintf("Action %d", tested) {
      t.Error(fmt.Sprintf("Action %d != %s", tested, entry.def))
    }

    fmt.Println(entry)

    /*tim, err := time.Parse(
      "2017-01-25T18:00:46.52576271+01:00",
      entry.tim,
    )

    if err != nil {
      t.Error(err)
    }*/
  }

  count, err := queue.ReadHistory(callback, 0, -1)

  if err != nil {
    t.Fatal(err)
  }

  if count != 3 || tested != 3 {
    t.Fatal(fmt.Sprintf("Count 3 != %d != %d", count, tested))
  }
}
