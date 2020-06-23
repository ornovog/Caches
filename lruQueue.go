package caches

import "container/list"

type queue struct {
	lRUList list.List
	pointersArray[] *list.Element
}

func (q *queue) Init(length uint32){
	q.pointersArray = make([]*list.Element,length)
}

func (q *queue) Update(index uint32){
	if  q.pointersArray[index] == nil{
		element := q.lRUList.PushFront(index)
		q.pointersArray[index] = element
	} else {
		element := q.pointersArray[index]
		q.lRUList.MoveToFront(element)
	}
}

func (q *queue) Back() uint32{
	return q.lRUList.Back().Value.(uint32)
}