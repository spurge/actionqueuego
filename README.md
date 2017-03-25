Action Queue (go version)
=========================

[![Build Status](https://semaphoreci.com/api/v1/spurge/actionqueuego/branches/master/shields_badge.svg)](https://semaphoreci.com/spurge/actionqueuego)

This is a very lightweight and non distributed data event stream library and server. The event stream is designed as a queue. You define you're queue with a filename (everything is written to disk) which the queue will append data to. Then you'll start adding data (Actions) to that queue. At the other end, you'll have some data consumers. Those consumers can consume any data in the stream and/or tailing.

## Setup

```bash
go get -u github.com/kardianos/govendor
govendor sync
```

## Test

`govendor test +local`

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
