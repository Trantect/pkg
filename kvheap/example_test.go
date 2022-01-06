package kvheap

import (
	"container/heap"
	"fmt"
)

func ExampleKVHeap() {
	h := New(map[string]int64{
		"a": 12,
		"b": 13,
	})
	h.Push(StrIntKV{"c", int64(10)})
	fmt.Printf("max v is: %v\n", (*h)[0])
	for h.Len() > 0 {
		fmt.Printf("%v ", heap.Pop(h))
	}

	// Output:
	// max v is: {b 13}
	// {b 13} {a 12} {c 10}
}
