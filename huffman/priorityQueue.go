package huffman

// this is the example from container/heap

import (
	"container/heap"
	"fmt"
)

// An Item is something we manage in a priority queue.
type Item[T interface{}] struct {
	value    T   // The value of the item; arbitrary.
	priority int // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue[T interface{}] []*Item[T]

func (pq PriorityQueue[T]) Len() int { return len(pq) }

func (pq PriorityQueue[T]) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue[T]) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue[T]) Push(x any) {
	n := len(*pq)
	item := x.(*Item[T])
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue[T]) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // don't stop the GC from reclaiming the item eventually
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue[T]) update(item *Item[T], value T, priority int) {
	item.value = value
	item.priority = priority
	heap.Fix(pq, item.index)
}

type Queuable interface {
	comparable
	fmt.Stringer
}

func makeFromMap[T Queuable](items map[T]int) *PriorityQueue[T] {
	pq := make(PriorityQueue[T], len(items))
	i := 0
	for value, priority := range items {
		pq[i] = &Item[T]{
			value:    value,
			priority: priority,
			index:    i,
		}
		log.Debugf("%s has priority %d and index %d", value, priority, i)
		i++
	}
	heap.Init(&pq)

	return &pq
}
