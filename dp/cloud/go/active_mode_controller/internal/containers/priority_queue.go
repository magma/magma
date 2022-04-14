package containers

import "container/heap"

type PriorityQueue struct {
	h items
}

func (p *PriorityQueue) Top() int {
	return p.h[0].value
}

func (p *PriorityQueue) Push(value int) *Item {
	item := &Item{value: value}
	heap.Push(&p.h, item)
	return item
}

func (p *PriorityQueue) Delete(key *Item) {
	heap.Remove(&p.h, key.index)
}

type Item struct {
	value int
	index int
}

type items []*Item

func (h items) Len() int {
	return len(h)
}

func (h items) Less(i, j int) bool {
	return h[i].value > h[j].value
}

func (h items) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}

func (h *items) Push(x interface{}) {
	n := len(*h)
	item := x.(*Item)
	item.index = n
	*h = append(*h, item)
}

func (h *items) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	item.index = -1
	*h = old[:n-1]
	return item
}
