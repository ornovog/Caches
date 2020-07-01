package caches

import (
	"container/list"
	"sync"
)

type lruQueue struct {
	lRUList list.List
	pointersArray[] *list.Element
	rWM sync.RWMutex
}

func (q *lruQueue) Init(length uint32){
	q.pointersArray = make([]*list.Element,length)
}

func (q *lruQueue) Update(index uint32){
	q.rWM.Lock()
	if  q.pointersArray[index] == nil{
		element := q.lRUList.PushFront(index)
		q.pointersArray[index] = element
	} else {
		element := q.pointersArray[index]
		q.lRUList.MoveToFront(element)
	}
	q.rWM.Unlock()
}

func (q *lruQueue) Back() uint32{
	q.rWM.RLock()
	val := q.lRUList.Back().Value.(uint32)
	q.rWM.RUnlock()
	return val
}