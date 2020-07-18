package LFU

import "container/heap"

type Item struct {
	indexInCache uint32
	numOfUses    int
	indexInHeap  int
}

type LeastFrequentUsesHeap []*Item

func (lfuHeap LeastFrequentUsesHeap) Len() int { return len(lfuHeap) }

func (lfuHeap LeastFrequentUsesHeap) Less(i, j int) bool {
	return lfuHeap[i].numOfUses > lfuHeap[j].numOfUses
}

func (lfuHeap LeastFrequentUsesHeap) Swap(i, j int) {
	lfuHeap[i], lfuHeap[j] = lfuHeap[j], lfuHeap[i]
	lfuHeap[i].indexInHeap = i
	lfuHeap[j].indexInHeap = j
}

func (lfuHeap *LeastFrequentUsesHeap) Push(x interface{}) {
	n := len(*lfuHeap)
	item := x.(*Item)
	item.indexInHeap = n
	*lfuHeap = append(*lfuHeap, item)
}

func (lfuHeap *LeastFrequentUsesHeap) Pop() interface{} {
	old := *lfuHeap
	n := len(old)
	item := old[n-1]
	old[n-1] = nil        // avoid memory leak
	item.indexInHeap = -1 // for safety
	*lfuHeap = old[0 : n-1]
	return item
}

func (lfuHeap LeastFrequentUsesHeap) Top() interface{} {
	old := lfuHeap
	n := len(lfuHeap)
	item := old[n-1]
	return item
}

func (lfuHeap *LeastFrequentUsesHeap) IncrementUses(item *Item) {
	item.numOfUses++
	heap.Fix(lfuHeap, item.indexInHeap)
}

func (lfuHeap *LeastFrequentUsesHeap) RestartUses(item *Item) {
	item.numOfUses = 1
	heap.Fix(lfuHeap, item.indexInHeap)
}
