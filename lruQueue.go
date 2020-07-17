package caches

import (
	"container/list"
)

type lruQueue struct {
	lRUList list.List
	pointersArray[] *list.Element
}

func (q *lruQueue) Init(length uint32){
	q.pointersArray = make([]*list.Element,length)
}

func (q *lruQueue) Update(index Address){
	if  q.pointersArray[index] == nil{
		element := q.lRUList.PushFront(index)
		q.pointersArray[index] = element
	} else {
		element := q.pointersArray[index]
		q.lRUList.MoveToFront(element)
	}
}

func (q *lruQueue) Back() uint32{
	val := q.lRUList.Back().Value.(uint32)
	return val
}