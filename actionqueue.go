package actionqueue

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"time"
)

type ActionQueue struct {
	filename string
	pos      int
	file     *os.File
	writer   *bufio.Writer
}

type ActionEntry struct {
	pos int
	def string
	tim string
}

type HistoryCallback func(entry *ActionEntry, err error)

type fileReaderCallback func(reader *bufio.Reader) (count int, err error)

func NewActionQueue(filename string) (*ActionQueue, error) {
	file, err := os.OpenFile(
		filename,
		os.O_CREATE|os.O_RDWR|os.O_APPEND,
		0660,
	)

	return &ActionQueue{
		filename,
		0,
		file,
		bufio.NewWriter(file),
	}, err
}

func (q *ActionQueue) AddAction(def string) (int, error) {
	data, err := json.Marshal(map[string]interface{}{
		"def": def,
		"tim": time.Now().String(),
	})

	if err != nil {
		return q.pos, err
	}

	if _, err := q.writer.Write(data); err != nil {
		return q.pos, err
	}

	if err := q.writer.WriteByte('\n'); err != nil {
		return q.pos, err
	}

	q.pos++

	defer q.writer.Flush()

	return q.pos, nil
}

func (q *ActionQueue) ReadHistory(
	cb HistoryCallback,
	from int,
	to int,
) (int, error) {
	callback := func(reader *bufio.Reader) (int, error) {
		_, count := readLines(cb, reader, 0, from, to)
		return count, nil
	}

	return readFile(q.filename, callback)
}

func (q *ActionQueue) TailHistory(
	cb HistoryCallback,
	from int,
	done chan bool,
) (int, error) {
	callback := func(reader *bufio.Reader) (int, error) {
		count := 0
		pos := 0

	loop:
		for {
			select {
			case <-done:
				break loop
			default:
				p, c := readLines(cb, reader, pos, from, -1)

				count += c
				pos = p

				time.Sleep(1 * time.Millisecond)
			}
		}

		return count, nil
	}

	return readFile(q.filename, callback)
}

func (q *ActionQueue) Close() {
	q.file.Close()
}

func readFile(filename string, cb fileReaderCallback) (int, error) {
	file, err := os.Open(filename)
	defer file.Close()

	if err != nil {
		return 0, err
	}

	return cb(bufio.NewReader(file))
}

func readLines(
	cb HistoryCallback,
	reader *bufio.Reader,
	pos int,
	from int,
	to int,
) (int, int) {
	count := 0
	bytes, _, err := reader.ReadLine()

	for len(bytes) > 0 && err != io.EOF && (to < 0 || pos <= to) {
		if pos >= from {
			var data map[string]string
			err := json.Unmarshal(bytes, &data)
			entry := ActionEntry{
				pos,
				data["def"],
				data["tim"],
			}

			cb(&entry, err)
			count++
		}

		bytes, _, err = reader.ReadLine()
		pos++
	}

	return pos, count
}
