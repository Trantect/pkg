package kvheap

import (
	"container/heap"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStrIntHeap_PopKV(t *testing.T) {
	m := map[string]int64{
		"a": 11,
		"b": 12,
		"c": 13,
	}
	h := New(m)
	assert.Equal(t, int64(13), heap.Pop(h).(StrIntKV).Value)
	assert.Equal(t, int64(12), heap.Pop(h).(StrIntKV).Value)
	assert.Equal(t, int64(11), heap.Pop(h).(StrIntKV).Value)
}
