package queue

import (
	"container/list"
	"sync"
)

var Queue *queue

func init() {
	Queue = &queue{}
	Queue.list = list.New().Init()
}

type queue struct {
	Lock sync.RWMutex
	List *list.List
}

func (q *queue) Enqueue(responce interface{}) {
	q.Lock.Lock()
	defer q.Lock.Unlock()
	q.List.PushBack(responce)
}

func (q *queue) Dequeue() interface{} {
	q.Lock.Lock()
	defer q.Lock.Unlock()
	element := q.List.Front()
	q.List.Remove(element)
	return element.Value
}

func (q *queue) Len() int {
	q.Lock.Lock()
	defer q.Lock.Unlock()
	return q.List.Len()
}
