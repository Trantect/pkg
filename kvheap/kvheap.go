package kvheap

import "container/heap"

// StrIntKV is a key-value pair.
// type of key is string and value is int64.
type StrIntKV struct {
	Key   string
	Value int64
}

// StrIntHeap is a heap of StrIntKV.
// It implements heap.Interface.
type StrIntHeap []StrIntKV

func New(m map[string]int64) *StrIntHeap {
	h := &StrIntHeap{}
	heap.Init(h)
	for k, v := range m {
		heap.Push(h, StrIntKV{k, v})
	}
	return h
}

// Less is greater-than here so that we can pop *larger* items.
func (h StrIntHeap) Less(i, j int) bool { return h[i].Value > h[j].Value }
func (h StrIntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h StrIntHeap) Len() int           { return len(h) }

func (h *StrIntHeap) Push(x interface{}) {
	*h = append(*h, x.(StrIntKV))
}

func (h *StrIntHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
