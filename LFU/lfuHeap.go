package LFU

import (
	"container/heap"
)

type lfuHeap struct {
	heap             LeastFrequentUsesHeap
	cacheIndexToItem map[uint32]*Item
}

func (l *lfuHeap) Init(length uint32) {
	l.heap = make(LeastFrequentUsesHeap, length)
	heap.Init(&l.heap)
}

func (l lfuHeap) LFUIndex() uint32 {
	return l.heap.Top().(uint32)
}

func (l *lfuHeap) IncrementInIndex(index uint32) {
	item := l.cacheIndexToItem[index]
	l.heap.IncrementUses(item)
}

func (l *lfuHeap) RestartIndex(index uint32) {
	item := l.cacheIndexToItem[index]
	l.heap.RestartUses(item)
}
