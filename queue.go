package tinycache

import "container/heap"

type priorityQueue []*cacheItem

func newPriorityQueue() *priorityQueue {
	queue := &priorityQueue{}
	heap.Init(queue)
	return queue
}

func (pq *priorityQueue) pushItem(item *cacheItem) {
	heap.Push(pq, item)
}

func (pq *priorityQueue) updateItem(item *cacheItem) {
	heap.Fix(pq, item.index)
}

func (pq *priorityQueue) removeItem(item *cacheItem) {
	heap.Remove(pq, item.index)
}

func (pq *priorityQueue) popItem() *cacheItem {
	if pq.Len() == 0 {
		return nil
	}
	return heap.Pop(pq).(*cacheItem)
}

func (pq priorityQueue) Len() int {
	return len(pq)
}

func (pq priorityQueue) Less(i, j int) bool {
	if pq[i].expireAt.IsZero() {
		return false
	}
	if pq[j].expireAt.IsZero() {
		return true
	}
	return pq[i].expireAt.Before(pq[j].expireAt)
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *priorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*cacheItem)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[0 : n-1]
	return item
}
