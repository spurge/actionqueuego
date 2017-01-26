Action Queue (go version)
=========================

[![Build Status](https://semaphoreci.com/api/v1/projects/e86409eb-a13b-4ce7-ad0b-43987c0e3dc5/1151050/shields_badge.svg)](https://semaphoreci.com/houseagency/actionqueue-go)

Implements the Action Queue in go.

## 1. Create a queue

```go
queue, err := NewActionQueue(filename)
```

## 2. Add an action to the queue

```go
queue.AddAction("{\"some-property\":\"with-a-value\"}")
```

## 3. Read all actions

```go
callback := func(entry *ActionEntry, err error) {
  // entry.pos
  // entry.tim
  // entry.def
}

count, err := queue.ReadHistory(callback, 0, -1)
```

## 4. Tailing actions

```go
done := make(chan bool)

callback := func(entry *ActionEntry, err error) {
  // entry.pos
  // entry.tim
  // entry.def
}

go queue.TailHistory(callback, 0, done)
done <- true
```
