package actionqueue

import (
  "bufio"
  "encoding/json"
  "os"
  "time"
)

type ActionQueue struct {
  filename string
  position int
  file os.File
  writer bufio.Writer
}

type ActionEntry struct {
  pos int
  def string
  tim string
}

type HistoryCallback func(action ActionEntry, err error)

func NewActionQueue(filename string) (ActionQueue, error) {
  file, err := os.OpenFile(
    filename,
    os.O_CREATE|os.O_RDWR|os.O_APPEND,
    0660,
  )

  return ActionQueue{
    filename,
    0,
    *file,
    *bufio.NewWriter(file),
  }, err
}

func (q ActionQueue) AddAction(def string) (int, error) {
  data, err := json.Marshal(map[string]interface{}{
    "pos": q.position,
    "def": def,
    "tim": time.Now().String(),
  })

  if err != nil {
    return q.position, err
  }

  if _, err := q.writer.Write(data); err != nil {
    return q.position, err
  }

  if err := q.writer.WriteByte('\n'); err != nil {
    return q.position, err
  }

  q.position++

  defer q.writer.Flush()

  return q.position, nil
}

func (q ActionQueue) ReadHistory(cb HistoryCallback, from int, to int) (int, error) {
  file, err := os.Open(q.filename)

  if err != nil {
    return 0, err
  }

  count := 0
  position := 0
  scanner := bufio.NewScanner(file)

  for scanner.Scan() {
    if (position >= from && (to < 0 || position <= to)) {
      var entry ActionEntry
      err := json.Unmarshal(scanner.Bytes(), &entry)
      cb(entry, err)
      count++
    }

    position++;
  }

  return count, scanner.Err()
}

func (q ActionQueue) Close() {
  q.file.Close()
}
