package LRU

import (
	"Caches/mainMemory"
	"container/list"
)

type LruQueue struct {
	lRUList list.List
	pointersArray[] *list.Element
}

func (q *LruQueue) Init(length uint32){
	q.pointersArray = make([]*list.Element,length)
}

func (q *LruQueue) Update(index mainMemory.Address){
	if  q.pointersArray[index] == nil{
		element := q.lRUList.PushFront(index)
		q.pointersArray[index] = element
	} else {
		element := q.pointersArray[index]
		q.lRUList.MoveToFront(element)
	}
}

func (q *LruQueue) Back() uint32{
	val := q.lRUList.Back().Value.(uint32)
	return val
}